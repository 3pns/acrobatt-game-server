package model

import (
	_ "encoding/json"
	log "github.com/Sirupsen/logrus"
	"sync"
	"time"
)

type server struct {
	clients        map[int]*Client
	currentGames   map[int]*Game
	lobbies        map[int]*Lobby
	lobbyFactory   *LobbyFactory
	gameFactory    *GameFactory
	clientFactory  *ClientFactory
	hubFactory     *HubFactory
	cleanerChannel chan *Client
	hub            *Hub
}

//thread safe singleton pattern
var instance *server
var once sync.Once

func GetServer() *server {
	once.Do(func() {
		instance = &server{}
		instance.clients = make(map[int]*Client)
		instance.currentGames = make(map[int]*Game)
		instance.lobbies = make(map[int]*Lobby)
		instance.lobbyFactory = NewLobbyFactory()
		instance.gameFactory = NewGameFactory()
		instance.clientFactory = NewClientFactory()
		instance.hubFactory = NewHubFactory()
		instance.cleanerChannel = make(chan *Client, 100)
	})
	return instance
}

func (serv *server) GetLobbyFactory() *LobbyFactory {
	return serv.lobbyFactory
}

func (serv *server) GetGameFactory() *GameFactory {
	return serv.gameFactory
}

func (serv *server) GetClientFactory() *ClientFactory {
	return serv.clientFactory
}

func (serv *server) GetHubFactory() *HubFactory {
	return serv.hubFactory
}

func (serv *server) Start() {
	//Démarrage du Chat Hub DU Serveur
	hub := GetServer().GetHubFactory().NewHub()
	hub.Clients = serv.clients
	hub.HolderType = "server"
	hub.HolderId = -1
	go hub.Start()
	serv.hub = hub
	for {
		time.Sleep(5 * time.Second)
		serv.broadcastLobbies()
		serv.broadcastGames()
		serv.broadcastHomeClients()
	}
}

func (serv *server) Process(request Request) {
	var client = request.Client
	if request.Type == "CreateLobby" {
		client.State.Event("create_lobby")
		lobby := GetServer().GetLobbyFactory().NewLobby(client)
		GetServer().AddLobby(lobby)
		client.CurrentLobby = lobby
		client.UPTrace("CreateLobby")
	} else if request.Type == "JoinLobby" {
		index := request.DataToInt()
		if serv.lobbies[index] != nil {
			serv.lobbies[index].Join(client)
			client.UPTrace("JoinLobby")
			var req = NewRequestWithCallbackId("JoinLobbyConfirmed", request.CallbackId)
			client.RequestChannel <- req
		} else {
			var req = NewRequestWithCallbackId("JoinLobbyRefused", request.CallbackId)
			client.RequestChannel <- req
		}
	} else if request.Type == "ObserveGame" {
		index := request.DataToInt()
		if serv.currentGames[index] != nil {
			serv.currentGames[index].JoinAsObserver(client)
		}
		client.UPTrace("ObserveGame")
	} else if request.Type == "BroadcastMessage" {
		serv.hub.RequestChannel <- request
		client.UPTrace("BroadcastMessage")
	}
}

func (serv *server) AddClient(client *Client) {
	serv.clients[client.Id] = client
}

func (serv *server) RemoveClient(client *Client) {
	delete(serv.clients, client.Id)
}

func (serv *server) CleanClient(client *Client) {
	serv.cleanerChannel <- client
}

func (serv *server) StartCleaner() {
	client := &Client{}
	for {
		client = <-serv.cleanerChannel

		if serv.clients[client.Id] != nil && client != nil && serv.clients[client.Id] == client && serv.clients[client.Id].CurrentGame != nil {
			serv.clients[client.Id].CurrentGame.RemoveClient(client)
		}
		if serv.clients[client.Id] != nil && client != nil && serv.clients[client.Id] == client && serv.clients[client.Id].CurrentLobby != nil {
			serv.clients[client.Id].CurrentLobby.RemoveClient(client)
		}
		if serv.clients[client.Id] != nil && client != nil && serv.clients[client.Id] == client {
			delete(serv.clients, client.Id)
		}
		client.Id = -1
	}
	// on vérifie à chaque étape
}

func (serv *server) AddGame(game *Game) {
	serv.currentGames[game.Id] = game
}

func (serv *server) RemoveGame(game *Game) {
	delete(serv.currentGames, game.Id)
}

func (serv *server) AddLobby(lobby *Lobby) {
	serv.lobbies[lobby.Id] = lobby
}

func (serv *server) RemoveLobby(lobby *Lobby) {
	//serv.lobbies[lobby.Id] = nil
	delete(serv.lobbies, lobby.Id)
}

func (serv *server) broadcastLobbies() {
	request := NewRequest("Broadcast")
	lobbiesSlice := LobbySlice{}
	lobbiesSlice.Lobbies = serv.lobbiesSlice()
	log.Info("lobbiesSlice : ")
	log.Info(lobbiesSlice)
	request.MarshalData(lobbiesSlice)
	serv.broadcastRequest(&request)
}

func (serv *server) broadcastGames() {
	request := NewRequest("Broadcast")
	gamesSlice := GameSlice{}
	gamesSlice.Games = serv.gamesSlice()
	log.Info("gamesSlice : ")
	log.Info(gamesSlice)
	request.MarshalData(gamesSlice)
	serv.broadcastRequest(&request)
}

func (serv *server) broadcastHomeClients() {
	request := NewRequest("Broadcast")
	clientsSlice := ClientSlice{}
	clientsSlice.Clients = serv.invitableClientsSlice()
	log.Info("HomeClients : ")
	log.Info(clientsSlice)
	request.MarshalData(clientsSlice)
	serv.broadcastInvitableClientRequest(&request)
}

func (serv *server) sanetizeClients() {
	for key := range serv.clients {
		now := time.Now()
		err := serv.clients[key].Conn.WriteControl(9, []byte("PING"), now.Add(time.Duration(10)*time.Second))
		if err != nil {
			log.Info("Server->sanetize->client[", key, "]")
			serv.CleanClient(serv.clients[key])
		}
	}
}

func (serv *server) broadcastRequest(request *Request) {
	for index, _ := range serv.clients {
		if serv.clients[index].State.Current() == "home" {
			serv.clients[index].RequestChannel <- *request
		}
	}
}

func (serv *server) broadcastInvitableClientRequest(request *Request) {
	for index, _ := range serv.clients {
		if serv.clients[index].State.Current() == "home" || serv.clients[index].State.Current() == "lobby" {
			serv.clients[index].RequestChannel <- *request
		}
	}
}

func (serv *server) lobbiesSlice() []*Lobby {
	lobbyslice := []*Lobby{}
	for _, lobby := range serv.lobbies {
		lobbyslice = append(lobbyslice, lobby)
	}
	return lobbyslice
}

func (serv *server) gamesSlice() []*Game {
	gamesSlices := []*Game{}
	for _, game := range serv.currentGames {
		gamesSlices = append(gamesSlices, game)
	}
	return gamesSlices
}

func (serv *server) invitableClientsSlice() []*Client {
	clientsSlices := []*Client{}
	for _, client := range serv.clients {
		if client.State.Current() == "home" || client.State.Current() == "lobby" {
			clientsSlices = append(clientsSlices, client)
		}
	}
	return clientsSlices
}

func (serv *server) reconnectClient(client *Client) bool {
	if serv.clients[client.Id] != nil {
		serv.clients[client.Id].StopAndSendDCByOtherClientRequest()
		client.MyLobbyId = serv.clients[client.Id].MyLobbyId
		client.MyGameId = serv.clients[client.Id].MyGameId
		client.State = serv.clients[client.Id].State
		client.CurrentGame = serv.clients[client.Id].CurrentGame
		client.CurrentLobby = serv.clients[client.Id].CurrentLobby
		client.RequestChannel = serv.clients[client.Id].RequestChannel
		serv.SwapClients(serv.clients[client.Id], client)
		return true
	} else {
		return false
	}
}

func (serv *server) SwapClients(oldClient *Client, newClient *Client) bool {
	if serv.clients[newClient.Id] != nil {
		serv.clients[newClient.Id] = newClient
	}
	if newClient.CurrentLobby != nil {
		newClient.CurrentLobby.SwapClients(oldClient, newClient)
	}
	if newClient.CurrentGame != nil {
		newClient.CurrentGame.SwapClients(oldClient, newClient)
	}
	return true
}
