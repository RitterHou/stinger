package common

import "fmt"

func ByteFormat(length uint64) string {
	size := float64(length)
	units := []string{"B", "K", "M", "G", "T"}

	for _, unit := range units {
		if size < 1024 {
			return fmt.Sprintf("%7.2f%s", size, unit)
		}
		size = size / 1024.0
	}
	return "Size larger than 1024TB."
}
