package build

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"runtime"
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

func (b *Build) BuildDir() string {
	if os.Getenv("DEBUG") == "true" {
		return path.Join(os.TempDir(), b.id)
	}
	return "build"
}

func (b *Build) Run() error {
	if os.Getenv("DEBUG") == "true" {
		return b.runDebug()
	}
	return b.runProd()
}

func (b *Build) runProd() error {
	fail := func(err error) error {
		return fmt.Errorf("build.runProd: %w", err)
	}
	start := time.Now()
	outDir := b.BuildDir()
	utils.ClearTerminal()
	fmt.Println("generating static assets...")
	if err := b.resetOutDir(); err != nil {
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
	if err := utils.WriteFile(path.Join(outDir, "main.js"), out.Bytes()); err != nil {
		return fail(err)
	}
	if err := utils.CopyDirectory("public", outDir, b.minify); err != nil {
		return fail(err)
	}
	if err := b.compile(outDir); err != nil {
		return fail(err)
	}
	dur := time.Since(start).Round(time.Microsecond * 100)
	utils.ClearTerminal()
	color.Green("Built successfully in %s!\n\n", dur)
	b.reportBuildSizes(outDir)
	fmt.Println()
	return nil
}

func (b *Build) resetOutDir() error {
	fail := func(err error) error {
		return fmt.Errorf("build.resetOutDir: %w", err)
	}
	dir := b.BuildDir()
	if err := os.RemoveAll(dir); err != nil {
		return fail(err)
	}
	if err := utils.Mkdir(dir); err != nil {
		return fail(err)
	}
	return nil
}

func (b *Build) runDebug() error {
	fail := func(err error) error {
		return fmt.Errorf("build.runDebug: %w", err)
	}
	outDir := b.BuildDir()
	start := time.Now()
	utils.ClearTerminal()
	fmt.Println("generating static assets...")
	if !b.staticAssetsCopied {
		if err := b.resetOutDir(); err != nil {
			return fail(err)
		}
		wasmExec, err := os.ReadFile(path.Join(runtime.GOROOT(), "misc", "wasm", "wasm_exec.js"))
		if err != nil {
			return fail(err)
		}
		bundle := bytes.Join([][]byte{files.DebugJS, wasmExec, files.WasmFetchJS}, []byte("\n"))
		if err := utils.WriteFile(path.Join(outDir, "main.js"), bundle); err != nil {
			return fail(err)
		}
		b.staticAssetsCopied = true
	}
	if err := utils.CopyDirectory("public", outDir, b.minify); err != nil {
		return fail(err)
	}
	if err := b.compile(outDir); err != nil {
		return fail(err)
	}
	dur := time.Since(start).Round(time.Microsecond * 100)
	utils.ClearTerminal()
	color.Green("Built successfully in %s!\n\n", dur)
	fmt.Printf("View in your browser at %s\n\n", utils.DevServerUrl())
	fmt.Print("To create a build for production, use ")
	color.Blue("gouix build\n\n")
	fmt.Printf("Press Ctrl+C to stop\n")
	fmt.Println()
	return nil
}

func (b *Build) reportBuildSizes(dir string) error {
	fail := func(err error) error {
		return fmt.Errorf("build.reportBuildSizes: %w", err)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fail(err)
	}
	fmt.Printf("%s%s%s\n", utils.PadRight("file", 15), utils.PadLeft("raw", 15), utils.PadLeft("gzip", 15))
	fmt.Printf("---------------------------------------------\n")
	for _, entry := range entries {
		fullPath := path.Join(dir, entry.Name())
		if entry.IsDir() {
			return b.reportBuildSizes(fullPath)
		}
		fi, err := os.Stat(fullPath)
		if err != nil {
			return fail(err)
		}
		gzipSize, err := utils.GzipSize(fullPath)
		if err != nil {
			return fail(err)
		}
		fmt.Printf(
			"%s%s%s\n",
			utils.PadRight(fi.Name(), 15),
			utils.PadLeft(utils.FormatFileSize(fi.Size()), 15),
			utils.PadLeft(utils.FormatFileSize(gzipSize), 15),
		)
	}
	return nil
}

func (b *Build) compile(outDir string) error {
	fmt.Println("compiling src...")
	if err := utils.Command("go", "build", "-o", path.Join(outDir, "main.wasm"), `-ldflags=-s -w`, path.Join("src", "main.go")); err != nil {
		return err
	}
	// TODO: make it to where you can opt in to use wasm-opt
	// if os.Getenv("DEBUG") != "true" {
	// 	fmt.Println("optimizing build...")
	// 	if err := utils.Command("wasm-opt", "-Oz", "--enable-bulk-memory", "-o", path.Join(outDir, "main.wasm"), path.Join(outDir, "main.wasm")); err != nil {
	// 		return nil
	// 	}
	// }
	return nil
}
