package build

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"path"
	"runtime"
	"strconv"
	"time"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/twharmon/gouix/files"
	"github.com/twharmon/gouix/utils"

	"github.com/fatih/color"
	"github.com/twharmon/gouid"
)

type Build struct {
	id                 string
	staticAssetsCopied bool
	minify             *minify.M
}

func New() *Build {
	b := &Build{
		id: gouid.String(8, gouid.MixedCaseAlpha),
	}
	if os.Getenv("DEBUG") != "true" {
		b.minify = minify.New()
		b.minify.AddFunc("text/css", css.Minify)
		b.minify.AddFunc("text/html", html.Minify)
		b.minify.AddFunc("application/javascript", js.Minify)
	}
	return b
}

func (b *Build) TmpDir() string {
	return path.Join(os.TempDir(), b.id)
}

func (b *Build) Run() error {
	fail := func(err error) error {
		return fmt.Errorf("build.Run: %w", err)
	}
	start := time.Now()
	outDir := "build"
	utils.ClearTerminal()
	fmt.Println("generating static assets...")
	resetOutDir := func(dir string) error {
		if err := os.RemoveAll(outDir); err != nil {
			return fail(err)
		}
		if err := utils.Mkdir(outDir); err != nil {
			return fail(err)
		}
		return nil
	}
	if os.Getenv("DEBUG") == "true" {
		outDir = b.TmpDir()
		if !b.staticAssetsCopied {
			if err := resetOutDir(outDir); err != nil {
				return fail(err)
			}
			wasmExec, err := os.ReadFile(path.Join(runtime.GOROOT(), "misc", "wasm", "wasm_exec.js"))
			if err != nil {
				return fail(err)
			}
			bundle := bytes.Join([][]byte{files.DebugJS, wasmExec, files.WasmFetchJS}, []byte("\n"))
			if err := utils.WriteFile(path.Join(outDir, "wasm.js"), bundle); err != nil {
				return fail(err)
			}
			b.staticAssetsCopied = true
		}
	} else {
		if err := resetOutDir(outDir); err != nil {
			return fail(err)
		}
		wasmExec, err := os.ReadFile(path.Join(runtime.GOROOT(), "misc", "wasm", "wasm_exec.js"))
		if err != nil {
			return fail(err)
		}
		bundle := bytes.Join([][]byte{wasmExec, files.WasmFetchJS}, []byte("\n"))
		out := new(bytes.Buffer)
		if err := b.minify.Minify("application/javascript", out, bytes.NewBuffer(bundle)); err != nil {
			return fail(err)
		}
		if err := utils.WriteFile(path.Join(outDir, "wasm.js"), out.Bytes()); err != nil {
			return fail(err)
		}
	}
	if err := utils.CopyDirectory("public", outDir, b.minify); err != nil {
		return fail(err)
	}
	fmt.Println("compiling src...")
	if err := b.compile(outDir); err != nil {
		return fail(err)
	}
	dur := time.Since(start).Round(time.Microsecond * 100)
	utils.ClearTerminal()
	color.Green("Built successfully in %s!\n\n", dur)
	if os.Getenv("DEBUG") == "true" {
		fmt.Printf("View in your browser at %s\n\n", utils.DevServerUrl())
		fmt.Print("To create a build for production, use ")
		color.Blue("gouix build\n\n")
		fmt.Printf("Press Ctrl+C to stop\n")
	} else {
		b.reportBuildSizes(outDir)
	}
	fmt.Println()
	return nil
}

func (b *Build) reportBuildSizes(dir string) error {
	fail := func(err error) error {
		return fmt.Errorf("build.reportBuildSizes: %w", err)
	}
	padIntLeft := func(i int) string {
		s := strconv.Itoa(i)
		for len(s) < 10 {
			s = " " + s
		}
		return s
	}
	padStrRight := func(s string) string {
		max := 20
		if len(s) > max {
			s = s[len(s)-max:]
		}
		for len(s) < max {
			s = s + " "
		}
		return s
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fail(err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			return b.reportBuildSizes(path.Join(dir, entry.Name()))
		}
		fi, err := os.Stat(path.Join(dir, entry.Name()))
		if err != nil {
			return fail(err)
		}
		size := float64(fi.Size())
		fmt.Printf("\t%s%s KB\n", padStrRight(fi.Name()), padIntLeft(int(math.Round(size/1000))))
	}
	return nil
}

func (b *Build) compile(outDir string) error {
	return utils.Command("go", "build", "-o", path.Join(outDir, "wasm.wasm"), path.Join("src", "main.go"))
}
