package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

const rapidAPIHost = "youtube-mp36.p.rapidapi.com"

type rapidDownloadResponse struct {
	Link    string `json:"link"`
	Message string `json:"msg"`
	Status  string `json:"status"`
}

func RapidDownloadURL(videoID string) (string, error) {
	key := strings.TrimSpace(os.Getenv("RAPIDAPI_KEY"))
	if key == "" {
		return "", fmt.Errorf("RAPIDAPI_KEY não está configurada no servidor")
	}

	client := &http.Client{Timeout: 30 * time.Second}
	for attempt := 0; attempt < 45; attempt++ {
		endpoint := "https://" + rapidAPIHost + "/dl?id=" + videoID
		request, err := http.NewRequest(http.MethodGet, endpoint, nil)
		if err != nil {
			return "", fmt.Errorf("failed to create RapidAPI request: %w", err)
		}
		request.Header.Set("x-rapidapi-key", key)
		request.Header.Set("x-rapidapi-host", rapidAPIHost)
		request.Header.Set("Content-Type", "application/json")

		response, err := client.Do(request)
		if err != nil {
			return "", fmt.Errorf("RapidAPI request failed: %w", err)
		}
		var payload rapidDownloadResponse
		decodeErr := json.NewDecoder(response.Body).Decode(&payload)
		response.Body.Close()
		if decodeErr != nil {
			return "", fmt.Errorf("invalid RapidAPI response: %w", decodeErr)
		}
		if response.StatusCode == http.StatusTooManyRequests {
			return "", fmt.Errorf("limite da RapidAPI atingido; verifique a quota do plano YouTube MP3")
		}
		if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
			return "", fmt.Errorf("RapidAPI returned HTTP %d: %s", response.StatusCode, payload.Message)
		}
		if payload.Status == "ok" && payload.Link != "" {
			return payload.Link, nil
		}
		if payload.Status == "fail" {
			return "", fmt.Errorf("RapidAPI conversion failed: %s", payload.Message)
		}
		time.Sleep(time.Second)
	}
	return "", fmt.Errorf("RapidAPI conversion timed out")
}
