package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"youtube_converter/services"
	"youtube_converter/utils"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

var config = services.Config{
	APIKey:     os.Getenv("API_KEY_YOUTUBE"),
	MaxResults: 10,
}

type Response struct {
	Error   bool   `json:"error"`
	Message string `json:"message,omitempty"`
	File    string `json:"file,omitempty"`
	services.VideoResult
}

func ConvertHandler(c *gin.Context) {
	youtubeLink := c.Query("youtubelink")
	format := c.DefaultQuery("format", "mp3")
	taskID := c.Query("taskId")

	if youtubeLink == "" {
		c.JSON(http.StatusBadRequest, Response{Error: true, Message: "Parâmetro 'youtubelink' ausente"})
		return
	}

	if !utils.IsValidFormat(format) {
		c.JSON(http.StatusBadRequest, Response{Error: true, Message: "Formato inválido: apenas mp3, m4a ou mp4 são suportados"})
		return
	}
	format = strings.ToLower(format)

	videoID, err := utils.ExtractVideoID(youtubeLink)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: true, Message: "Nenhum vídeo especificado"})
		return
	}

	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithAPIKey(config.APIKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: true, Message: fmt.Sprintf("Erro ao inicializar cliente da API YouTube: %s", err.Error())})
		return
	}

	video, err := service.Videos.List([]string{"contentDetails", "snippet"}).Id(videoID).Do()
	if err != nil || len(video.Items) == 0 {
		c.JSON(http.StatusInternalServerError, Response{Error: true, Message: "Erro ao obter detalhes do vídeo"})
		return
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

	if rapidURL, err := services.RapidDownloadURL(videoID); err != nil {
		c.JSON(http.StatusBadGateway, Response{Error: true, Message: err.Error()})
		return
	} else if rapidURL != "" {
		c.JSON(http.StatusOK, Response{Error: false, File: rapidURL, VideoResult: videoResult})
		return
	}

	title := utils.SanitizeTitle(v.Snippet.Title)
	if title == "" {
		title = videoID
	}
	filePath := filepath.Join(DownloadFolder, title+"."+format)

	if _, err := os.Stat(filePath); err == nil {
		log.Printf("Arquivo já existente: %s", filePath)
		c.JSON(http.StatusOK, Response{
			Error:       false,
			File:        "/download/" + title + "." + format,
			VideoResult: videoResult,
		})
		return
	}

	log.Printf("Iniciando conversão: %s (%s) para formato %s", videoResult.Title, videoResult.FullLink, format)
	start := time.Now()

	if err := services.DownloadVideo(youtubeLink, format, filePath, taskID); err != nil {
		log.Printf("Erro durante a conversão do vídeo %s: %v", videoID, err)
		c.JSON(http.StatusInternalServerError, Response{Error: true, Message: err.Error()})
		return
	}

	log.Printf("Conversão concluída em %v. Arquivo salvo em: %s", time.Since(start), filePath)

	c.JSON(http.StatusOK, Response{
		Error:       false,
		File:        "/download/" + title + "." + format,
		VideoResult: videoResult,
	})
}
