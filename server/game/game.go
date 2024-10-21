package game

import (
	"time"
)

type Game struct {
	id        uint16
	gameState GameState
}

func (g *Game) run() {
	for range time.Tick(16 * time.Millisecond) {
		print("hi")
	}
}

type GameState struct {
	orcs        []Orc
	projectiles []Projectile
}

type Orc struct {
	hitbox Hitbox
	point  Point
}

type Hitbox struct {
	vertices []Point
}

type Point struct {
	x uint16
	y uint16
}

type Player struct {
	id         uint16
	conn       Connection
	isReady    bool
	controller PlayerController
}

type PlayerController struct {
	inputCh <-chan PlayerInput
	orc     Orc
}

type Projectile struct {
	hitbox Hitbox
}
