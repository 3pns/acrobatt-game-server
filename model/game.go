package model

import (
	. "../utils"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
)

type Game struct {
	Id int
	Clients        []*Client
	board          *Board
	RequestChannel chan Request
}

type GameFactory struct {
	Id       int
}

func NewGameFactory() *GameFactory {
	var factory = new(GameFactory)
	factory.Id = 0
	return factory
}

type GameSlice struct {
	Games []*Game `json:"games"`
}

func (factory *GameFactory) NewGame(clients []*Client) *Game {
	var game = Game{factory.Id, clients, nil, make(chan Request, 100)}
	factory.Id++
	return &game
}

func (game *Game) Board() Board {
	return *game.board
}

func StartDemo(client *Client) {
	fmt.Println("Client : Switching to Game State ")
	var client0 = client
	var client1 = NewAiClient()
	var client2 = NewAiClient()
	var client3 = NewAiClient()
	clients := []*Client{client0, client1, client2, client3}
	var game = GetServer().GetGameFactory().NewGame(clients)
	go game.Start()
	fmt.Println("Client : StartDemo GO !!!")
}

func (game *Game) Start() {
	fmt.Println("Starting Game")
	var board Board
	board.InitBoard()
	board.InitPieces()
	board.InitPlayers()
	game.board = &board

	for index, _ := range game.Clients {
		game.Clients[index].CurrentGame = game
		if game.Clients[index].IsAi() {
			game.Clients[index].Ai.Player = game.board.Players[index]
			go game.Clients[index].Start()
		} else {
			game.Clients[index].State.Event("join_game")
		}
	}

	request := Request{}
	for {
		request = <-game.RequestChannel
		fmt.Println("Game[" + string(game.Id) + "]: new request detected")
		player := board.Players[request.Client.GameId()]
		isPlayerTurn := player == board.PlayerTurn
		conn := &websocket.Conn{}
		if !request.Client.IsAi() {
			conn = request.Client.Conn
		}
		if request.Type == "PlacePiece" && isPlayerTurn {
			fmt.Println("Message de type PlacePiece detected !")
			piece := Piece{}
			json.Unmarshal(request.Data, &piece)
			fmt.Println(piece)
			placed := player.PlacePiece(piece, &board, false)
			if placed {
				var req = Request{"PlacementConfirmed", "", nil, request.CallbackId, nil}
				WriteTextMessage(conn, req.Marshal())
				game.board.NextTurn()
				game.board.PrintBoard()
				game.BroadcastRefresh()
			} else {
				fmt.Println("PlacementRefused")
				var req = Request{"PlacementRefused", "", nil, request.CallbackId, nil}
				WriteTextMessage(conn, req.Marshal())
			}
		} else if request.Type == "PlaceRandom" && isPlayerTurn {
			fmt.Println("Message de type PlaceRandom detected !")
			player.PlaceRandomPieceWithIAEasy(&board, false)
			game.board.NextTurn()
			game.board.PrintBoard()
			game.BroadcastRefresh()
		} else if request.Type == "Fetch" {
			var req = Request{"Fetch", "", nil, request.CallbackId, nil}
			req.MarshalData(game.Board())
			WriteTextMessage(conn, req.Marshal())
		} else if request.Type == "FetchPlayer" {
			var req = Request{"FetchPlayer", "Player", nil, request.CallbackId, nil}
			req.MarshalData(player)
			WriteTextMessage(conn, req.Marshal())
		} else if request.Type == "Concede" && isPlayerTurn {
			player.Concede()
			game.BroadcastConcede(player)
			game.board.NextTurn()
			game.BroadcastRefresh()
		} else if request.Type == "Quit" && !player.HasPlaceabePieces(game.board) {
			request.Client.State.Event("quit_demo")
		}
		if game.IsGameOver() {
			game.BroadcastGameOver()
			game.DisconnectPlayers()
			return
		}
	}
	return
}

func (game *Game) BroadcastConcede(player *Player) {
	var req = Request{"Concede", "", nil, "", nil}
	req.MarshalData(*player)
	game.BroadCastRequest(req)
}

func (game *Game) BroadcastRefresh() {
	var req = Request{"Refresh", "", nil, "", nil}
	req.MarshalData(game.Board())
	game.BroadCastRequest(req)
}

func (game *Game) BroadcastGameOver() {
	var req = Request{"GameOver", "", nil, "", nil}
	game.BroadCastRequest(req)
}

func (game *Game) BroadCastRequest(request Request) {
	for index, _ := range game.Clients {
		if game.Clients[index].IsAi() {
			game.Clients[index].Ai.RequestChannel <- request
		} else {
			WriteTextMessage(game.Clients[index].Conn, request.Marshal())
		}
	}
}

func (game *Game) IsGameOver() bool {
	for index, _ := range game.board.Players {
		if game.board.Players[index].HasPlaceabePieces(game.board) {
			return false
		}
	}
	return true
}

func (game *Game) DisconnectPlayers() {
	for index, _ := range game.Clients {
		if !game.Clients[index].IsAi() {
			if game.Clients[index].IsAuthenticated() {
				game.Clients[index].CurrentGame = nil
				game.Clients[index].State.Event("quit_game")
			} else {
				game.Clients[index].CurrentGame = nil
				game.Clients[index].State.Event("quit_demo")
			}
		}
	}
}
