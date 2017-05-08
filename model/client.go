package model

import (
	"bytes"
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	"github.com/looplab/fsm"
	"strconv"
	. "../jsonapi"
	. "../utils"
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
	trace          string          `json:"-"`
}

type ClientFactory struct {
	id int
}

func NewClientFactory() *ClientFactory {
	var factory = new(ClientFactory)
	factory.id = 1
	return factory
}

func (factory *ClientFactory) NewClient(conn *websocket.Conn) *Client {
	var client Client
	client.Id = -1//factory.id
	//factory.id++
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
				client.UpdateTrace("quiting demo : " + e.FSM.Current())
			},
			"authenticate": func(e *fsm.Event) {
				client.UpdateTrace("authenticating : " + e.FSM.Current())
				GetServer().AddClient(&client)
			},
			"join_lobby": func(e *fsm.Event) { client.UpdateTrace("joining lobby : " + e.FSM.Current()) },
			"create_lobby": func(e *fsm.Event) {
				client.UpdateTrace("creating lobby : " + e.FSM.Current())
				lobby := GetServer().GetLobbyFactory().NewLobby(&client)
				GetServer().AddLobby(lobby)
			},
			"quit_lobby": func(e *fsm.Event) {
				client.CurrentLobby = nil
				client.UpdateTrace("quiting lobby : " + e.FSM.Current())
			},
			"join_game": func(e *fsm.Event) {
				client.CurrentLobby = nil
				client.UpdateTrace("joining game : " + e.FSM.Current())
			},
			"quit_game": func(e *fsm.Event) {
				client.CurrentGame = nil
				client.UpdateTrace("quiting game : " + e.FSM.Current())
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
func (client *Client) Trace() string {
	return client.trace
}
func (client *Client) Tracing() bool {
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

func (client *Client) Authenticate(auth AuthenticateJson) bool{
	marshalledAuth, err := json.Marshal(auth)
    if err != nil {
        log.Error("%s\n", err)
        return false
    }
	resp, _, err := ApiRequest("POST", "manager/authenticate_player", marshalledAuth)
    if err != nil {
        return false
    }
	if resp.StatusCode == 200 {
		client.Id = auth.PlayerId
		client.State.Event("authenticate")
		client.UPTrace("success")
		return true
	} else {
		client.UPTrace("failure:"+strconv.Itoa(auth.PlayerId) + ":" + auth.AccessToken + ":" + auth.Client)
		return false
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
			log.Warn("read: ", err)
			//TODO on atterit ici et sa affiche websocket: close 1005 (no status)  lorsqu'un utilisateur ferme la fenetre ou a temporairement plus de raison
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
				log.Info("Server->Writer->", err.Error(), "->Client is being removed from Server")
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
