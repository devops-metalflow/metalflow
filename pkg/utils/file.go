package utils

import (
	"os"
	"strconv"
	"strings"
)

// GetWorkDir 获取程序运行目录
func GetWorkDir() string {
	pwd, _ := os.Getwd()
	return pwd
}

const (
	BYTE = 1 << (10 * iota)
	KILOBYTE
	MEGABYTE
	GIGABYTE
	TERABYTE
	PETABYTE
	EXABYTE
)

func ByteSize(bytes int64) string {
	unit := ""
	value := float64(bytes)

	switch {
	case bytes >= EXABYTE:
		unit = "E"
		value /= EXABYTE
	case bytes >= PETABYTE:
		unit = "P"
		value /= PETABYTE
	case bytes >= TERABYTE:
		unit = "T"
		value /= TERABYTE
	case bytes >= GIGABYTE:
		unit = "G"
		value /= GIGABYTE
	case bytes >= MEGABYTE:
		unit = "M"
		value /= MEGABYTE
	case bytes >= KILOBYTE:
		unit = "K"
		value /= KILOBYTE
	case bytes >= BYTE:
		unit = "B"
	case bytes == 0:
		return "0"
	}

	result := strconv.FormatFloat(value, 'f', 1, 64) // nolint:gomnd
	result = strings.TrimSuffix(result, ".0")
	return result + unit
}
