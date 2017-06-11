package model

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"strconv"
)

type Lobby struct {
	Id             int             `json:"id"`
	Name           string          `json:"name"`
	Clients        map[int]*Client      `json:"clients"`
	AIClients      map[int]*Client `json:"-"`
	Master         *Client         `json:"master"`
	Seats          map[int]*Client `json:"seats"`
	game           *Game           `json:"-"`
	RequestChannel chan Request    `json:"-"`
	hub *Hub `json:"-"`
}

type LobbyFactory struct {
	Id int
}

func NewLobbyFactory() *LobbyFactory {
	var factory = new(LobbyFactory)
	factory.Id = 1
	return factory
}

type LobbySlice struct {
	Lobbies []*Lobby `json:"lobbies"`
}

func (factory *LobbyFactory) NewLobby(client *Client) *Lobby {
	var lobby Lobby
	lobby.Id = factory.Id
	factory.Id++
	lobby.Name = "TEST"
	lobby.AIClients = make(map[int]*Client)
	lobby.AIClients[0] = NewAiClient()
	lobby.AIClients[1] = NewAiClient()
	lobby.AIClients[2] = NewAiClient()
	lobby.AIClients[3] = NewAiClient()
	lobby.Clients = make(map[int]*Client)
	lobby.Clients[client.Id] = client //append(lobby.Clients, client)
	lobby.Master = client
	lobby.game = GetServer().GetGameFactory().NewGame()
	lobby.RequestChannel = make(chan Request, 100)
	lobby.Seats = make(map[int]*Client)
	go lobby.Start()
	return &lobby
}

func (lobby *Lobby) Start() {
	hub := GetServer().GetHubFactory().NewHub()
	hub.Clients = lobby.Clients
	hub.HolderType = "lobby"
	hub.HolderId = lobby.Id
	go hub.Start()
	lobby.hub = hub
	for {
		request, more := <-lobby.RequestChannel
		if more {
			var client = request.Client
			if client != nil {
				client.UpdateTrace("Lobby[" + strconv.Itoa(lobby.Id) + "]->")
			}
			if request.Type == "BroadcastMessage" {
					lobby.hub.RequestChannel <-request
			} else if request.Type == "Start" && (client == lobby.Master) {
				if lobby.Seats[0] != nil && lobby.Seats[1] != nil && lobby.Seats[2] != nil && lobby.Seats[3] != nil {

					//Ajout des players
					for key := range lobby.Seats {
						lobby.game.Clients[key] = lobby.Seats[key]
						lobby.game.Clients[key].CurrentGame = lobby.game
					}

					//Ajout des observateurs
					for key := range lobby.Clients {
						if lobby.ClientIsObservater(lobby.Clients[key]){
							lobby.game.Observers[key] = lobby.Clients[key]
						}
					}
										client.UPTrace("Start")
					lobby.broadcastStart()
					go lobby.game.Start()
					GetServer().RemoveLobby(lobby)
					lobby.hub.Stop()
					GetServer().AddGame(lobby.game)
					return
				}
			} else if request.Type == "FetchLobby" {
				client.UPTrace("FetchLobby")
				var req = NewRequestWithCallbackId("FetchLobby", request.CallbackId)
				req.MarshalData(*lobby)
				client.RequestChannel <- req
			} else if request.Type == "Invitation" {
				client.UpdateTrace("Invitation->")
				data := map[string]Client{}
				if err := json.Unmarshal(request.Data, &data); err != nil {
					client.UPTrace("UnmarshallError")
					log.Error(err)
					return
				}
				if GetServer().clients[data["recipient"].Id] != nil {
					client.UPTrace("Sent")
					GetServer().clients[data["recipient"].Id].RequestChannel <- request
				} else {
					client.UpdateTrace("InvitationFailed->")
					request.Type = "InvitationFailed"
					client.RequestChannel <- request
				}
			} else if request.Type == "Sit" {
				client.UpdateTrace("Sit")
				seatNumber := request.DataToInt()
				if lobby.Seats[seatNumber] == nil {
					client.UPTrace("->success")
					lobby.unsit(client)
					lobby.Seats[seatNumber] = client
					lobby.broadcast()
				} else {
					client.UPTrace("->failure")
				}
			} else if request.Type == "Unsit" {
				client.UpdateTrace("Unsit")
				if lobby.unsit(client) {
					client.UPTrace("->success")
					lobby.broadcast()
				} else {
					client.UPTrace("->failure")
				}
			} else if request.Type == "SitAI" && lobby.isMaster(client) {
				client.UpdateTrace("SitAI")
				seatNumber := request.DataToInt()
				if lobby.Seats[seatNumber] == nil {
					lobby.Seats[seatNumber] = lobby.AIClients[seatNumber]
					client.UPTrace("->success")
					lobby.broadcast()
				} else {
					client.UPTrace("->failure")
				}

			} else if request.Type == "UnsitAI" && lobby.isMaster(client) {
				client.UpdateTrace("UnsitAI")
				seatNumber := request.DataToInt()
				if lobby.Seats[seatNumber] != nil && lobby.Seats[seatNumber].IsAi() {
					client.UPTrace("->success")
					lobby.Seats[seatNumber] = nil
					lobby.broadcast()
				} else {
					client.UPTrace("->failure")
				}

			} else if request.Type == "Quit" {
				client.UPTrace("->Quit")
				lobby.unsit(client)
				client.State.Event("quit_lobby")
				lobby.RemoveClient(client)
			} else {
				client.PrintTrace()
			}
		} else {
			log.Info("Destroy Lobby[" + strconv.Itoa(lobby.Id) + "] RequestChannel")
			return
		}
		/*
		fmt.Println("#################### PRINTING SEATS ##########################")
		fmt.Println("seat 0 : ")
		fmt.Println(lobby.Seats[0])
		fmt.Println("seat 1 : ")
		fmt.Println(lobby.Seats[1])
		fmt.Println("seat 2 : ")
		fmt.Println(lobby.Seats[2])
		fmt.Println("seat 3 : ")
		fmt.Println(lobby.Seats[3])*/
	}
}

func (lobby *Lobby) isMaster(client *Client) bool {
	if lobby.Master == client {
		return true
	}
	return false
}

func (lobby *Lobby) Join(client *Client) {
	lobby.Clients[client.Id] = client
	client.CurrentLobby = lobby
	client.State.Event("join_lobby")
	lobby.broadcast()
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
	req := NewRequest("Broadcast")
	req.MarshalData(*lobby)
	lobby.broadcastRequest(&req)
}

func (lobby *Lobby) broadcastStart() {
	req := NewRequest("Start")
	lobby.broadcastRequest(&req)
}

func (lobby *Lobby) broadcastRequest(request *Request) {
	for index, _ := range lobby.Clients {
		//if lobby.Clients[index].State.Current() == "lobby" {
		lobby.Clients[index].RequestChannel <- *request
		//}
	}
}

func (lobby *Lobby) RemoveClient(client *Client) {
	lobby.unsit(client)
	delete(lobby.Clients, client.Id)
	//si il s'agit du Master, un autre client devient Master, sinon on supprime le lobby
	if lobby.isMaster(client) && len(lobby.Clients) > 0 {
		lobby.Master = lobby.Clients[0]
	}
	if len(lobby.Clients) > 0 {
		lobby.broadcast()
	} else {
		lobby.Destroy()
	}
	return
}

func (lobby *Lobby) Destroy() {
	GetServer().RemoveLobby(lobby)
	if lobby.RequestChannel != nil {
		close(lobby.RequestChannel)
		lobby.RequestChannel = nil
		lobby.hub.Stop()
	}

}

func (lobby *Lobby) SwapClients(oldClient *Client, newClient *Client) bool {

	// swap seats
	for index, _ := range lobby.Seats {
		if lobby.Seats[index] == oldClient {
			lobby.Seats[index] = newClient
		}
	}

	//swap list des clients
	for index, _ := range lobby.Clients {
		if len(lobby.Clients) > index && lobby.Clients[index] == oldClient {
			lobby.Clients[index] = newClient
		}
	}

	//swap Master
	if lobby.isMaster(oldClient) {
		lobby.Master = newClient
	}

	return true
}

func (lobby *Lobby) ClientIsObservater(client *Client) bool {
	for key := range lobby.Seats {
		if lobby.Seats[key] == client {
			return false
		}
	}
	return true
}
