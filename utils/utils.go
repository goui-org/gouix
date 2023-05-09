package utils

import (
	"bytes"
	"fmt"
	"mime"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"

	"github.com/tdewolff/minify/v2"
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

var minifiable = map[string]bool{
	"text/css":               true,
	"text/html":              true,
	"application/javascript": true,
}

func CopyFile(src string, dst string, m *minify.M) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	if m != nil {
		ty := mime.TypeByExtension(path.Ext(src))
		ty = strings.TrimSpace(strings.Split(ty, ";")[0])
		if minifiable[ty] {
			out := new(bytes.Buffer)
			if err := m.Minify(ty, out, bytes.NewBuffer(data)); err != nil {
				return err
			}
			return os.WriteFile(dst, out.Bytes(), 0755)
		}
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

func CopyDirectory(scrDir, dest string, m *minify.M) error {
	entries, err := os.ReadDir(scrDir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		sourcePath := filepath.Join(scrDir, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		fileInfo, err := os.Stat(sourcePath)
		if err != nil {
			return err
		}

		stat, ok := fileInfo.Sys().(*syscall.Stat_t)
		if !ok {
			return fmt.Errorf("failed to get raw syscall.Stat_t data for '%s'", sourcePath)
		}

		switch fileInfo.Mode() & os.ModeType {
		case os.ModeDir:
			if err := createIfNotExists(destPath, 0755); err != nil {
				return err
			}
			if err := CopyDirectory(sourcePath, destPath, m); err != nil {
				return err
			}
		case os.ModeSymlink:
			if err := copySymLink(sourcePath, destPath); err != nil {
				return err
			}
		default:
			if err := CopyFile(sourcePath, destPath, m); err != nil {
				return err
			}
		}

		if err := os.Lchown(destPath, int(stat.Uid), int(stat.Gid)); err != nil {
			return err
		}

		fInfo, err := entry.Info()
		if err != nil {
			return err
		}

		isSymlink := fInfo.Mode()&os.ModeSymlink != 0
		if !isSymlink {
			if err := os.Chmod(destPath, fInfo.Mode()); err != nil {
				return err
			}
		}
	}
	return nil
}

func exists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}

	return true
}

func createIfNotExists(dir string, perm os.FileMode) error {
	if exists(dir) {
		return nil
	}

	if err := os.MkdirAll(dir, perm); err != nil {
		return fmt.Errorf("failed to create directory: '%s', error: '%s'", dir, err.Error())
	}

	return nil
}

func copySymLink(source, dest string) error {
	link, err := os.Readlink(source)
	if err != nil {
		return err
	}
	return os.Symlink(link, dest)
}
