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
        <link rel="stylesheet" href="main.css"> 
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
	count, setCount := goui.UseState(0)

	goui.UseEffect(func() goui.EffectTeardown {
		godom.Console.Log("count is %d", count)
		return nil
	}, count)

	return goui.Element("div", goui.Attributes{
		Class: "app",
		Children: []*goui.Node{
			goui.Element("button", goui.Attributes{
				Class:    "app-btn",
				Children: goui.Text("increment").Slice(),
				OnClick: func(e *godom.MouseEvent) {
					setCount(func(c int) int { return c + 1 })
				},
			}),
			goui.Element("p", goui.Attributes{
				Children: goui.Text("count: %d", count).Slice(),
			}),
		},
	})
}
`)

var MainCSS = []byte(`body {
    margin: 0;
    font-family: 'Roboto', 'Helvetica Neue', sans-serif;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
}

.app {
    background-color: #282c34;
    min-height: 100vh;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    font-size: calc(10px + 2vmin);
    color: white;
}

.app-btn {
    color: #fff;
    padding: 10px 20px;
    border-radius: 5px;
    background: #25697b;
    border: 1px solid #25697b;
    cursor: pointer;
    text-transform: uppercase;
}

.app-btn:hover {
    background: #1f5c6d;
    border: 1px solid #1f5c6d;
}

.app-btn:active {
    background: #174f5f;
    border: 1px solid #174f5f;
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
