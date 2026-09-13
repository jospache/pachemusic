package handlers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"net/http"
	"os"
	"youtube_converter/services"
	"youtube_converter/utils"
)

var getConfig = services.Config{
	APIKey:     os.Getenv("API_KEY_YOUTUBE"),
	MaxResults: 10,
}

func SearchHandler(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, services.SearchResult{Error: true, Message: "Parâmetro 'q' ausente"})
		return
	}

	maxResults := getConfig.MaxResults
	if maxParam := c.Query("max_results"); maxParam != "" {
		if _, err := fmt.Sscanf(maxParam, "%d", &maxResults); err != nil {
			c.JSON(http.StatusBadRequest, services.SearchResult{Error: true, Message: "Parâmetro 'max_results' inválido"})
			return
		}
	}

	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithAPIKey(getConfig.APIKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, services.SearchResult{Error: true, Message: fmt.Sprintf("Erro ao inicializar cliente da API YouTube: %s", err.Error())})
		return
	}

	call := service.Search.List([]string{"id", "snippet"}).
		Q(query).
		MaxResults(int64(maxResults)).
		Type("video")

	response, err := call.Do()
	if err != nil {
		c.JSON(http.StatusInternalServerError, services.SearchResult{Error: true, Message: fmt.Sprintf("Erro ao buscar vídeos no YouTube: %s", err.Error())})
		return
	}

	var results []services.VideoResult
	for _, item := range response.Items {
		videoID := item.Id.VideoId
		video, err := service.Videos.List([]string{"contentDetails", "snippet"}).Id(videoID).Do()
		if err != nil || len(video.Items) == 0 {
			continue
		}
		v := video.Items[0]

		uploadDate := v.Snippet.PublishedAt[:10]
		formattedDate := fmt.Sprintf("%s/%s/%s", uploadDate[8:10], uploadDate[5:7], uploadDate[0:4])

		duration := utils.ParseDuration(v.ContentDetails.Duration)

		videoResult := services.VideoResult{
			ID:          videoID,
			Channel:     v.Snippet.ChannelTitle,
			Title:       v.Snippet.Title,
			FullLink:    fmt.Sprintf("https://youtube.com/watch?v=%s", videoID),
			Duracao:     duration,
			PublicadoEm: formattedDate,
		}
		results = append(results, videoResult)
	}

	c.JSON(http.StatusOK, services.SearchResult{
		Error:   false,
		Results: results,
	})
}
