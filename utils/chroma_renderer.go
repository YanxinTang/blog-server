package utils

import (
	"io"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/russross/blackfriday"
)

// ChromaRenderer implement interface of blackfriday renderer
type ChromaRenderer struct {
	html *blackfriday.HTMLRenderer
}

// RenderNode is the main content
func (r *ChromaRenderer) RenderNode(w io.Writer, node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
	if node.Type == blackfriday.CodeBlock {
		lang := string(node.CodeBlockData.Info)
		lexer := lexers.Get(lang)
		if lexer == nil {
			lexer = lexers.Fallback
		}
		lexer = chroma.Coalesce(lexer)
		style := chroma.MustNewStyle("none", chroma.StyleEntries{})
		formatter := html.New(html.WithClasses(true), html.WithLineNumbers(true), html.TabWidth(4), html.LineNumbersInTable(true))
		iterator, _ := lexer.Tokenise(nil, string(node.Literal))
		formatter.Format(w, style, iterator)
		return blackfriday.GoToNext
	}
	return r.html.RenderNode(w, node, entering)
}

// RenderHeader is a symmetric counterpart of RenderHeader.
func (r *ChromaRenderer) RenderHeader(w io.Writer, ast *blackfriday.Node) {}

// RenderFooter is a symmetric counterpart of RenderHeader.
func (r *ChromaRenderer) RenderFooter(w io.Writer, ast *blackfriday.Node) {}

// NewChromaRenderer creates and configures an HTMLRenderer object, which satisfies the Renderer interface.
func NewChromaRenderer() *ChromaRenderer {
	return &ChromaRenderer{
		html: blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
			Flags: blackfriday.TOC,
		}),
	}
}
