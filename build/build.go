package build

import (
	"bytes"
	"fmt"
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
	return &Build{
		id: gouid.String(8, gouid.MixedCaseAlpha),
	}
}

func (b *Build) TmpDir() string {
	return path.Join(os.TempDir(), b.id)
}

func (b *Build) Run() error {
	start := time.Now()
	outDir := "public"
	utils.ClearTerminal()
	fmt.Println("generating static assets...")
	if os.Getenv("DEBUG") == "true" {
		outDir = b.TmpDir()
		if !b.staticAssetsCopied {
			if err := os.RemoveAll(outDir); err != nil {
				return fmt.Errorf("build.Run: %w", err)
			}
			if err := utils.Mkdir(outDir); err != nil {
				return fmt.Errorf("build.Run: %w", err)
			}
			wasmExec, err := os.ReadFile(path.Join(runtime.GOROOT(), "misc", "wasm", "wasm_exec.js"))
			if err != nil {
				return fmt.Errorf("build.All: %w", err)
			}
			bundle := bytes.Join([][]byte{files.DebugJS, wasmExec, files.WasmFetchJS}, []byte("\n"))
			if err := utils.WriteFile(path.Join(outDir, "wasm.js"), bundle); err != nil {
				return fmt.Errorf("build.All: %w", err)
			}
			b.staticAssetsCopied = true
		}
		if err := utils.CopyFile(path.Join("public", "index.html"), path.Join(outDir, "index.html")); err != nil {
			return fmt.Errorf("build.All: %w", err)
		}
	} else {
		wasmExec, err := os.ReadFile(path.Join(runtime.GOROOT(), "misc", "wasm", "wasm_exec.js"))
		if err != nil {
			return fmt.Errorf("build.All: %w", err)
		}
		bundle := bytes.Join([][]byte{wasmExec, files.WasmFetchJS}, []byte("\n"))
		if err := utils.WriteFile(path.Join(outDir, "wasm.js"), bundle); err != nil {
			return fmt.Errorf("build.All: %w", err)
		}
	}
	fmt.Println("compiling src...")
	if err := b.compile(outDir); err != nil {
		return fmt.Errorf("build.Run: %w", err)
	}
	dur := time.Since(start).Round(time.Microsecond * 100)
	utils.ClearTerminal()
	color.Green("Built successfully in %s!\n\n", dur)
	fmt.Printf("View in your browser at %s\n\n", utils.DevServerUrl())
	fmt.Print("To create a build for production, use ")
	color.Blue("gouix build\n\n")
	fmt.Printf("Press Ctrl+C to stop\n\n")
	return nil
}

func (b *Build) compile(outDir string) error {
	return utils.Command("go", "build", "-o", path.Join(outDir, "wasm.wasm"), path.Join("src", "main.go"))
}
