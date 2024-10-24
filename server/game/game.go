package game

import (
	"fmt"
	"time"
)

const (
	ORC_HEALTH = 100
	ORC_HEIGHT = 100
	ORC_WIDTH  = 100
	ORC_SPEED  = 40

	// values for facing up/down
	ORC_SWING_HEIGHT_VERTICAL = 20
	ORC_SWING_WIDTH_VERTICAL  = 50

	ORC_SWING_HEIGHT_HORIZONTAL = 50
	ORC_SWING_WIDTH_HORIZONTAL  = 20

	ORC_SWING_DAMAGE = 40

	PROJECTILE_SPEED = 5
	MAP_WIDTH        = 1500
	MAP_HEIGHT       = 700

	TIME_UNIT_MS_DEV  = 1000
	TIME_UNIT_MS_PROD = 16

	SCORE_TO_WIN = 2
)

type Game struct {
	id        uint16
	gameState GameState
}

func (g *Game) run() {
	for range time.Tick(TIME_UNIT_MS_DEV * time.Millisecond) {

		for idx := range g.gameState.orcs {
			orc := &g.gameState.orcs[idx]
			// fmt.Printf("ID: %d", idx)
			// printLocation(&orc)
			if orc.action == Action_MELEE && !orc.isSwinging {
				swing := orc.swing()
				g.gameState.meleeSwings = append(g.gameState.meleeSwings, swing)
			}

			// reset orc action if no further inputs
			orc.action = Action_NONE
		}

		for idx := range g.gameState.meleeSwings {
			swing := &g.gameState.meleeSwings[idx]
			// for each swing check if it collides w/ any orc
			for j := range g.gameState.orcs {
				orc := &g.gameState.orcs[j]
				if swing.id == orc.id {
					continue
				}
				if swing.hurtbox.isActive &&
					swing.hurtbox.hitbox.collidesWith(&orc.hitbox) {
					swing.hurtbox.isActive = false
					orc.health -= swing.hurtbox.damage
					fmt.Printf("Orc %d got hit, remaining hp %d\n", orc.id, orc.health)
					if orc.health < 0 {
						fmt.Printf("Player %d died, Player %d gets a point\n", orc.id, swing.id)
						g.gameState.scores[idx]++
						orc.init(orc.id, 50+200*orc.id, 50+200*orc.id, Direction_NONE) //reset
					}
				}
			}

			swing.activeTicksRemaining--
			fmt.Printf("swing ticks remaining from orc %d: %d\n", swing.id, swing.activeTicksRemaining)
		}

		fmt.Printf("No. of swings in play: %d\n", len(g.gameState.meleeSwings))

		for len(g.gameState.meleeSwings) > 0 && g.gameState.meleeSwings[0].activeTicksRemaining == 0 {
			g.gameState.orcs[g.gameState.meleeSwings[0].id].isSwinging = false
			g.gameState.meleeSwings = g.gameState.meleeSwings[1:]
			fmt.Print("Removed swing\n")
		}

		for idx, score := range g.gameState.scores {
			if score == SCORE_TO_WIN {
				fmt.Println("Winner: ", idx)
				return
			}
		}
		//TODO: send state data back to clients to render
	}
}

type GameState struct {
	scores      []int
	orcs        []Orc // use idx of orcs array as their ID
	projectiles []Projectile
	meleeSwings []MeleeSwing
}

func (gs *GameState) init(players *[]*Player) {
	gs.orcs = make([]Orc, len(*players))
	gs.scores = make([]int, len(*players))
	gs.meleeSwings = make([]MeleeSwing, 0)
	for idx, player := range *players {
		gs.orcs[idx] = Orc{}
		gs.orcs[idx].init(idx, 50+200*idx, 50+200*idx, Direction_NONE)
		player.controller.orc = &gs.orcs[idx]
		go player.controller.handleGameInputs()
	}
}

type Orc struct {
	id         int
	health     int
	hitbox     Hitbox
	point      Point
	direction  Direction
	action     Action
	isSwinging bool
}

func (o *Orc) init(id int, x int, y int, dir Direction) {
	o.id = id
	o.point.x = x
	o.point.y = y
	o.health = ORC_HEALTH

	// when we update location, hitbox point also updated
	// reason for this (might not be legit) is i want separation of location / hitbox
	o.hitbox.point = &o.point
	o.hitbox.height = ORC_HEIGHT
	o.hitbox.width = ORC_WIDTH

	o.direction = dir
	o.isSwinging = false
}

func (o *Orc) updateLocation(direction Direction) {
	switch direction {
	case Direction_UP:
		print("going up")
		o.point.y = max(0, o.point.y-ORC_SPEED)
	case Direction_DOWN:
		print("going down")
		o.point.y = min(MAP_HEIGHT, o.point.y+ORC_SPEED)
	case Direction_LEFT:
		print("going left")
		o.point.x = max(0, o.point.x-ORC_SPEED)
	case Direction_RIGHT:
		print("going right")
		o.point.x = min(MAP_WIDTH, o.point.x+ORC_SPEED)
	}
	printLocation(o)
}

func (o *Orc) updateAction(action Action) {
	switch action {
	case Action_MELEE:
		print("orc swing")
		o.action = Action_MELEE
	case Action_BLINK:
		print("orc blink")
		o.action = Action_BLINK
	case Action_PROJECTILE:
		print("orc fire")
		o.action = Action_PROJECTILE
	case Action_NONE:
		print("orc do nothing")
		o.action = Action_NONE
	}
}

func (o *Orc) swing() MeleeSwing {
	o.isSwinging = true
	swingLocation := o.getSwingLocation()
	swingHurtbox := Hurtbox{}

	if o.direction == Direction_UP || o.direction == Direction_DOWN {
		swingHurtbox.init(ORC_SWING_DAMAGE, swingLocation,
			ORC_SWING_HEIGHT_VERTICAL,
			ORC_SWING_WIDTH_VERTICAL)
	} else {
		swingHurtbox.init(ORC_SWING_DAMAGE, swingLocation,
			ORC_SWING_HEIGHT_HORIZONTAL,
			ORC_SWING_WIDTH_HORIZONTAL)
	}

	swing := MeleeSwing{
		id:                   o.id,
		activeTicksRemaining: 5,
		hurtbox:              swingHurtbox,
	}

	return swing
}

func (o *Orc) getSwingLocation() Point {
	swingPoint := o.point
	switch o.direction {
	case Direction_UP:
		swingPoint.y += o.hitbox.height/2 + ORC_SWING_HEIGHT_VERTICAL/2
	case Direction_DOWN:
		swingPoint.y += o.hitbox.height/2 + ORC_SWING_HEIGHT_VERTICAL/2
	case Direction_LEFT:
		swingPoint.x += o.hitbox.width/2 + ORC_SWING_HEIGHT_HORIZONTAL/2
	case Direction_RIGHT:
		swingPoint.x += o.hitbox.width/2 + ORC_SWING_HEIGHT_HORIZONTAL/2
	}
	return swingPoint
}

type MeleeSwing struct {
	id                   int
	activeTicksRemaining int
	hurtbox              Hurtbox
}

type Hurtbox struct {
	isActive bool
	damage   int
	hitbox   Hitbox
}

func (hb *Hurtbox) init(damage int, point Point, height int, width int) {
	hitbox := Hitbox{}
	hitbox.init(point, height, width)
	hb.hitbox = hitbox
	hb.damage = damage
	hb.isActive = true // true as default behaviour
}

type Hitbox struct {
	point  *Point
	height int
	width  int
}

func (hb *Hitbox) init(point Point, height int, width int) {
	hb.point = &point
	hb.height = height
	hb.width = width
}

func (hb *Hitbox) getTopRight() Point {
	return Point{x: hb.point.x + hb.width, y: hb.point.y - hb.height}
}

func (hb *Hitbox) getBottomleft() Point {
	return Point{x: hb.point.x - hb.width, y: hb.point.y + hb.height}
}

func (hb *Hitbox) getTopLeft() Point {
	return Point{x: hb.point.x - hb.width, y: hb.point.y - hb.height}
}

func (hb *Hitbox) getBottomRight() Point {
	return Point{x: hb.point.x + hb.width, y: hb.point.y + hb.height}
}

// Collision detection algorithm
// note that canvas top left (0,0), bottom right (CANVAS_WIDTH, CANVAS_HEIGHT)
func (hb *Hitbox) collidesWith(other *Hitbox) bool {
	hbTR := hb.getTopRight()
	hbBL := hb.getBottomleft()
	hbTL := hb.getTopLeft()
	hbBR := hb.getBottomRight()

	otherTR := other.getTopRight()
	otherBL := other.getBottomleft()
	otherTL := other.getTopLeft()
	otherBR := other.getBottomRight()

	isOverlappingVertice := ((hbTR.x >= otherBL.x && hbTR.y <= otherBL.y) &&
		(hbBL.x <= otherTR.x && hbBL.y >= otherBL.y)) ||
		((hbTL.x <= otherBR.x && hbTL.y <= otherBR.y) &&
			(hbBR.x >= otherTL.x && hbBR.y >= otherTL.y))

	return isOverlappingVertice
}

type Point struct {
	x int
	y int
}

func (p Point) String() string {
	return fmt.Sprintf("x: %d, y: %d", p.x, p.y)
}

type Projectile struct {
	hurtbox Hurtbox
}

func printLocation(o *Orc) {
	fmt.Printf("x: %d, y: %d\n", o.point.x, o.point.y)
}
