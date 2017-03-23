package model

import (
	. "../utils"
	"sync"
	"time"
	"fmt"
)

type server struct {
	currentGames   []*Game
	clients []*Client
	lobbies  map[int]*Lobby
	lobbyFactory *LobbyFactory
	gameFactory *GameFactory
}

//thread safe singleton pattern
var instance *server
var once sync.Once

func GetServer() *server {
    once.Do(func() {
        instance = &server{}
        instance.currentGames = []*Game{}
        instance.clients = []*Client{}
        instance.lobbies = make(map[int]*Lobby)
        instance.lobbyFactory = NewLobbyFactory()
        instance.gameFactory = NewGameFactory()
    })
    return instance
}

func (serv *server) GetLobbyFactory() *LobbyFactory {
 	return serv.lobbyFactory
}

func (serv *server) GetGameFactory() *GameFactory {
 	return serv.gameFactory
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
	serv.clients = append(serv.clients, client)
}

func (serv *server) RemoveClient(client *Client) {
	for index, _ := range serv.clients {
		if serv.clients[index] == client {
			copy(serv.clients[index:], serv.clients[index+1:])
			serv.clients[len(serv.clients)-1] = nil
			serv.clients = serv.clients[:len(serv.clients)-1]
		}
	}
}

func (serv *server) AddGame(game *Game) {
	serv.currentGames = append(serv.currentGames, game)
}

func (serv *server) RemoveGame(game *Game) {
	for index, _ := range serv.currentGames {
		if serv.currentGames[index] == game {
			copy(serv.currentGames[index:], serv.currentGames[index+1:])
			serv.currentGames[len(serv.currentGames)-1] = nil
			serv.currentGames = serv.currentGames[:len(serv.currentGames)-1]
		}
	}
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