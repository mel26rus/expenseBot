//go:build windows

package platform

import "golang.org/x/sys/windows/svc"

func IsService() (bool, error) {
	return svc.IsWindowsService()
}
