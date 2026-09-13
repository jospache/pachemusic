package utils

import (
	"fmt"
	"regexp"
	"strings"
)

func ParseDuration(duration string) string {
	re := regexp.MustCompile(`PT(\d+H)?(\d+M)?(\d+S)?`)
	matches := re.FindStringSubmatch(duration)
	if len(matches) < 4 {
		return "00:00:00"
	}

	hours := strings.TrimSuffix(matches[1], "H")
	minutes := strings.TrimSuffix(matches[2], "M")
	seconds := strings.TrimSuffix(matches[3], "S")

	if hours == "" {
		hours = "00"
	} else if len(hours) == 1 {
		hours = "0" + hours
	}
	if minutes == "" {
		minutes = "00"
	} else if len(minutes) == 1 {
		minutes = "0" + minutes
	}
	if seconds == "" {
		seconds = "00"
	} else if len(seconds) == 1 {
		seconds = "0" + seconds
	}

	return fmt.Sprintf("%s:%s:%s", hours, minutes, seconds)
}

func SanitizeTitle(title string) string {
	reg := regexp.MustCompile(`[^\w\s-]`)
	sanitized := reg.ReplaceAllString(title, "")
	sanitized = strings.ReplaceAll(sanitized, " ", "_")
	return sanitized
}
