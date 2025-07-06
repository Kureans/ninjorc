package game

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type Connection struct {
	socket  *websocket.Conn
	gameCh  chan<- GameInputBatch
	lobbyCh chan<- LobbyInputClient
}

// for alt serialisation protocols, use readMessage then a separate serialisation fn
func (conn *Connection) getNextPacket() Packet {
	var packet Packet
	err := conn.socket.ReadJSON(&packet) //some kind of optimistic marshalling? doesn't throw error if data property doesn't exist, probably default initializer called
	if err != nil {
		log.Fatal(err)
	}
	return packet
}

func (conn *Connection) sendPacket(packet Packet) {
	err := conn.socket.WriteJSON(packet)
	if err != nil {
		log.Fatal(err)
	}
}

type Packet struct {
	Id   int
	Type string
	Size int
	Data []interface{}
}

func printPayload(p *Packet) {
	switch p.Type {
	case "A":
		ctx, ok := p.Data[0].(GameInitContext)
		if !ok {
			print("Value is not a GameInitContext")
		}
		for idx, point := range ctx.IdToOrcLocations {
			fmt.Printf("Orc %d: x: %d, y: %d\n", idx, point.X, point.Y)
		}
	case "L":
		li, ok := p.Data[0].(LobbyInputClient)
		if !ok {
			print("Value is not a LobbyInput")
		}
		fmt.Print("Is Ready? ", li.IsReady)
	case "G":
		fmt.Print("Game Payload TODO")
	}
}
