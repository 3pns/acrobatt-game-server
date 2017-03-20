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
	lobbies []*Lobby
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
		fmt.Println("Lobbies : ")
		fmt.Println(serv.lobbies)
		fmt.Println(GetServer().lobbies)
		fmt.Println(LobbySlice{GetServer().lobbies})
	}
}

func (serv *server) Process(request Request) {
	var client = request.Client

	if request.Type == "CreateLobby" {
		client.State.Event("create_lobby")
	} else if request.Type == "JoinLobby" {
		//TODO JOIN le lobby envoyé dans la requête
		client.State.Event("join_lobby")
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

func (serv *server) AddLobby(lobby *Lobby) {
	serv.lobbies = append(serv.lobbies, lobby)
}

func (serv *server) broadcastLobbies() {
	request := Request{"Broadcast", "[]Lobby", nil, "", nil}
	lobbyslice := LobbySlice{}
	lobbyslice.Lobbies = serv.lobbies
	fmt.Println(lobbyslice)
	request.MarshalData(lobbyslice)
	serv.broadcastRequest(&request)
}

func (serv *server) broadcastGames() {
	request := Request{"Broadcast", "[]Game", nil, "", nil}
	request.MarshalData(GameSlice{serv.currentGames})
	serv.broadcastRequest(&request)
}

func (serv *server) broadcastRequest(request *Request) {
	for index, _ := range serv.clients {
		if serv.clients[index].State.Current() == "home" {
			WriteTextMessage(serv.clients[index].Conn, request.Marshal())
		}
	}
}