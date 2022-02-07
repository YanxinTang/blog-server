package controller

import (
	"bytes"
	"fmt"
	"html/template"
)

type Pagination struct {
	Total   uint64
	Current uint64
	PerPage uint64
	URL     string
	Links   *[]Link
}

type Link struct {
	Active bool
	URL    string
	Index  uint64
}

const tmpl string = `
{{if .links}}
	<ul class="pagination">
    {{range .links}}
			<li class="pagination-item {{if .Active}}active{{end}}">
				<a class="pagination-item__link" href="{{.URL}}">{{.Index}}</a>
			</li>
    {{end}}
	</ul>
{{end}}
`

func NewPagination(total, current, perpage uint64, url string) *Pagination {
	pagination := &Pagination{
		total,
		current,
		perpage,
		url,
		nil,
	}
	return pagination
}

func (p *Pagination) Render() template.HTML {
	var out bytes.Buffer
	tPagination := template.Must(template.New("pagination").Parse(tmpl))
	tMap := map[string]interface{}{
		"links": p.Links,
	}
	tPagination.Execute(&out, tMap)
	return template.HTML(out.String())
}

func (p *Pagination) Init() *Pagination {
	var i uint64
	var links []Link
	for i = 0; i*p.PerPage < p.Total; i++ {
		links = append(links, Link{
			Index:  i + 1,
			URL:    fmt.Sprintf(p.URL, i+1),
			Active: i+1 == p.Current,
		})
	}
	p.Links = &links
	return p
}
