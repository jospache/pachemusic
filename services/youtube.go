package services

import "os"

type SearchResult struct {
	Error   bool          `json:"error"`
	Message string        `json:"message,omitempty"`
	Results []VideoResult `json:"results,omitempty"`
}

type VideoResult struct {
	ID          string `json:"id"`
	Channel     string `json:"channel"`
	Title       string `json:"title"`
	FullLink    string `json:"full_link"`
	Duracao     string `json:"duracao"`
	PublicadoEm string `json:"publicado_em"`
}

type Config struct {
	APIKey     string
	MaxResults int
}

var getConfig = Config{
	APIKey:     os.Getenv("API_KEY_YOUTUBE"),
	MaxResults: 10,
}
