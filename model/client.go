package model

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/looplab/fsm"
)

type Client struct {
	Conn        *websocket.Conn
	token       string
	State       *fsm.FSM
	Ai          *AI
	CurrentGame *Game
}

func NewClient(conn *websocket.Conn) *Client {
	var client Client
	client.Conn = conn
	client.token = ""
	client.Ai = nil

	client.State = fsm.NewFSM(
		"start",
		fsm.Events{
			{Name: "create_demo", Src: []string{"start"}, Dst: "game"},
			{Name: "quit_demo", Src: []string{"game"}, Dst: "start"},
		},
		fsm.Callbacks{
			"create_demo": func(e *fsm.Event) { StartDemo(&client) },
			"quit_demo":   func(e *fsm.Event) { fmt.Println("quiting demo : " + e.FSM.Current()) },
		},
	)

	return &client
}

func NewAiClient() *Client {
	var client Client
	client.Conn = nil
	client.token = ""
	client.State = nil
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

func (client *Client) Authenticate(token string) bool {
	if token == "" {
		return false
	} else {
		client.token = token
		return true
	}
}

func (client *Client) Start() {
	if client.IsAi() {
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
			request.Dispatch()
		}
	}
}
