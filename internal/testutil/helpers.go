//go:build !production

package testutil

func IntPtr(v int) *int {
	return &v
}
