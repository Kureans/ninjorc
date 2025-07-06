package game

import (
	"errors"
)

func toLobbyInputClient(data *[]interface{}) (LobbyInputClient, error) {
	item, ok := (*data)[0].(map[string]interface{})

	if !ok {
		return LobbyInputClient{}, errors.New("failed to convert to LobbyInputClient")
	}

	isReady, ok := item["IsReady"].(bool)
	if !ok {
		return LobbyInputClient{}, errors.New("failed to convert to LobbyInputClient")
	}

	return LobbyInputClient{IsReady: isReady}, nil
}

// type TaggedUnion struct {
// 	tag   string
// 	value interface{}
// }

func toGameInputBatch(packet *Packet) (GameInputBatch, error) {
	gameInputs := make([]GameInput, packet.Size)
	for idx, item := range packet.Data {
		item, ok := item.(map[string]interface{})
		var input GameInput
		if !ok {
			return GameInputBatch{}, errors.New("value is not a GameInput")
		}

		if item["tag"] == "D" {
			input = GameInput{Direction: Direction(int(item["value"].(float64)))}
		} else if item["tag"] == "A" {
			input = GameInput{Action: Action(int(item["value"].(float64)))}
		} else {
			return GameInputBatch{}, errors.New("value is not a GameInput")
		}

		gameInputs[idx] = input
	}

	return GameInputBatch{
		size:   packet.Size,
		inputs: gameInputs,
	}, nil
}
