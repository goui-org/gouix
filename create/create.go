package create

import (
	"fmt"
	"path"

	"github.com/twharmon/goui-cli/files"
	"github.com/twharmon/goui-cli/utils"
)

func Create(name string) error {
	fail := func(err error) error {
		return fmt.Errorf("create.Create: %w", err)
	}
	fmt.Printf("Creating %s...\n", name)
	if err := utils.Mkdir(path.Join(name, ".vscode"), path.Join(name, "src", "app"), path.Join(name, "public")); err != nil {
		return fail(err)
	}
	if err := utils.WriteFile(path.Join(name, ".vscode", "settings.json"), files.VSCodeSettingsJSON); err != nil {
		return fail(err)
	}
	if err := utils.WriteFile(path.Join(name, "public", "index.html"), files.IndexHTML); err != nil {
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
	fmt.Println("-------------------------------")
	fmt.Printf("\tcd %s\n", name)
	fmt.Println("\tgo get main/src")
	fmt.Println("\tgoui serve")
	return nil
}
