package main

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func MDToHTML(md []byte) []byte {
	p := parser.New()
	doc := p.Parse(md)
	renderer := html.NewRenderer(html.RendererOptions{})
	return markdown.Render(doc, renderer)
}
