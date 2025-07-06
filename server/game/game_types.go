package game

// REMEMBER TYPES SENT OVER THE WIRE VIA JSON NEED TO HAVE PROPERTIES CAPITALISED

type LobbyInputServer struct {
	CanStartGame bool
}

type LobbyInputClient struct {
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

type GameInitContext struct {
	OrcCount         int
	ClientId         int
	IdToOrcLocations map[int]Point
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
