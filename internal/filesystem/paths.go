package filesystem

import (
	"os"
	"path/filepath"
)

func BaseDir(configFile string) string {

	// Если запущено из IDE / go run
	if wd, err := os.Getwd(); err == nil {

		if _, err := os.Stat(filepath.Join(wd, configFile)); err == nil {
			return wd
		}
	}

	// Если запущен готовый exe (служба)
	exe, err := os.Executable()
	if err != nil {
		return "."
	}

	return filepath.Dir(exe)
}
