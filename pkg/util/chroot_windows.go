//go:build windows
// +build windows

package util

const RootFS = "none"

// Chroot chroot()s to the new path.
// NB: All file paths after this call are effectively relative to
// `rootfs`
func Chroot(rootfs string) error {
	return errors.New("chroot is not implemented on Windows")
}
