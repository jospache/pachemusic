package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const rapidAPIHost = "youtube-mp310.p.rapidapi.com"

type rapidDownloadResponse struct {
	DownloadURL string `json:"downloadUrl"`
}

func RapidDownloadURL(videoLink string) (string, error) {
	key := strings.TrimSpace(os.Getenv("RAPIDAPI_KEY"))
	if key == "" {
		return "", fmt.Errorf("RAPIDAPI_KEY não está configurada no servidor")
	}

	endpoint := "https://" + rapidAPIHost + "/download/mp3?url=" + url.QueryEscape(videoLink)
	request, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create RapidAPI request: %w", err)
	}
	request.Header.Set("x-rapidapi-key", key)
	request.Header.Set("x-rapidapi-host", rapidAPIHost)

	client := &http.Client{Timeout: 45 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		return "", fmt.Errorf("RapidAPI request failed: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		if response.StatusCode == http.StatusTooManyRequests {
			return "", fmt.Errorf("limite da RapidAPI atingido: o plano gratuito permite 220 pedidos por dia; aguarde a renovação da quota ou altere o plano")
		}
		return "", fmt.Errorf("RapidAPI returned HTTP %d", response.StatusCode)
	}

	var payload rapidDownloadResponse
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return "", fmt.Errorf("invalid RapidAPI response: %w", err)
	}
	if payload.DownloadURL == "" {
		return "", fmt.Errorf("RapidAPI did not return a download URL")
	}
	return payload.DownloadURL, nil
}
