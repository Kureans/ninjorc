package game

import "github.com/gorilla/websocket"

type Connection struct {
	socket   *websocket.Conn
	clientCh chan<- ClientInput
	playerCh chan<- PlayerInput
}

type ClientInput struct {
	isReady bool
}

type PlayerInput struct {
	direction Direction
	action    Action
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
