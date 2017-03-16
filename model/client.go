package model

import (
	"github.com/gorilla/websocket"
	"encoding/json"
	"fmt"
)

type Client struct {
	Conn        *websocket.Conn
	token       string
	State       string
	Ai          *AI
	CurrentGame *Game
}

func NewClient(conn *websocket.Conn) *Client {
	var client Client
	client.Conn = conn
	client.token = ""
	client.State = "Start"
	client.Ai = nil
	return &client
}

func NewAiClient() *Client {
	var client Client
	client.Conn = nil
	client.token = ""
	client.State = "Start"
	ai := NewIA(&client)
	client.Ai = &ai
	return &client
}

func (client *Client) IsAi() bool {
	if client.Ai == nil {
		return false
	} else {
		return true
	}
}

func (client *Client) IsAuthenticated() bool {
	if client.token == "" {
		return false
	} else {
		//TODO check token validity
		return true
	}
}

func (client *Client) GameId() int {
	for index, _ := range client.CurrentGame.Clients {
		if client == client.CurrentGame.Clients[index] {
			return index
		}
	}
	return -1
}

func (client *Client) Start() {
	if client.IsAi(){
		go client.Ai.Start()
		return
	}
	var conn = client.Conn
	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("read: ", err)
			return
		}
		if mt == websocket.TextMessage {
			request := Request{}
			json.Unmarshal(message, &request)
			request.Client = client
			if request.Type == "PlacePiece" || request.Type == "PlaceRandom" || request.Type == "Fetch" || request.Type == "FetchPlayer" {
				client.CurrentGame.RequestChannel <- request
			}
		}
	}
}