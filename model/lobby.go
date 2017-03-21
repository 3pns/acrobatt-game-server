package model

import (
	"fmt"
	"github.com/gorilla/websocket"
	. "../utils"
)

type Lobby struct {
	Id int
	Name string
	Clients []*Client `json:"clients"`
	AIClients map[int]*Client`json:"-"`
	Master *Client
	Seats map[int]*Client `json:"seats"`
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
	lobby.AIClients = make(map[int]*Client)
	lobby.AIClients[0] = NewAiClient()
	lobby.AIClients[1] = NewAiClient()
	lobby.AIClients[2] = NewAiClient()
	lobby.AIClients[3] = NewAiClient()
	lobby.Clients = []*Client{client}
	lobby.Master = client
	lobby.game = GetServer().GetGameFactory().NewGame()
	lobby.RequestChannel = make(chan Request, 100)
	lobby.Seats = make(map[int]*Client)
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
			if lobby.Seats[0] != nil && lobby.Seats[1] != nil && lobby.Seats[2] != nil && lobby.Seats[3] != nil {
				lobby.game.Clients = []*Client{lobby.Seats[0], lobby.Seats[1], lobby.Seats[2], lobby.Seats[3]}
				go lobby.game.Start()
				GetServer().RemoveLobby(lobby)
				GetServer().AddGame(lobby.game)
				return
			}
		} else if request.Type == "Fetch" {
			var req = Request{"Fetch", "", nil, request.CallbackId, nil}
			req.MarshalData(lobby)
			WriteTextMessage(conn, req.Marshal())
		}  else if request.Type == "Sit" {
			seatNumber := request.DataToInt()
			if lobby.Seats[seatNumber] == nil {
				lobby.Seats[seatNumber] = client
				lobby.broadcast()
			}
		}  else if request.Type == "Unsit" {
			if lobby.unsit(client){
				lobby.broadcast()
			}
		}  else if request.Type == "SitAI" && lobby.isMaster(client) {
			seatNumber := request.DataToInt()
			if lobby.Seats[seatNumber] == nil {
				lobby.Seats[seatNumber] = lobby.AIClients[seatNumber]
				lobby.broadcast()
			}
			
		}  else if request.Type == "UnsitAI" && lobby.isMaster(client) {
			seatNumber := request.DataToInt()
			if lobby.Seats[seatNumber].IsAi() {
				lobby.Seats[seatNumber] = nil
				lobby.broadcast()
			}
			
		}  else if request.Type == "Quit" {
			lobby.unsit(client)
			lobby.RemoveClient(client)
			client.State.Event("quit_lobby")
			if len(lobby.Clients) == 0 {
				GetServer().RemoveLobby(lobby)
				return
			}
		}
	}
}

func (lobby *Lobby) isMaster(client *Client) bool {
	if lobby.Master == client {
		return true
	}
	return false
}

func (lobby *Lobby) Join(client *Client) {
	lobby.Clients = append(lobby.Clients, client)
	client.CurrentLobby = lobby
	client.State.Event("join_lobby")
}

func (lobby *Lobby) unsit(client *Client) bool {
	for index, _ := range lobby.Seats {
		if lobby.Seats[index] == client {
			lobby.Seats[index] = nil
			return true
		}
	}
	return false
}

func (lobby *Lobby) broadcast() {
	var req = Request{"broadcast", "Lobby", nil, "", nil}
	req.MarshalData(lobby)
	lobby.broadcastRequest(&req)
}

func (lobby *Lobby) broadcastRequest(request *Request) {
	for index, _ := range lobby.Clients {
		if lobby.Clients[index].State.Current() == "lobby" {
			WriteTextMessage(lobby.Clients[index].Conn, request.Marshal())
		}
	}
}

func (lobby *Lobby) RemoveClient(client *Client) {
	for index, _ := range lobby.Clients {
		if lobby.Clients[index] == client {
			copy(lobby.Clients[index:], lobby.Clients[index+1:])
			lobby.Clients[len(lobby.Clients)-1] = nil
			lobby.Clients = lobby.Clients[:len(lobby.Clients)-1]
		}
	}
}