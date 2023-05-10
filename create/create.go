package create

import (
	"fmt"
	"os"
	"path"

	"github.com/fatih/color"
	"github.com/twharmon/gouix/files"
	"github.com/twharmon/gouix/utils"
)

func Create(name string) error {
	fail := func(err error) error {
		return fmt.Errorf("create.Create: %w", err)
	}
	fmt.Printf("creating %s...\n", name)
	if err := utils.Mkdir(path.Join(name, ".vscode"), path.Join(name, "src", "app"), path.Join(name, "public")); err != nil {
		return fail(err)
	}
	if err := utils.WriteFile(path.Join(name, ".vscode", "settings.json"), files.VSCodeSettingsJSON); err != nil {
		return fail(err)
	}
	if err := utils.WriteFile(path.Join(name, "public", "index.html"), files.IndexHTML); err != nil {
		return fail(err)
	}
	if err := utils.WriteFile(path.Join(name, "public", "main.css"), files.MainCSS); err != nil {
		return fail(err)
	}
	if err := utils.WriteFile(path.Join(name, "go.mod"), files.GoMOD); err != nil {
		return fail(err)
	}
	if err := utils.WriteFile(path.Join(name, "src", "main.go"), files.MainGO); err != nil {
		return fail(err)
	}
	if err := utils.WriteFile(path.Join(name, "src", "app", "app.go"), files.AppGO); err != nil {
		return fail(err)
	}
	if err := utils.WriteFile(path.Join(name, ".gitignore"), files.GitIgnore); err != nil {
		return fail(err)
	}
	goGetSuccess := goGet(name)
	utils.ClearTerminal()
	color.Green("Successfully created %s!\n\n", name)
	fmt.Printf("To get started, run the following commands:\n\n")
	color.Blue("\tcd %s\n", name)
	if !goGetSuccess {
		color.Blue("\tgo get main/src\n")
	}
	color.Blue("\tgouix serve\n\n")
	return nil
}

func goGet(name string) bool {
	if err := os.Chdir(name); err != nil {
		return false
	}
	if err := utils.Command("go", "get", "main/src"); err != nil {
		return false
	}
	return true
}
