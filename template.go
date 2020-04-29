package txtx

import (
	"bytes"
	"html/template"
	"io"

	"golang.org/x/net/html"
)

type Template struct {
	indexName string

	*template.Template

	XTemplates template.HTML
}

func New(tmpl *template.Template, r io.Reader) (*Template, error) {
	indexNode, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	xtmps := make(map[string]*html.Node)

	if err := findAndRemoveTemplates(xtmps, indexNode); err != nil {
		return nil, err
	}

	var tmp Template

	var indexWithoutTemplates bytes.Buffer
	if err := html.Render(&indexWithoutTemplates, indexNode); err != nil {
		return nil, err
	}

	tmp.Template, err = tmpl.Parse(html.UnescapeString(indexWithoutTemplates.String()))
	if err != nil {
		return nil, err
	}

	for id, nd := range xtmps {
		var buf bytes.Buffer
		if err := html.Render(&buf, nd.FirstChild); err != nil {
			return nil, err
		}

		t, err := template.New(id).Parse(html.UnescapeString(buf.String()))
		if err != nil {
			return nil, err
		}

		if tmp.Template, err = tmp.Template.AddParseTree(id, t.Tree); err != nil {
			return nil, err
		}
	}

	for _, nd := range xtmps {
		var xt bytes.Buffer
		if err := html.Render(&xt, nd); err != nil {
			return nil, err
		}
		tmp.XTemplates += "\n" + template.HTML(xt.String())
	}

	return &tmp, nil
}

func findAndRemoveTemplates(xtmps map[string]*html.Node, node *html.Node) error {
	for ; node != nil; node = node.NextSibling {
		if node.Type == html.ElementNode &&
			node.Data == "script" &&
			hasAttTypeXTemplate(node.Attr) {

			xtmps[getID(node.Attr)] = node

			if node.Parent != nil {
				next := node.PrevSibling
				node.Parent.RemoveChild(node)
				node = next
			}

			continue
		}

		if node.FirstChild != nil {
			if err := findAndRemoveTemplates(xtmps, node.FirstChild); err != nil {
				return err
			}
		}
	}

	return nil
}

func hasAttTypeXTemplate(attr []html.Attribute) bool {
	for _, at := range attr {
		if at.Key == "type" && at.Val == "text/x-template" {
			return true
		}
	}
	return false
}

func getID(attr []html.Attribute) string {
	for _, at := range attr {
		if at.Key == "id" {
			return at.Val
		}
	}
	return ""
}
