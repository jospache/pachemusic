package utils

import (
	"errors"
	"regexp"
)

const SupportedFormats = "mp3,mp4"

func ExtractVideoID(link string) (string, error) {
	re := regexp.MustCompile(`(?:v=|v/|youtu\.be/)([^&\n]+)`)
	matches := re.FindStringSubmatch(link)
	if len(matches) < 2 {
		return "", errors.New("ID do vídeo não encontrado")
	}
	return matches[1], nil
}
