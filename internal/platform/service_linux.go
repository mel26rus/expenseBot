//go:build !windows

package platform

func IsService() (bool, error) {
	return false, nil
}
