package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
	"youtube_converter/common"
)

var (
	upgrade = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	clients   = make(map[*websocket.Conn]string)
	broadcast = common.Broadcast
	mutex     sync.Mutex
)

func HandleConnections(c *gin.Context) {
	taskID := c.Query("taskId")
	ws, err := upgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Erro ao atualizar para WebSocket: %v", err)
		return
	}
	defer func(ws *websocket.Conn) {
		err := ws.Close()
		if err != nil {
			log.Printf("Erro ao fechar conexão WebSocket: %v", err)
		}
	}(ws)

	log.Printf("Nova conexão WebSocket estabelecida para o cliente: %s", taskID)

	mutex.Lock()
	clients[ws] = taskID
	mutex.Unlock()

	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			if err := ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("Erro ao enviar ping: %v", err)
				return
			}
		}
	}()

	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			log.Printf("Erro ao ler mensagem: %v", err)
			break
		}
	}

	mutex.Lock()
	delete(clients, ws)
	mutex.Unlock()
}

func HandleMessages() {
	for msg := range broadcast {
		log.Printf("Enviando mensagem para o cliente: %s", msg.TaskID)

		mutex.Lock()
		for client, id := range clients {
			if id == msg.TaskID {
				err := client.WriteJSON(msg)
				if err != nil {
					log.Printf("Erro ao enviar mensagem: %v", err)
					err := client.Close()
					if err != nil {
						log.Printf("Erro ao fechar conexão do cliente: %v", err)
					}
					delete(clients, client)
				}
			}
		}
		mutex.Unlock()
	}
}
