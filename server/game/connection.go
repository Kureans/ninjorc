package game

import "github.com/gorilla/websocket"

type Connection struct {
	socket  *websocket.Conn
	gameCh  chan<- GameInputBatch
	lobbyCh chan<- LobbyInput
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
