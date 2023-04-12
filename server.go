package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type WebSocketConnections struct {
	connectionsList []*websocket.Conn
}

var connections WebSocketConnections

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// atualiza a conexão HTTP para uma conexão WebSocket
	wsConnection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	connections.connectionsList = append(connections.connectionsList, wsConnection)

	// le as mensagens recebidas do WebSocket
	for {
		messageType, p, err := wsConnection.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		// envia as mensagens recebidas para todas as conexões WebSocket abertas
		for _, c := range connections.connectionsList {
			err = c.WriteMessage(messageType, p)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func main() {

	http.HandleFunc("/", wsHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// para abrir o servidor usando o wscat: wscat -c ws://localhost:8080
