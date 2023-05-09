package utils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func Mkdir(dirs ...string) error {
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("utils.Mkdir: %w", err)
		}
	}
	return nil
}

func WriteFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0755)
}

func CopyFile(src string, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0755)
}

func Command(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("%s: %s", err, stderr.String())
	}
	return nil
}

func DevServerUrl() string {
	return fmt.Sprintf("http://localhost:%s", os.Getenv("PORT"))
}

func DevServerPort() string {
	return fmt.Sprintf(":%s", os.Getenv("PORT"))
}

var clear map[string]func()

func init() {
	clear = make(map[string]func())
	clear["linux"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["darwin"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func ClearTerminal() {
	value, ok := clear[runtime.GOOS]
	if ok {
		value()
	}
}
