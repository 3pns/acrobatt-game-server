package model

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/looplab/fsm"
)

type Client struct {
	Id       int `json:"id"`
	Conn        *websocket.Conn `json:"-"`
	token       string `json:"-"`
	State       *fsm.FSM `json:"-"`
	Ai          *AI `json:"-"`
	CurrentGame *Game `json:"-"`
	CurrentLobby *Lobby `json:"-"`
}

type ClientFactory struct {
	id       int
}

func NewClientFactory() *ClientFactory {
	var factory = new(ClientFactory)
	factory.id = 0
	return factory
}

func (factory *ClientFactory) NewClient(conn *websocket.Conn) *Client {
	var client Client
	client.Id = factory.id
	factory.id++
	client.Ai = nil
	client.Conn = conn
	client.token = ""

	client.State = fsm.NewFSM(
		"start",
		fsm.Events{
			{Name: "create_demo", Src: []string{"start"}, Dst: "game"},
			{Name: "quit_demo", Src: []string{"game"}, Dst: "start"},
			{Name: "authenticate", Src: []string{"start"}, Dst: "home"},
			{Name: "join_lobby", Src: []string{"home"}, Dst: "lobby"},
			{Name: "create_lobby", Src: []string{"home"}, Dst: "lobby"},
			{Name: "quit_lobby", Src: []string{"lobby"}, Dst: "home"},
			{Name: "join_game", Src: []string{"lobby"}, Dst: "game"},
			{Name: "quit_game", Src: []string{"game"}, Dst: "home"},
		},
		fsm.Callbacks{
			"create_demo": func(e *fsm.Event) { StartDemo(&client) },
			"quit_demo":   func(e *fsm.Event) {
				client.CurrentGame = nil
				fmt.Println("quiting demo : " + e.FSM.Current()) 
				},
			"authenticate":   func(e *fsm.Event) { 
				fmt.Println("authenticating : " + e.FSM.Current())
				GetServer().AddClient(&client)
				},
			"join_lobby":   func(e *fsm.Event) { fmt.Println("joining lobby : " + e.FSM.Current()) },
			"create_lobby":   func(e *fsm.Event) { 
				fmt.Println("creating lobby : " + e.FSM.Current())
				lobby := GetServer().GetLobbyFactory().NewLobby(&client)
				GetServer().AddLobby(lobby)
				},
			"quit_lobby":   func(e *fsm.Event) { 
				client.CurrentLobby = nil
				fmt.Println("quiting lobby : " + e.FSM.Current()) 
				},
			"join_game":   func(e *fsm.Event) { 
				client.CurrentLobby = nil
				fmt.Println("joining game : " + e.FSM.Current()) 
				},
			"quit_game":   func(e *fsm.Event) { 
				client.CurrentGame = nil
				fmt.Println("quiting game : " + e.FSM.Current()) 
				},
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
	for key := range client.CurrentGame.Clients {
		if client.CurrentGame.Clients[key] == client {
			return key
		}
	}
	return -1
}

func (client *Client) Authenticate(token string) bool {
	if token == "" {
		return false
	} else {
		client.token = token
		client.State.Event("authenticate")
		return true
	}
}

func (client *Client) Start() {
	if client.IsAi() {
		go client.Ai.Start()
		return
	}
	var conn = client.Conn
	defer client.Conn.Close()
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
			fmt.Print("Client[",client.Id,"]{" + request.Type +"-"+request.DataType+"}")
			request.Dispatch()
		}
	}
}

func (client *Client) JoinLobby (lobby *Lobby) {
	client.CurrentLobby = lobby
	client.State.Event("join_lobby")
}

func (client *Client) LeaveLobby () {
	client.State.Event("quit_lobby")
}
