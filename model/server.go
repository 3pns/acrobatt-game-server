package model

import (
	. "../utils"
	"sync"
	"time"
)

type server struct {
	currentGames   []*Game
	clients []*Client
	clientsInHome  []*Client
	lobbies []*Lobby
}

//thread safe singleton pattern
var instance *server
var once sync.Once

func GetServer() *server {
    once.Do(func() {
        instance = &server{}
        instance.currentGames = []*Game{}
        instance.clients = []*Client{}
        //instance.clientsInHome = []*Client{}
    })
    return instance
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
	request.MarshalData(LobbySlice{serv.lobbies})
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