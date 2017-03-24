package model

import (
	"encoding/json"
	"fmt"
	_ "time"
)

type AI struct {
	RequestChannel chan Request
	Difficulty     string
	client         *Client
	Player         *Player
}

func NewIA(client *Client) AI {
	var ai AI
	ai.RequestChannel = make(chan Request, 100)
	ai.Difficulty = "easy"
	ai.Player = nil
	ai.client = client
	return ai
}

func (ai *AI) Start() {
	request := Request{}
	board := Board{}
	for {
		request = <-ai.RequestChannel
		if request.Type == "Refresh" {
			json.Unmarshal(request.Data, &board)
			if board.PlayerTurn != nil && board.PlayerTurn.Id == ai.Player.Id {
				var req =	NewRequest ("PlaceRandom")
				req.Client = ai.client
				ai.client.CurrentGame.RequestChannel <- req
			}
		} else if request.Type == "GameOver" {
			fmt.Println("Closing AI")
			return
		}
	}
}
