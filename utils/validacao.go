package utils

import "strings"

func IsValidFormat(format string) bool {
	switch strings.ToLower(format) {
	case "mp3", "mp4":
		return true
	default:
		return false
	}
}
