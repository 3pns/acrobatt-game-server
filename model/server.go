package model

import (
	. "../utils"
	"sync"
	"time"
	"fmt"
)

type server struct {
	clients map[int]*Client
  currentGames   map[int]*Game
	lobbies  map[int]*Lobby
	lobbyFactory *LobbyFactory
	gameFactory *GameFactory
  clientFactory *ClientFactory
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

func (serv *server) Start() {
	for {
		time.Sleep(5 * time.Second)
		serv.broadcastLobbies()
		serv.broadcastGames()
	}
}

func (serv *server) Process(request Request) {
	var client = request.Client

	if request.Type == "CreateLobby" {
		client.State.Event("create_lobby")
	} else if request.Type == "JoinLobby" {
		index := request.DataToInt()
		if serv.lobbies[index] != nil {
			serv.lobbies[index].Join(client)
		}
	}
}

func (serv *server) AddClient(client *Client) {
  serv.clients[client.Id] = client
}

func (serv *server) RemoveClient(client *Client) {
  delete(serv.clients, client.Id)
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
	request := Request{"Broadcast", "ListLobby", nil, "", nil}
	lobbiesSlice := LobbySlice{}
	lobbiesSlice.Lobbies = serv.lobbiesSlice()
	fmt.Println(lobbiesSlice)
	request.MarshalData(lobbiesSlice)
	serv.broadcastRequest(&request)
}

func (serv *server) broadcastGames() {
	request := Request{"Broadcast", "ListGame", nil, "", nil}
	gamesSlice := GameSlice{}
	gamesSlice.Games = serv.gamesSlice()
	fmt.Println(gamesSlice)
	request.MarshalData(gamesSlice)
	serv.broadcastRequest(&request)
}

func (serv *server) broadcastRequest(request *Request) {
	for index, _ := range serv.clients {
		if serv.clients[index].State.Current() == "home" {
			WriteTextMessage(serv.clients[index].Conn, request.Marshal())
		}
	}
}

func (serv *server) lobbiesSlice()[]*Lobby{
	lobbyslice := []*Lobby{}
  for _,lobby := range serv.lobbies {
    lobbyslice = append(lobbyslice, lobby)
  }
  return lobbyslice
}

func (serv *server) gamesSlice()[]*Game{
	gamesSlices := []*Game{}
  for _,game := range serv.currentGames {
    gamesSlices = append(gamesSlices, game)
  }
  return gamesSlices
}