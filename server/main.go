package main

import (
	"log"
	"net/http"

	"github.com/Kureans/ninjorc/server/game"
	"github.com/gorilla/websocket"
)

type Client struct {
	clientId int
	conn     *websocket.Conn
}

func main() {
	manager := game.LobbyManager{}
	manager.Init()
	http.HandleFunc("/", manager.HandleNewClient)
	print("Listening on port 8080...\n")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
