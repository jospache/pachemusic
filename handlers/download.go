package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const DownloadFolder = "./downloads/"

func DownloadFileHandler(c *gin.Context) {
	fileName := filepath.Base(c.Param("filename"))
	if fileName == "." || fileName == ".." {
		c.JSON(http.StatusBadRequest, Response{Error: true, Message: "Nome de arquivo inválido"})
		return
	}
	filePath := filepath.Join(DownloadFolder, fileName)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, Response{Error: true, Message: "File not found"})
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%q", fileName))
	c.File(filePath)

	go func() {
		time.Sleep(5 * time.Second)
		if err := os.Remove(filePath); err != nil {
			log.Printf("Erro ao apagar o arquivo %s: %s", filePath, err.Error())
		}
	}()
}
