package build

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/twharmon/gouix/files"
	"github.com/twharmon/gouix/utils"

	"github.com/fatih/color"
	"github.com/twharmon/gouid"
)

type Build struct {
	id                 string
	staticAssetsCopied bool
}

func New() *Build {
	b := &Build{
		id: gouid.String(8, gouid.MixedCaseAlpha),
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
	if os.Getenv("DEBUG") == "true" {
		outDir = b.TmpDir()
		if !b.staticAssetsCopied {
			if err := os.RemoveAll(outDir); err != nil {
				return fail(err)
			}
			if err := utils.Mkdir(outDir); err != nil {
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
		if err := utils.Mkdir(outDir); err != nil {
			return fail(err)
		}
		wasmExec, err := os.ReadFile(path.Join(runtime.GOROOT(), "misc", "wasm", "wasm_exec.js"))
		if err != nil {
			return fail(err)
		}
		bundle := bytes.Join([][]byte{wasmExec, files.WasmFetchJS}, []byte("\n"))
		if err := utils.WriteFile(path.Join(outDir, "wasm.js"), bundle); err != nil {
			return fail(err)
		}
	}
	if err := utils.CopyDirectory("public", outDir); err != nil {
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
		fmt.Printf("Press Ctrl+C to stop\n\n")
	} else {
		wasmFI, err := os.Stat(path.Join(outDir, "wasm.wasm"))
		if err != nil {
			return fail(err)
		}
		jsFI, err := os.Stat(path.Join(outDir, "wasm.js"))
		if err != nil {
			return fail(err)
		}
		wasmSize := float64(wasmFI.Size())
		jsSize := float64(jsFI.Size())
		fmt.Printf("\twasm.wasm:\t%d KB\n", int(math.Round(wasmSize/1000)))
		fmt.Printf("\twasm.js:\t%d KB\n", int(math.Round(jsSize/1000)))
		fmt.Println()
	}
	return nil
}

func (b *Build) compile(outDir string) error {
	return utils.Command("go", "build", "-o", path.Join(outDir, "wasm.wasm"), path.Join("src", "main.go"))
}
