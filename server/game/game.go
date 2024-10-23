package game

import (
	"fmt"
	"log"
	"time"
)

const (
	ORC_SPEED        = 2
	PROJECTILE_SPEED = 5
	MAP_WIDTH        = 1500
	MAP_HEIGHT       = 700
)

type Game struct {
	id        uint16
	gameState GameState
}

func (g *Game) run() {
	for range time.Tick(16 * time.Millisecond) {
		// for  := range g.gameState.orcs {
		// 	// print("ID: ", idx)
		// 	// printLocation(&orc)
		// }
	}
}

type GameState struct {
	orcs        []Orc
	projectiles []Projectile
}

func (gs *GameState) init(players *[]*Player) {
	gs.orcs = make([]Orc, len(*players))
	for idx, player := range *players {
		gs.orcs[idx] = Orc{}
		gs.orcs[idx].init(200*idx, 200*idx)
		player.controller.orc = &gs.orcs[idx]
		go player.controller.handleGameInputs()
	}
}

type Orc struct {
	hitbox Hitbox
	point  Point
}

func (o *Orc) init(x int, y int) {
	o.point.x = x
	o.point.y = y
	//hard-coding orcs to be 100x100 px
	o.hitbox.vertices = []Point{
		Point{x: x - 100, y: y - 100},
		Point{x: x + 100, y: y - 100},
		Point{x: x + 100, y: y + 100},
		Point{x: x + 100, y: y + 100},
	}
}

func (o *Orc) updateLocation(direction Direction) {
	switch direction {
	case Direction_UP:
		o.point.x = max(0, o.point.x-ORC_SPEED)
	case Direction_DOWN:
		o.point.x = min(MAP_HEIGHT, o.point.x+ORC_SPEED)
	case Direction_LEFT:
		o.point.y = max(0, o.point.y-ORC_SPEED)
	case Direction_RIGHT:
		o.point.y = min(MAP_WIDTH, o.point.y+ORC_SPEED)
	}
	printLocation(o)
}

func (o *Orc) updateAction(action Action) {
	switch action {
	case Action_MELEE:
		print("orc swing")
	case Action_BLINK:
		print("orc blink")
	case Action_PROJECTILE:
		print("orc fire")
	case Action_NONE:
		print("orc do nothing")
	}
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func printLocation(o *Orc) {
	fmt.Printf("x: %d, y: %d\n", o.point.x, o.point.y)
}

type Hitbox struct {
	vertices []Point
}

type Point struct {
	x int
	y int
}

type Player struct {
	id         int
	conn       Connection
	isReady    bool
	controller PlayerController
}

func (p *Player) handleLobbyInputs(lobbyCh <-chan LobbyInput) {
	for input := range lobbyCh {
		if input.IsReady {
			p.isReady = true
		} else {
			p.isReady = false
		}
	}
}

func (p *Player) routeInputs() {
	var packet Packet
	for {
		err := p.conn.socket.ReadJSON(&packet)
		if err != nil {
			log.Fatal(err)
			return
		}
		fmt.Printf("ID: %d, Type: %s\n", packet.Id, packet.Type)
		switch packet.Type {
		case "L":
			print("Lobby\n")
			p.isReady = packet.Data[0].Lobby.IsReady
		case "G":
			print("Game\n")
			gameInputs := make([]GameInput, packet.Size)

			for idx, item := range packet.Data {
				gameInputs[idx] = item.Game
			}
			p.conn.gameCh <- GameInputBatch{
				size:   packet.Size,
				inputs: gameInputs,
			}
		default:
			panic("Invalid Packet Type")
		}

	}

}

type PlayerController struct {
	inputCh <-chan GameInputBatch
	orc     *Orc
}

func (pc *PlayerController) handleGameInputs() {
	for batch := range pc.inputCh {
		for _, input := range batch.inputs {
			pc.orc.updateLocation(input.Direction)
			pc.orc.updateAction(input.Action)
		}
	}
}

type Projectile struct {
	hitbox Hitbox
}
