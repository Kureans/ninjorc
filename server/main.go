package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Client struct {
	clientId int
	conn     *websocket.Conn
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func initClient(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
}

func main() {
	http.HandleFunc("/", initClient)
	go http.ListenAndServe(":8080", nil)
}
