package game

import (
	"fmt"
)

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
	for {
		packet := p.conn.getNextPacket()
		fmt.Printf("ID: %d, Type: %s, Size: %d\n", packet.Id, packet.Type, packet.Size)
		switch packet.Type {
		case "L":
			print("Lobby\n")
			if p.isReady {
				print("Player is Ready")
			} else {
				print("Player is not ready")
			}
			l, ok := packet.Data[0].(LobbyInput)
			if ok {
				p.isReady = l.IsReady
			} else {
				print("Value is not a LobbyInput")
			}
		case "G":
			print("Game\n")
			gameInputs := make([]GameInput, packet.Size)

			for idx, item := range packet.Data {
				item, ok := item.(GameInput)
				if ok {
					gameInputs[idx] = item
				} else {
					print("Value is not a GameInput")
				}
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

func (p *Player) handleGameResponses(responseCh <-chan GameResponse) {
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
	for batch := range pc.inputCh {
		for _, input := range batch.inputs {
			pc.orc.updateLocation(input.Direction)
			pc.orc.updateAction(input.Action)
		}
	}
}
