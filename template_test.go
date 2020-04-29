package txtx_test

import (
	"bytes"
	"html/template"
	"io"
	"testing"

	"github.com/kylelemons/godebug/diff"
	"github.com/yosssi/gohtml"

	"github.com/crhntr/txtx"
)

var _ interface {
	ExecuteTemplate(wr io.Writer, name string, data interface{}) error
} = (*txtx.Template)(nil)

func TestNew(t *testing.T) {
	tmp, err := txtx.New(template.New("dashboard.gohtml"), bytes.NewBuffer([]byte(`<!DOCTYPE html>
<html lang="en" dir="ltr">
<head>
  <meta charset="utf-8">
  <title>{{.Title}}</title>
	<script type="text/go-template" id="paragraph">
		<a href="/{{.}}">
			{{.}}
		</a>
	</script>
	<script type="text/go-template" id="heading">
		<h1>{{.}}</h1>
	</script>
</head>
<body>
{{template "heading" .Title}}
{{- range .Pages -}}
	{{- template "paragraph" . -}}
{{- end -}}
{{.XTemplates}}
</body>
</html>`)))

	if err != nil {
		t.Error(err)
	}

	var res bytes.Buffer
	if err := tmp.ExecuteTemplate(&res, "dashboard.gohtml", struct {
		Title      string
		Pages      []string
		XTemplates template.HTML
	}{"Hello, world!", []string{"home", "about", "blog"}, tmp.XTemplates}); err != nil {
		t.Error(err)
	}

	got := gohtml.Format(res.String())
	exp := gohtml.Format(`<!DOCTYPE html><html lang="en" dir="ltr"><head>
  <meta charset="utf-8"/>
  <title>Hello, world!</title>
</head>
<body>
  <h1>
    Hello, world!
  </h1>
	<a href="/home">
		home
	</a>
	<a href="/about">
		about
	</a>
	<a href="/blog">
		blog
	</a>
	<script type="text/go-template" id="paragraph">
		<a href="/{{.}}">
			{{.}}
		</a>
	</script>
	<script type="text/go-template" id="heading">
		<h1>{{.}}</h1>
	</script>
</body></html>`)
	if got != exp {
		t.Fail()
		t.Log("got", got)
		t.Log("exp", exp)
		t.Log("diff", diff.Diff(exp, got))
	}

}
