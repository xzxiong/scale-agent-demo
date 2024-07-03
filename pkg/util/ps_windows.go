//go:build windows
// +build windows

package util

import "fmt"

func GetProcessList() ([]ProcessCmdLine, error) {
	return nil, fmt.Errorf("NOT SUPPORT.")
}
