/*
	Extract first paragraph from markdown as summary
*/

package utils

import (
	"bytes"
	"fmt"
	"io"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

// SummaryRenderer is the rendering interface
type SummaryRenderer struct {
}

// RenderNode extract first paragraph
func (r *SummaryRenderer) RenderNode(w io.Writer, node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
	if node.Parent != nil && node.Parent.Type == blackfriday.Paragraph {
		switch node.Type {
		case blackfriday.Text:
			w.Write(node.Literal)
		case blackfriday.Code:
			w.Write([]byte(fmt.Sprintf("```%s```", node.Literal)))
		}
		return blackfriday.GoToNext
	}
	if !entering && node.Type == blackfriday.Paragraph {
		return blackfriday.Terminate
	}
	return blackfriday.GoToNext
}

// RenderHeader can produce extra content before main document
func (r *SummaryRenderer) RenderHeader(w io.Writer, ast *blackfriday.Node) {
}

// RenderFooter is a symmetric counterpart of RenderHeader.
func (r *SummaryRenderer) RenderFooter(w io.Writer, ast *blackfriday.Node) {
}

// NewSummaryRenderer creates and configures an HTMLRenderer object, which satisfies the Renderer interface.
func NewSummaryRenderer() *SummaryRenderer {
	return &SummaryRenderer{}
}

func Summary(t string) string {
	input := bytes.Replace([]byte(t), []byte("\r"), nil, -1)
	output := blackfriday.Run(input, blackfriday.WithRenderer(NewSummaryRenderer()))
	return string(bluemonday.UGCPolicy().SanitizeBytes(output))
}
