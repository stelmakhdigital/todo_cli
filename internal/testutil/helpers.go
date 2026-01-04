//go:build !production

package testutil

func IntPtr(v int) *int {
	return &v
}

func StrPtr(v string) *string {
	return &v
}
