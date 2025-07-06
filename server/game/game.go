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
	MAP_WIDTH        = 960
	MAP_HEIGHT       = 640

	TIME_UNIT_MS_DEV  = 1000 // 1 tick
	TIME_UNIT_MS_PROD = 16   // around 64 tick

	SCORE_TO_WIN = 2
)

type Game struct {
	id                uint16
	gameState         GameState
	playerResponseChs []chan GameResponse
}

func (g *Game) initResponseChannels(players *[]*Player) {
	g.playerResponseChs = make([]chan GameResponse, len(*players))
	for idx, player := range *players {
		responseCh := make(chan GameResponse)
		g.playerResponseChs[idx] = responseCh
		go player.handleGameResponses(responseCh)
	}
}

func (g *Game) run() {
	for range time.Tick(TIME_UNIT_MS_PROD * time.Millisecond) {

		for idx := range g.gameState.orcs {
			orc := &g.gameState.orcs[idx]
			// fmt.Printf("ID: %d", idx)
			// printLocation(&orc)
			if orc.Action == Action_MELEE && !orc.IsSwinging {
				swing := orc.swing()
				g.gameState.meleeSwings = append(g.gameState.meleeSwings, swing)
			}

			// reset orc action if no further inputs
			orc.Action = Action_NONE
		}

		for idx := range g.gameState.meleeSwings {
			swing := &g.gameState.meleeSwings[idx]
			// for each swing check if it collides w/ any orc
			for j := range g.gameState.orcs {
				orc := &g.gameState.orcs[j]
				if swing.id == orc.Id {
					continue
				}
				if swing.hurtbox.isActive &&
					swing.hurtbox.hitbox.collidesWith(&orc.hitbox) {
					swing.hurtbox.isActive = false
					orc.Health -= swing.hurtbox.damage
					fmt.Printf("Orc %d got hit, remaining hp %d\n", orc.Id, orc.Health)
					if orc.Health < 0 {
						fmt.Printf("Player %d died, Player %d gets a point\n", orc.Id, swing.id)
						g.gameState.scores[idx]++
						orc.init(orc.Id, 50+200*orc.Id, 50+200*orc.Id, Direction_NONE) //reset
					}
				}
			}

			swing.activeTicksRemaining--
			fmt.Printf("swing ticks remaining from orc %d: %d\n", swing.id, swing.activeTicksRemaining)
		}

		// fmt.Printf("No. of swings in play: %d\n", len(g.gameState.meleeSwings))

		for len(g.gameState.meleeSwings) > 0 && g.gameState.meleeSwings[0].activeTicksRemaining == 0 {
			g.gameState.orcs[g.gameState.meleeSwings[0].id].IsSwinging = false
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
		for _, ch := range g.playerResponseChs {
			ch <- GameResponse{Orcs: g.gameState.orcs}
		}
	}
}

type GameState struct {
	orcCount    int
	scores      []int
	orcs        []Orc // use idx of orcs array as their ID
	projectiles []Projectile
	meleeSwings []MeleeSwing
}

func (gs *GameState) init(players *[]*Player) {
	gs.orcCount = len(*players)
	gs.orcs = make([]Orc, gs.orcCount)
	gs.scores = make([]int, gs.orcCount)
	gs.meleeSwings = make([]MeleeSwing, 0)
	idToOrcLocations := make(map[int]Point, gs.orcCount)

	for idx, _ := range *players {
		idToOrcLocations[idx] = Point{X: 50 + 150*idx, Y: 50 + 150*idx}
	}

	for idx, player := range *players {
		gs.orcs[idx] = Orc{}
		gs.orcs[idx].initWithPoint(idx, idToOrcLocations[idx], Direction_NONE)

		//TODO: Add util functions to encapsulate packet creation
		ctx := GameInitContext{
			OrcCount:         gs.orcCount,
			ClientId:         idx,
			IdToOrcLocations: idToOrcLocations,
		}

		payloadArr := make([]interface{}, 1)
		payloadArr[0] = ctx
		pkt := Packet{
			Id:   idx,
			Type: "A",
			Size: 1,
			Data: payloadArr,
		}

		player.conn.sendPacket(pkt)
		player.controller.orc = &gs.orcs[idx]
		go player.controller.handleGameInputs()
	}

}

type Orc struct {
	Id         int
	Health     int
	hitbox     Hitbox
	Point      Point
	Direction  Direction
	Action     Action
	IsSwinging bool
}

func (o *Orc) init(id int, X int, Y int, dir Direction) {
	o.Id = id
	o.Point.X = X
	o.Point.Y = Y
	o.Health = ORC_HEALTH

	// when we update location, hitbox point also updated
	// reason for this (might not be legit) is i want separation of location / hitbox
	o.hitbox.point = &o.Point
	o.hitbox.height = ORC_HEIGHT
	o.hitbox.width = ORC_WIDTH

	o.Direction = dir
	o.IsSwinging = false
}

func (o *Orc) initWithPoint(id int, point Point, dir Direction) {
	o.init(id, point.X, point.Y, dir)
}

func (o *Orc) updateLocation(direction Direction) {
	switch direction {
	case Direction_UP:
		print("going up")
		o.Point.Y = max(0, o.Point.Y-ORC_SPEED)
	case Direction_DOWN:
		print("going down")
		o.Point.Y = min(MAP_HEIGHT, o.Point.Y+ORC_SPEED)
	case Direction_LEFT:
		print("going left")
		o.Point.X = max(0, o.Point.X-ORC_SPEED)
	case Direction_RIGHT:
		print("going right")
		o.Point.X = min(MAP_WIDTH, o.Point.X+ORC_SPEED)
	}
	printLocation(o)
}

func (o *Orc) updateAction(action Action) {
	switch action {
	case Action_MELEE:
		print("orc swing")
		o.Action = Action_MELEE
	case Action_BLINK:
		print("orc blink")
		o.Action = Action_BLINK
	case Action_PROJECTILE:
		print("orc fire")
		o.Action = Action_PROJECTILE
		// case Action_NONE:
		// 	print("orc do nothing")
		// 	o.Action = Action_NONE
	}
}

func (o *Orc) swing() MeleeSwing {
	o.IsSwinging = true
	swingLocation := o.getSwingLocation()
	swingHurtbox := Hurtbox{}

	if o.Direction == Direction_UP || o.Direction == Direction_DOWN {
		swingHurtbox.init(ORC_SWING_DAMAGE, swingLocation,
			ORC_SWING_HEIGHT_VERTICAL,
			ORC_SWING_WIDTH_VERTICAL)
	} else {
		swingHurtbox.init(ORC_SWING_DAMAGE, swingLocation,
			ORC_SWING_HEIGHT_HORIZONTAL,
			ORC_SWING_WIDTH_HORIZONTAL)
	}

	swing := MeleeSwing{
		id:                   o.Id,
		activeTicksRemaining: 5,
		hurtbox:              swingHurtbox,
	}

	return swing
}

func (o *Orc) getSwingLocation() Point {
	swingPoint := o.Point
	switch o.Direction {
	case Direction_UP:
		swingPoint.Y += o.hitbox.height/2 + ORC_SWING_HEIGHT_VERTICAL/2
	case Direction_DOWN:
		swingPoint.Y += o.hitbox.height/2 + ORC_SWING_HEIGHT_VERTICAL/2
	case Direction_LEFT:
		swingPoint.X += o.hitbox.width/2 + ORC_SWING_HEIGHT_HORIZONTAL/2
	case Direction_RIGHT:
		swingPoint.X += o.hitbox.width/2 + ORC_SWING_HEIGHT_HORIZONTAL/2
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
	return Point{X: hb.point.X + hb.width, Y: hb.point.Y - hb.height}
}

func (hb *Hitbox) getBottomleft() Point {
	return Point{X: hb.point.X - hb.width, Y: hb.point.Y + hb.height}
}

func (hb *Hitbox) getTopLeft() Point {
	return Point{X: hb.point.X - hb.width, Y: hb.point.Y - hb.height}
}

func (hb *Hitbox) getBottomRight() Point {
	return Point{X: hb.point.X + hb.width, Y: hb.point.Y + hb.height}
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

	isOverlappingVertice := ((hbTR.X >= otherBL.X && hbTR.Y <= otherBL.Y) &&
		(hbBL.X <= otherTR.X && hbBL.Y >= otherBL.Y)) ||
		((hbTL.X <= otherBR.X && hbTL.Y <= otherBR.Y) &&
			(hbBR.X >= otherTL.X && hbBR.Y >= otherTL.Y))

	return isOverlappingVertice
}

type Point struct {
	X int
	Y int
}

func (p Point) String() string {
	return fmt.Sprintf("X: %d, Y: %d", p.X, p.Y)
}

type Projectile struct {
	hurtbox Hurtbox
}

func printLocation(o *Orc) {
	fmt.Printf("X: %d, Y: %d\n", o.Point.X, o.Point.Y)
}
