package main

import (
	"log"
	"os"
	"runtime"
	"youtube_converter/handlers"

	"github.com/gin-gonic/gin"
)

const DownloadFolder = "./downloads/"

func main() {
	if err := os.MkdirAll(DownloadFolder, 0755); err != nil {
		log.Fatalf("Failed to create download folder: %v", err)
	}

	maxProcess := runtime.NumCPU()
	runtime.GOMAXPROCS(maxProcess)
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	if err := router.SetTrustedProxies(nil); err != nil {
		log.Fatalf("Failed to set trusted proxies: %v", err)
	}

	router.Use(handlers.CORS())

	router.GET("/", func(c *gin.Context) {
		c.File("./site/index.html")
	})
	router.StaticFile("/script.js", "./site/script.js")
	router.StaticFile("/style.css", "./site/style.css")
	router.GET("/convert", handlers.ConvertHandler)
	router.GET("/download/:filename", handlers.DownloadFileHandler)
	router.GET("/search", handlers.SearchHandler)
	router.GET("/ws", handlers.HandleConnections)

	go handlers.HandleMessages()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
