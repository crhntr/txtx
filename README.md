# txtx

This package parses a go html template.
Finds "script" elements with the tag "text/go-template".
Provides a string of type template.HTML so that the template can be used again client side.
Useful for GOOS=js GOARCH=wasm + Go templates.

## (wip) Example
Please note github.com/crhntr/dom is unstable. the following may not compile.
It is included to convey why you might want to copy some code from this package.

```html
<html>
<head>
	<title>Some Page</title>
	<!-- Load/start wasm_exec.js and compiled wasm binary (for main included below)-->
</head>

<body>
    <div id="main">
        {{template "greeting" . }}
    </div>

    {{.XTemplates}}
</body>

<script type="text/go-template" id="greeting">
    <h1>{{.Message}}</h1>
</script>

</html>
```

```go
package server

import (
	"html/template"
	"net/http"
	"os"

	"github.com/crhntr/txtx"
)

func SomePage() http.HandlerFunc {	
    return func(res http.ResponseWriter, req *http.Request) {
    	f, _ := os.Open("pages/some-page/index.html")
        tmp, _ := txtx.New(template.New(""), f)
 
        var data = struct{
            XTemplates template.HTML
            Message string
        } {
            XTemplates: tmp.XTemplates,
            Message: "Hello, world!",
        }
        
        res.WriteHeader(http.StatusOK)
        _ = tmp.ExecuteTemplate(res, "index.html", data)
    }
}
```

```go
// +build js wasm

package main

import (
    "html/template"
    "net/http"
    "os"
    "time"

    "github.com/crhntr/dom"
)

func main() {
    tmp, _ := dom.LoadTemplates((*template.Template)(nil), "")
    time.Sleep(time.Second * 5)
    mainEl = dom.GetElementByID("main")
    greeting, _ := dom.NewElementFromTemplate(tmp, "greeting", struct{
        Message string
    } { "Hola, mundo!" })
    mainEl.SetInnerHTML("")
    mainEl.Call("appendChild", greeting)
}
```
