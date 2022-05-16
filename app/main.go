package main

import (
	"bytes"
	"errors"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/gomarkdown/markdown"
	"gopkg.in/yaml.v2"
)

type Params struct {
	Input  string
	Output string
}

// parse command line arguments
func parseArgs() *Params {
	switch len(os.Args) {
	case 2:
		return &Params{
			Input:  os.Args[1],
			Output: filepath.Base(os.Args[1]) + ".html",
		}
	case 3:
		return &Params{
			Input:  os.Args[1],
			Output: os.Args[2],
		}
	default: // invalid number of arguments
		panic(`Usage: template-gen <input> <output>
			invalid number of arguments`)
	}
}

const articleTempl = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
	.max-width-md {
		max-width: 768px;
	}
	.article img {
		display: block;
		margin: auto;
		height: 100%;
		width: 100%;
		object-fit: cover;
	}
	.article blockquote {
		padding: 0 1em;
		color: gray;
		border-left: .25em solid gray;
	}
	code.has-jax {
		font: inherit;
		font-size: 100%;
		background: inherit;
		border: inherit;
		color: #515151;
	}
	</style>
	<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css" rel="stylesheet"
		integrity="sha384-1BmE4kWBq78iYhFldvKuhfTAU6auU8tT94WrHftjDbrCEXSU1oBoqyl2QvZ6jIW3" crossorigin="anonymous">
	<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap-icons@1.8.1/font/bootstrap-icons.css">
	<script src="https://polyfill.io/v3/polyfill.min.js?features=es6"></script>
	<script id="MathJax-script" async src="https://cdn.jsdelivr.net/npm/mathjax@3/es5/tex-mml-chtml.js"></script>
	<script id="MathJax-script" type="text/javascript">
	MathJax = {
		tex: {
			inlineMath: [
				["$", "$"],
				["\\(", "\\)"],
			],
		},
		svg: {
			fontCache: "global",
		},
	};
	</script>
	<script type="text/javascript" id="MathJax-script" async
		src="https://cdn.jsdelivr.net/npm/mathjax@3/es5/tex-svg.js"></script>
</head>
<header
	class="container d-flex flex-wrap align-items-center justify-content-center justify-content-md-between py-3 mb-4 border-bottom">
	<a href="/" class="d-flex align-items-center col-md-3 mb-2 mb-md-0 text-dark text-decoration-none">
		<svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" fill="currentColor"
			class="bi bi-bootstrap-fill" viewBox="0 0 16 16">
			<path
				d="M6.375 7.125V4.658h1.78c.973 0 1.542.457 1.542 1.237 0 .802-.604 1.23-1.764 1.23H6.375zm0 3.762h1.898c1.184 0 1.81-.48 1.81-1.377 0-.885-.65-1.348-1.886-1.348H6.375v2.725z" />
			<path
				d="M4.002 0a4 4 0 0 0-4 4v8a4 4 0 0 0 4 4h8a4 4 0 0 0 4-4V4a4 4 0 0 0-4-4h-8zm1.06 12V3.545h3.399c1.587 0 2.543.809 2.543 2.11 0 .884-.65 1.675-1.483 1.816v.1c1.143.117 1.904.931 1.904 2.033 0 1.488-1.084 2.396-2.888 2.396H5.062z" />
		</svg>
	</a>

	<ul class="nav col-12 col-md-auto mb-2 justify-content-center mb-md-0">
		<li><a href="/courses" class="nav-link px-2 link-dark">
				<i class="bi bi-book"></i> Курси</a></li>
		<li><a href="/guides" class="nav-link px-2 link-dark">
				<i class="bi bi-question-circle"></i> Як допомогти</a></li>
		<li><a href="/about" class="nav-link px-2 link-dark">
				<i class="bi bi-people-fill"></i>
				Про нас</a></li>
	</ul>
</header>
<div class="article container-fluid max-width-md">
<h1 class="title">{{ .Title }}</h1>
<div class="author">
	<a href="#" class="btn">
		<i class="bi bi-person-circle"></i>
		{{ .Author }}
	</a>
</div>
<hr />
<p class="text-justify">{{ .Content }}</p>
</div>
<footer class="container d-flex justify-content-between align-items-center py-3 my-4 border-top">
    <div class="col-md-4 d-flex align-items-center">
        <a href="/" class="mb-3 me-2 mb-md-0 text-muted text-decoration-none lh-1">
            <svg class="bi" width="30" height="24">
                <use xlink:href="#bootstrap"></use>
            </svg>
        </a>
        <span class="text-muted">© 2021 Company, Inc</span>
    </div>

    <ul class="nav col-md-4 justify-content-end list-unstyled d-flex">
        <li class="ms-3">
            <a class="text-muted" href="#">
                <i class="bi bi-github"></i>
            </a>
        </li>
    </ul>
</footer>
</html>`

type Article struct {
	Title   string `yaml:"title"`
	Content template.HTML
	Author  string `yaml:"author"`
}

// parse buffer as yaml and return Article struct
func ParseArticleData(buf []byte) (*Article, error) {
	var data Article
	err := yaml.Unmarshal(buf, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// parse buffer to as markdown and yaml separated by '---'
func ParseArticle(buf []byte) (a *Article, err error) {
	yamlStart := bytes.Index(buf, []byte("---"))

	if yamlStart == -1 {
		return nil, errors.New("yaml separator not found")
	}
	yamlEnd := bytes.Index(buf[yamlStart+3:], []byte("---"))
	if yamlEnd == -1 {
		return nil, errors.New("yaml separator not found")
	}
	a, err = ParseArticleData(buf[yamlStart+3 : yamlStart+3+yamlEnd])
	if err != nil {
		return nil, err
	}
	a.Content = template.HTML(markdown.ToHTML(buf[yamlStart+6+yamlEnd:], nil, nil))
	log.Println(a.Content)
	return a, nil
}

// parse markdown to html and insert it to template
func BuildHTML(buff []byte) ([]byte, error) {
	a, err := ParseArticle(buff)
	if err != nil {
		return nil, err
	}
	html := template.Must(template.New("article").Parse(articleTempl))
	var buffHTML bytes.Buffer
	err = html.Execute(&buffHTML, a)
	if err != nil {
		return nil, err
	}
	return buffHTML.Bytes(), nil
}

func main() {
	params := parseArgs()
	input, err := ioutil.ReadFile(params.Input)
	if err != nil {
		panic(err)
	}
	output, err := BuildHTML(input)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(params.Output, []byte(output), 0644)
	if err != nil {
		panic(err)
	}
}
