package files

import _ "embed"

//go:embed wasmfetch.js
var WasmFetchJS []byte

//go:embed debug.js
var DebugJS []byte

//go:embed index.html
var IndexHTML []byte

//go:embed readme.md
var ReadmeMD []byte

//go:embed main.go_
var MainGO []byte

//go:embed app.go_
var AppGO []byte

//go:embed main.css
var MainCSS []byte

//go:embed go.mod_
var GoMOD []byte

//go:embed vscodesettings.json
var VSCodeSettingsJSON []byte

//go:embed gitignore
var GitIgnore []byte

//go:embed goui.yml
var GoUIYML []byte
