package model

import (
	"fmt"
)

type Lobby struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Clients []*Client `json:"clients"`
	AIClients map[int]*Client`json:"-"`
	Master *Client `json:"master"`
	Seats map[int]*Client `json:"seats"`
	game *Game `json:"-"`
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
	for {
		request = <-lobby.RequestChannel
		fmt.Println("Lobby[",lobby.Id,"]->")
		var client = request.Client
		if request.Type == "Start" && (client == lobby.Master) {
			if lobby.Seats[0] != nil && lobby.Seats[1] != nil && lobby.Seats[2] != nil && lobby.Seats[3] != nil {
		
				for key := range lobby.Seats {
					fmt.Println("clé[",key,"]->settings clients in game and current game")
					lobby.game.Clients[key] = lobby.Seats[key]
					lobby.game.Clients[key].CurrentGame = lobby.game
				}

				go lobby.game.Start()
				lobby.broadcastStart()
				GetServer().RemoveLobby(lobby)
				GetServer().AddGame(lobby.game)
				return
			}
		} else if request.Type == "FetchLobby" {
			var req = NewRequestWithCallbackId ("FetchLobby", request.CallbackId)
			req.MarshalData(*lobby)
			WriteTextMessage(client, req.Marshal())
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
	req := NewRequest ("Broadcast")
	req.MarshalData(*lobby)
	lobby.broadcastRequest(&req)
}

func (lobby *Lobby) broadcastStart() {
	req := NewRequest ("Start")
	lobby.broadcastRequest(&req)
}

func (lobby *Lobby) broadcastRequest(request *Request) {
	for index, _ := range lobby.Clients {
		if lobby.Clients[index].State.Current() == "lobby" {
			WriteTextMessage(lobby.Clients[index], request.Marshal())
		}
	}
}

func (lobby *Lobby) RemoveClient(client *Client) {
	lobby.unsit(client)
	for index, _ := range lobby.Clients {
		if lobby.Clients[index] == client {
			copy(lobby.Clients[index:], lobby.Clients[index+1:])
			lobby.Clients[len(lobby.Clients)-1] = nil
			lobby.Clients = lobby.Clients[:len(lobby.Clients)-1]
		}
	}
	//si il s'agit du Master, un autre client devient Master, sinon on supprime le lobby
	if lobby.isMaster(client) && len(lobby.Clients) > 0 {
		lobby.Master = lobby.Clients[0]
		return
	} else {
		GetServer().RemoveLobby(lobby)
		return
	}
}