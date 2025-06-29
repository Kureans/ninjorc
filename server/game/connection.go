package game

import (
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
	err := conn.socket.ReadJSON(&packet) //some kind of optimistic marshalling? doesn't throw error if data property doesn't exist, probably default initializer called
	if err != nil {
		log.Fatal(err)
	}
	return packet
}

// func (conn *Connection) sendPacket(response GameResponse) {
// 	for _, orc := range response.Orcs {
// 		fmt.Printf("ID: %d, Point: %v\n", orc.Id, orc.Point)
// 	}
// 	err := conn.socket.WriteJSON(response)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

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
	Data []PayloadUnion
}

type PayloadUnion struct {
	Lobby LobbyInput
	Game  GameInput
	State GameResponse
}

type LobbyInput struct {
	IsReady      bool
	canStartGame bool
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
