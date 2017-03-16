package model

import (
	. "../utils"
	"encoding/json"
	_ "flag"
	"fmt"
	"github.com/gorilla/websocket"
	_ "io"
	_ "net/http"
	_ "strings"
)

type Game struct {
	Clients        []*Client
	board          *Board
	RequestChannel chan Request
}

func NewGame(clients []*Client) Game {
	var game = Game{clients, nil, make(chan Request, 100)}
	return game
}

func (game *Game) Board() Board {
	return *game.board
}

func (game Game) Start() {
	go startGame(game)
}

func startGame(game Game) {
	fmt.Println("Starting Game")
	var board Board
	board.InitBoard()
	board.InitPieces()
	board.InitPlayers()
	game.board = &board

	for index, _ := range game.Clients {
		game.Clients[index].CurrentGame = &game
		if game.Clients[index].IsAi() {
			game.Clients[index].Ai.Player = game.board.Players[index]
		}
	}

	request := Request{}
	for {
		request = <-game.RequestChannel
		player := board.Players[request.Client.GameId()]
		isPlayerTurn := player == board.PlayerTurn
		conn := &websocket.Conn{}
		if !request.Client.IsAi(){
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
			req.MarshalData(game.Board().Players[request.Client.GameId()])
			WriteTextMessage(conn, req.Marshal())
		}
		if game.IsGameOver() {
			game.BroadcastGameOver()
			return
		}
	}
	return
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
