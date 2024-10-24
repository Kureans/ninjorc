package game

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type Connection struct {
	socket  *websocket.Conn
	gameCh  chan<- GameInputBatch
	lobbyCh chan<- LobbyInput
}

// for alt serialisation protocols, use readMessage then a separate serialisation fn
func (conn *Connection) getNextPacket() Packet {
	var packet Packet
	err := conn.socket.ReadJSON(&packet)
	if err != nil {
		log.Fatal(err)
	}
	return packet
}

func (conn *Connection) sendPacket(response GameResponse) {
	for _, orc := range response.Orcs {
		fmt.Printf("ID: %d, Point: %v\n", orc.Id, orc.Point)
	}
	err := conn.socket.WriteJSON(response)
	if err != nil {
		log.Fatal(err)
	}
}

type Packet struct {
	Id   int
	Type string
	Size int
	Data []InputUnion
}

type InputUnion struct {
	Lobby LobbyInput
	Game  GameInput
}

type LobbyInput struct {
	IsReady bool
}

type GameInput struct {
	Direction Direction
	Action    Action
}

type GameInputBatch struct {
	size   int
	inputs []GameInput
}

type GameResponse struct {
	Orcs []Orc
}

type Direction int8

const (
	Direction_NONE Direction = iota
	Direction_UP
	Direction_DOWN
	Direction_LEFT
	Direction_RIGHT
)

type Action int8

const (
	Action_NONE Action = iota
	Action_MELEE
	Action_PROJECTILE
	Action_BLINK
)
