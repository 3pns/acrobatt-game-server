package model

import (
	"fmt"
	"github.com/gorilla/websocket"
	. "../utils"
)

type Lobby struct {
	Id int
	Name string
	Clients []*Client 
	AIClients []*Client `json:"-"`
	Master *Client
	game *Game
	RequestChannel chan Request `json:"-"`
}

type LobbyFactory struct {
	Id       int
}

func NewLobbyFactory() *LobbyFactory {
	var factory = new(LobbyFactory)
	factory.Id = 0
	return factory
}

type LobbySlice struct {
	Lobbies []*Lobby `json:"lobbies"`
}

func (factory *LobbyFactory)NewLobby(client *Client) *Lobby {
	var lobby Lobby
	lobby.Id = factory.Id
	factory.Id++
	lobby.Name = "TEST"
	lobby.AIClients = []*Client{NewAiClient(), NewAiClient(), NewAiClient(), NewAiClient()}
	lobby.Clients = []*Client{client, lobby.AIClients[1], lobby.AIClients[2], lobby.AIClients[3]}
	lobby.Master = client
	lobby.game = GetServer().GetGameFactory().NewGame(lobby.Clients)
	lobby.RequestChannel = make(chan Request, 100)
	client.CurrentLobby = &lobby
	go lobby.Start()
	return &lobby
}

func (lobby *Lobby) Start() {
	request := Request{}
	conn := &websocket.Conn{}
	for {
		request = <-lobby.RequestChannel
		fmt.Println("Lobby[" + string(lobby.Id) + "]: new request detected")
		var client = request.Client
		conn = request.Client.Conn
		if request.Type == "Start" && (client == lobby.Master) {
			if len(lobby.game.Clients) == 4 && lobby.game.Clients[0] != nil && lobby.game.Clients[1] != nil && lobby.game.Clients[2] != nil && lobby.game.Clients[3] != nil {
				go lobby.game.Start()

			}
		} else if request.Type == "Fetch" {
			var req = Request{"Fetch", "", nil, request.CallbackId, nil}
			req.MarshalData(lobby)
			WriteTextMessage(conn, req.Marshal())
		}
	}
	
}