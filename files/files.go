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

var WasmFetchJS = []byte(`let elements = {};
let nodes = new Map();
window._GOUI_ELEMENTS = elements;
let randInt = () => Math.floor(Math.random() * 2e9);
let generateId = () => {
    let id = randInt();
    while (elements[id]) id = randInt();
    return id;
};
let createElement = (tag, clicks) => {
    let id = generateId();
    let el = document.createElement(tag);
    elements[id] = el;
    if (clicks) nodes.set(el, id);
    return id;
};
let createElementNS = (tag, ns) => {
    let id = generateId();
    elements[id] = document.createElementNS(ns, tag);
    return id;
};
let createTextNode = text => {
    let id = generateId();
    elements[id] = document.createTextNode(text);
    return id;
};
let decoder = new TextDecoder();
let memory;
let getString = (addr, len) => decoder.decode(memory.buffer.slice(addr, addr + len));
let go = new Go();
go.importObject.env = {
    createElement: (addr, len, clicks) => createElement(getString(addr, len), clicks),
    createElementNS: (addr, len, addr2, len2) => createElementNS(getString(addr, len), getString(addr2, len2)),
    createTextNode: (addr, len) => createTextNode(getString(addr, len)),
    appendChild: (parent, child) => {
        elements[parent].appendChild(elements[child]);
    },
    setStr: (node, addr, len, addr2, len2) => {
        elements[node][getString(addr, len)] = getString(addr2, len2);
    },
    setClass: (node, addr, len) => {
        elements[node].className = getString(addr, len);
    },
    setData: (node, addr, len) => {
        elements[node].data = getString(addr, len);
    },
    setAriaHidden: (node, bool) => {
        elements[node].ariaHidden = !!bool;
    },
    setBool: (node, addr, len, bool) => {
        elements[node][getString(addr, len)] = !!bool;
    },
    replaceWith: (oldNode, newNode) => {
        elements[oldNode].replaceWith(elements[newNode]);
        delete elements[oldNode];
    },
    removeAttribute: (node, addr, len) => {
        elements[node].removeAttribute(getString(addr, len));
    },
    removeNode: node => {
        elements[node].remove();
    },
    disposeNode: (node, clicks) => {
        delete elements[node];
        if (clicks) nodes.delete(elements[node]);
    },
    moveBefore: (parent, nextKeyMatch, start, movingDomNode) => {
        let mdm = elements[movingDomNode];
        let par = elements[parent];
        let curr = par.childNodes[start];
        if (mdm === curr) return;
        let oldPos = mdm.nextSibling;
        par.insertBefore(mdm, curr);
        if (curr !== par.lastChild && !nextKeyMatch) {
            par.insertBefore(curr, oldPos);
        }
    },
    mount: (node, addr, len) => {
        let root = document.querySelector(getString(addr, len));
        root.appendChild(elements[node]);
    },
};

WebAssembly.instantiateStreaming(fetch('main.wasm'), go.importObject).then(o => {
    let instance = o.instance;
    let exports = instance.exports;
    window.addEventListener('click', e => {
        let target = e.target;
        while (target) {
            let node = nodes.get(target);
            if (node) {
                window._GOUI_EVENT = e;
                exports.callClickListener(node);
            }
            target = target.parentNode;
        }
    });
    memory = exports.memory;
    go.run(instance)
})`)

var IndexHTML = []byte(`<!DOCTYPE html>
<html lang="en">
    <head>
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <title>GoUI App</title>
        <link rel="stylesheet" href="main.css"> 
    </head>
    <body>
        <div id="root"></div>
    </body>
</html>`)

var ReadmeMD = []byte("# GoUI App\n\nIntall gouix\n```\ngo install github.com/twharmon/gouix@latest\n```\n\nStart the development server\n```\ngouix serve\n```\n\nCreate a production build\n```\ngouix build\n```\n")

var MainGO = []byte(`package main

import (
	"main/src/app"

	"github.com/twharmon/goui"
)

func main() {
	goui.Mount("#root", goui.Component(app.App, nil))
}
`)

var AppGO = []byte(`package app

import (
	"fmt"

	"github.com/twharmon/goui"
)

func App(goui.NoProps) *goui.Node {
	count, setCount := goui.UseState(0)

	goui.UseEffect(func() goui.EffectTeardown {
		goui.Console.Log("count is %d", count)
		return nil
	}, goui.Deps{count})

	handleIncrement := goui.UseCallback(func(e *goui.MouseEvent) {
		setCount(func(c int) int { return c + 1 })
	}, goui.Deps{})

	return goui.Element("div", &goui.Attributes{
		Class: "app",
		Children: goui.Children{
			goui.Element("button", &goui.Attributes{
				Class:    "app-btn",
				Children: goui.Children{goui.Text("increment")},
				OnClick: handleIncrement,
			}),
			goui.Element("p", &goui.Attributes{
				Children: goui.Children{goui.Text(fmt.Sprintf("count: %d", count))},
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

go 1.21

require (
    github.com/twharmon/goui v0.2.2
)
`)

var VSCodeSettingsJSON = []byte(`{
    "go.toolsEnvVars": {
        "GOARCH":"wasm",
        "GOOS":"js",
    },
    "go.installDependenciesWhenBuilding": false,
}`)

var GitIgnore = []byte(`build`)

var GoUIYML = []byte(`server:
    port: 3000
    # proxy: https://api.com
build:
    # wasm_opt: true # must have wasm_opt installed`)
