package build

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/goui-org/gouix/config"
	"github.com/goui-org/gouix/files"
	"github.com/goui-org/gouix/utils"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"

	"github.com/fatih/color"
	"github.com/twharmon/gouid"
)

type Build struct {
	id                 string
	staticAssetsCopied bool
	minify             *minify.M
	config             *config.Config
	prod               bool
}

func New(cfg *config.Config) *Build {
	b := &Build{
		id:     gouid.String(8, gouid.Secure32Char),
		config: cfg,
		prod:   os.Getenv("DEBUG") != "true",
	}
	if b.prod {
		b.minify = minify.New()
		b.minify.AddFunc("text/css", css.Minify)
		b.minify.AddFunc("text/html", html.Minify)
		b.minify.AddFunc("application/javascript", js.Minify)
	}
	return b
}

func (b *Build) ReplaceConfig(cfg *config.Config) {
	b.config = cfg
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
	out, err := exec.Command("tinygo", "env", "TINYGOROOT").Output()
	if err != nil {
		return fail(err)
	}
	wasmExecPath := path.Join(strings.TrimSpace(string(out)), "targets", "wasm_exec.js")

	wasmExec, err := os.ReadFile(wasmExecPath)
	if err != nil {
		return fail(err)
	}
	bundle := bytes.Join([][]byte{wasmExec, files.WasmFetchJS}, []byte("\n"))
	if err := b.bundleIndexHTML(bundle, outDir); err != nil {
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

func (b *Build) bundleIndexHTML(js []byte, outDir string) error {
	fail := func(err error) error {
		return fmt.Errorf("build.copyIndexHTML: %w", err)
	}
	indexHTML, err := os.Open(path.Join("public", "index.html"))
	if err != nil {
		return fail(err)
	}
	indexHTMLBytes, err := io.ReadAll(indexHTML)
	if err != nil {
		return fail(err)
	}
	script := []byte("<script>")
	script = append(script, js...)
	script = append(script, []byte("</script></body>")...)
	indexHTMLBytes = bytes.Replace(indexHTMLBytes, []byte("</body>"), script, 1)
	outIndexHTML, err := os.Create(path.Join(outDir, "index.html"))
	if err != nil {
		return fail(err)
	}
	if b.minify != nil {
		if err := b.minify.Minify("text/html", outIndexHTML, bytes.NewReader(indexHTMLBytes)); err != nil {
			return fail(err)
		}
	} else {
		if _, err := outIndexHTML.Write(indexHTMLBytes); err != nil {
			return fail(err)
		}
	}
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
		out, err := exec.Command("tinygo", "env", "TINYGOROOT").Output()
		if err != nil {
			return fail(err)
		}
		wasmExecPath := path.Join(strings.TrimSpace(string(out)), "targets", "wasm_exec.js")
		wasmExec, err := os.ReadFile(wasmExecPath)
		if err != nil {
			return fail(err)
		}
		bundle := bytes.Join([][]byte{wasmExec, files.DebugJS, files.WasmFetchJS}, []byte("\n"))
		if err := b.bundleIndexHTML(bundle, outDir); err != nil {
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
	fmt.Printf("View in your browser at http://localhost:%d\n\n", b.config.Server.Port)
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
	src := path.Join("src", "main.go")
	out := path.Join(outDir, "main.wasm")

	parts := []string{
		"build",
		"-target=wasm",
		// "-gc=precise",
		"-o",
		out,
	}
	if b.prod {
		parts = append(
			parts,
			fmt.Sprintf("-panic=%s", b.config.Build.Panic),
			fmt.Sprintf("-opt=%s", b.config.Build.Opt),
		)
		if !b.config.Build.Debug {
			parts = append(parts, "-no-debug")
		}
	}
	parts = append(parts, src)
	if err := utils.Command(b.config.Build.CompilerPath, parts...); err != nil {
		return err
	}

	if b.prod && b.config.Build.WASMOpt {
		parts := []string{"-O4", "-n", "--enable-bulk-memory", "-o", out}
		if b.config.Build.NoTraps {
			parts = append(parts, "-tnh")
		}
		parts = append(parts, out)
		return utils.Command("wasm-opt", parts...)
	}
	return nil
}
