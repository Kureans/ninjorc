package game

import "fmt"

type Player struct {
	id         int
	conn       Connection
	isReady    bool
	controller PlayerController
}

func (p *Player) handleLobbyInputs(lobbyCh <-chan LobbyInputClient) {
	for input := range lobbyCh {
		if input.IsReady {
			p.isReady = true
		} else {
			p.isReady = false
		}
	}
}

func (p *Player) routeInputs() {
	for {
		packet := p.conn.getNextPacket()
		switch packet.Type {
		case "L":
			input, err := toLobbyInputClient(&packet.Data)
			if err != nil {
				print(err)
			}
			p.isReady = input.IsReady

		case "G":
			input, err := toGameInputBatch(&packet)
			if err != nil {
				print(err)
			}
			p.conn.gameCh <- input

		default:
			panic("Invalid Packet Type")
		}

	}
}

func (p *Player) handleGameResponses(responseCh <-chan GameResponse) { // tells indiv players how to draw the orcs
	for response := range responseCh {
		payloadArr := make([]interface{}, 1)
		payloadArr[0] = response
		packet := Packet{p.id, "G", len(response.Orcs), payloadArr}
		p.conn.sendPacket(packet)
	}
}

type PlayerController struct {
	inputCh <-chan GameInputBatch
	orc     *Orc
}

func (pc *PlayerController) handleGameInputs() {
	fmt.Printf("Buffer size: %d", len(pc.inputCh))
	for batch := range pc.inputCh {
		for _, input := range batch.inputs {
			pc.orc.updateLocation(input.Direction)
			pc.orc.updateAction(input.Action)
		}
	}
}
