package files

var DebugJS = []byte(`const ws = new WebSocket('ws://' + window.location.host + '/hot')
ws.onmessage = e => {
	if (e.data === 'reload') return window.location.reload()
	const overlay = document.createElement('div')
	overlay.style = 'position: fixed; left: 0; right: 0; top: 0; bottom: 0; background: #000c; color: #e77; font-size: 18px'
	const msg = document.createElement('div')
	msg.style = 'width: 100%; max-width: 600px; margin: auto; margin-top: 5vh; line-height: 200%; padding: 0 20px;'
	msg.innerText = e.data
	overlay.appendChild(msg)
	document.body.appendChild(overlay)
}
ws.onopen = () => ws.send('loaded')
ws.onclose = () => window.close()
`)

var WasmFetchJS = []byte(`const go = new Go()
const fetched = fetch('/wasm.wasm')
if ('instantiateStreaming' in WebAssembly) {
    WebAssembly.instantiateStreaming(fetched, go.importObject).then(o => go.run(o.instance))
} else {
    fetched.then(r => r.arrayBuffer()).then(bytes =>
        WebAssembly.instantiate(bytes, go.importObject).then(o => go.run(o.instance))
    )
}`)

var IndexHTML = []byte(`<!DOCTYPE html>
<html lang="en">
    <head>
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <title>GoUI App</title>
    </head>
    <body>
        <div id="root"></div>
        <script src="wasm.js"></script>
    </body>
</html>`)

var MainGO = []byte(`package main

import (
	"main/src/app"

	"github.com/twharmon/goui"
)

func main() {
	goui.Mount("#root", goui.Component(app.App, goui.NoProps))
}
`)

var AppGO = []byte(`package app

import (
	"github.com/twharmon/godom"
	"github.com/twharmon/goui"
)

func App(_ any) *goui.Node {
	cnt, setCnt := goui.UseState(0)

	goui.UseEffect(func() goui.EffectTeardown {
		godom.Console.Log("count changed to %d", cnt)
		return nil
	}, cnt)

	return goui.Element("div", goui.Attributes{
		Children: []*goui.Node{
			goui.Element("button", goui.Attributes{
				Children: goui.Text("increment").Slice(),
				OnClick: func(e *godom.MouseEvent) {
					setCnt(func(c int) int { return c + 1 })
				},
			}),
			goui.Element("p", goui.Attributes{
				Children: goui.Text("cnt: %d", cnt).Slice(),
			}),
		},
	})
}
`)

var GoMOD = []byte(`module main

go 1.20

require (
    github.com/twharmon/godom v0.0.5
    github.com/twharmon/goui v0.0.2
)
`)

var VSCodeSettingsJSON = []byte(`{
    "go.toolsEnvVars": {
        "GOARCH":"wasm",
        "GOOS":"js",
    },
    "go.installDependenciesWhenBuilding": false,
}`)

var GitIgnore = []byte(`.DS_Store
build`)
