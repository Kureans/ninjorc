package game

import (
	"fmt"
	"time"
)

const (
	ORC_HEIGHT       = 100
	ORC_WIDTH        = 100
	ORC_SPEED        = 40
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
		for idx, orc := range g.gameState.orcs {
			fmt.Printf("ID: %d", idx)
			printLocation(&orc)
			for i, other := range g.gameState.orcs {
				if idx == i {
					continue
				}
				if orc.hitbox.collidesWith(&other.hitbox) {
					print("collision detected!!!!")
				}
			}
		}
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
	// should also contain other data like health
	hitbox Hitbox
	point  Point
}

func (o *Orc) init(x int, y int) {
	o.point.x = x
	o.point.y = y
	// when we update location, hitbox point also updated
	// reason for this (might not be legit) is i want separation of location / hitbox
	o.hitbox.point = &o.point
	o.hitbox.height = ORC_HEIGHT
	o.hitbox.width = ORC_WIDTH
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

type Hitbox struct {
	point  *Point
	height int
	width  int
}

// Collision detection algorithm:
// HB1.TR(x,y) > HB2.BL(x,y) &&
// HB1.BL(x,y) < HB2.TR(x,y)
// note that canvas top left (0,0), bottom right (CANVAS_WIDTH, CANVAS_HEIGHT)
func (hb *Hitbox) collidesWith(other *Hitbox) bool {
	hbTR := Point{x: hb.point.x + hb.width, y: hb.point.y - hb.height}
	hbBL := Point{x: hb.point.x - hb.width, y: hb.point.y + hb.height}
	otherTR := Point{x: other.point.x + other.width, y: other.point.y - other.height}
	otherBL := Point{x: other.point.x - other.width, y: other.point.y + other.height}

	isOverlappingVertice := (hbTR.x >= otherBL.x && hbTR.y <= otherBL.y) &&
		(hbBL.x >= otherTR.x && hbBL.y <= otherBL.y)

	return isOverlappingVertice
}

type Point struct {
	x int
	y int
}

type Projectile struct {
	hitbox Hitbox
}

func printLocation(o *Orc) {
	fmt.Printf("x: %d, y: %d\n", o.point.x, o.point.y)
}
