package common

import "fmt"

func ByteFormat(length uint64) string {
	size := float64(length)
	if size < 1024 {
		return fmt.Sprintf("%7.2fB", size)
	}
	size = size / 1024.0
	if size < 1024 {
		return fmt.Sprintf("%7.2fK", size)
	}
	size = size / 1024.0
	if size < 1024 {
		return fmt.Sprintf("%7.2fM", size)
	}
	size = size / 1024.0
	if size < 1024 {
		return fmt.Sprintf("%7.2fG", size)
	}
	return "Too larger size."
}
