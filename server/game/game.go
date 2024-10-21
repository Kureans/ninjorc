package game

import (
	"fmt"
	"time"
)

type Game struct {
	id        uint16
	gameState GameState
}

func (g *Game) run() {
	for range time.Tick(16 * time.Millisecond) {
		for idx, orc := range g.gameState.orcs {
			print("ID: ", idx)
			printLocation(&orc)
		}
	}
}

type GameState struct {
	orcs        []Orc
	projectiles []Projectile
}

func (gs *GameState) init(players *[]Player) {
	gs.orcs = make([]Orc, len(*players))
	for idx, player := range *players {
		gs.orcs[idx] = Orc{}
		gs.orcs[idx].init(200*idx, 200*idx)
		player.controller.orc = &gs.orcs[idx]
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

type PlayerController struct {
	inputCh <-chan PlayerInput
	orc     *Orc
}

type Projectile struct {
	hitbox Hitbox
}
