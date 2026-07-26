package filesystem

import (
	"os"
	"path/filepath"

	"golang.org/x/sys/windows/svc"
)

func BaseDir(configFile string) string {

	isService, _ := svc.IsWindowsService()

	// Если запущено из IDE / go run
	if wd, err := os.Getwd(); err == nil {

		if _, err := os.Stat(filepath.Join(wd, configFile)); err == nil && !isService {
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
