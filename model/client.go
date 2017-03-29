package model

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/looplab/fsm"
	log "github.com/Sirupsen/logrus"
	"strconv"
	"bytes"
)

type Client struct {
	Id             int             `json:"id"`
	Conn           *websocket.Conn `json:"-"`
	token          string          `json:"-"`
	State          *fsm.FSM        `json:"-"`
	Ai             *AI             `json:"-"`
	CurrentGame    *Game           `json:"-"`
	CurrentLobby   *Lobby          `json:"-"`
	RequestChannel chan Request    `json:"-"`
	trace string						   `json:"-"`
}

type ClientFactory struct {
	id int
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
	client.RequestChannel = make(chan Request, 100)

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
			"quit_demo": func(e *fsm.Event) {
				client.CurrentGame = nil
				fmt.Println("quiting demo : " + e.FSM.Current())
			},
			"authenticate": func(e *fsm.Event) {
				fmt.Println("authenticating : " + e.FSM.Current())
				GetServer().AddClient(&client)
			},
			"join_lobby": func(e *fsm.Event) { fmt.Println("joining lobby : " + e.FSM.Current()) },
			"create_lobby": func(e *fsm.Event) {
				fmt.Println("creating lobby : " + e.FSM.Current())
				lobby := GetServer().GetLobbyFactory().NewLobby(&client)
				GetServer().AddLobby(lobby)
			},
			"quit_lobby": func(e *fsm.Event) {
				client.CurrentLobby = nil
				fmt.Println("quiting lobby : " + e.FSM.Current())
			},
			"join_game": func(e *fsm.Event) {
				client.CurrentLobby = nil
				fmt.Println("joining game : " + e.FSM.Current())
			},
			"quit_game": func(e *fsm.Event) {
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

func (client *Client) StartTrace() {
	var buffer bytes.Buffer
	buffer.WriteString("Client[")
	buffer.WriteString(strconv.Itoa(client.Id))
	buffer.WriteString("]")
	client.trace = buffer.String()
}
func (client *Client) UpdateTrace(trace string) {
	var buffer bytes.Buffer
	buffer.WriteString(client.trace)
	buffer.WriteString(trace)
	client.trace = buffer.String()
}
func (client *Client) PrintTrace() {
	log.Info(client.trace)
	client.trace = ""
}

func (client *Client) UPTrace(trace string) {
	client.UpdateTrace(trace)
	client.PrintTrace()
}
func (client *Client) Trace() string{
	return client.trace
}
func (client *Client) Tracing() bool{
	if client.trace != "" {
		return true
	}
	return false
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
		client.UPTrace(token)
		return false
	} else {
		client.token = token
		client.State.Event("authenticate")
		client.UPTrace("success")
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
			client.StartTrace()
			request.Dispatch()
		}
	}
}

func (client *Client) StartWriter() {
	request := Request{}
	for {
		request = <-client.RequestChannel
		if client.Tracing() {
			client.UpdateTrace("->Writer->Sending")
		}
		err := client.Conn.WriteMessage(websocket.TextMessage, request.Marshal())
		if err != nil {
			if client.Tracing() {
				client.UpdateTrace("->")
				client.UpdateTrace(err.Error())
				client.UpdateTrace("->Client is being removed from Server")
				client.PrintTrace()
			} else {
				log.Info("Server->Writer->",err.Error(),"->Client is being removed from Server")
			}
			GetServer().CleanClient(client)
			return
		}
		if client.Tracing() {
			client.UpdateTrace("->Sent")
			client.PrintTrace()
		}

	}
}

func (client *Client) JoinLobby(lobby *Lobby) {
	client.CurrentLobby = lobby
	client.State.Event("join_lobby")
}

func (client *Client) LeaveLobby() {
	client.State.Event("quit_lobby")
}
