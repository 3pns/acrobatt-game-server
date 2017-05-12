package model

import (
	. "../jsonapi"
	. "../utils"
	"encoding/json"
	"fmt"
	"strconv"
)

type Game struct {
	Id             int             `json:"id"`
	Clients        map[int]*Client `json:"clients"`
	board          *Board          `json:"-"`
	RequestChannel chan Request    `json:"-"`
}

type GameFactory struct {
	Id int
}

func NewGameFactory() *GameFactory {
	var factory = new(GameFactory)
	factory.Id = 1
	return factory
}

type GameSlice struct {
	Games []*Game `json:"games"`
}

func (factory *GameFactory) NewGame() *Game {
	var game = Game{factory.Id, make(map[int]*Client), nil, make(chan Request, 100)}
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
	var game = GetServer().GetGameFactory().NewGame()
	game.Clients[0] = client0
	game.Clients[1] = client1
	game.Clients[2] = client2
	game.Clients[3] = client3
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

	for index := range game.Clients {
		game.Clients[index].CurrentGame = game
		if game.Clients[index].IsAi() {
			game.Clients[index].Ai.Player = game.board.Players[index]
			go game.Clients[index].Start()
		} else {
			game.board.Players[index].SetApiId(game.Clients[index].Id)
			game.Clients[index].State.Event("join_game")
		}
	}
	game.BroadcastRefresh()
	request := Request{}
	for {
		request = <-game.RequestChannel
		//fmt.Println(request)
		//fmt.Println(request.Client)
		player := board.Players[request.Client.GameId()]
		client := request.Client
		client.UpdateTrace("Game[" + string(game.Id) + "]->")
		isPlayerTurn := player == board.PlayerTurn
		if request.Type == "PlacePiece" && isPlayerTurn {
			client.UpdateTrace("PlacePiece->")
			piece := Piece{}
			json.Unmarshal(request.Data, &piece)
			fmt.Println(piece)
			placed := player.PlacePiece(piece, &board, false)
			if placed {
				var req = NewRequestWithCallbackId("PlacementConfirmed", request.CallbackId)
				client.UpdateTrace("PlacementConfirmed->")
				client.RequestChannel <- req
				game.board.NextTurn()
				game.board.PrintBoard()
				game.BroadcastRefresh()
			} else {
				var req = NewRequestWithCallbackId("PlacementRefused", request.CallbackId)
				client.UpdateTrace("PlacementRefused->")
				client.RequestChannel <- req
			}
		} else if request.Type == "PlaceRandom" && isPlayerTurn {
			client.UPTrace("PlaceRandom")
			player.PlaceRandomPieceWithIAEasy(&board, false)
			game.board.NextTurn()
			game.board.PrintBoard()
			game.BroadcastRefresh()
		} else if request.Type == "Fetch" {
			client.UpdateTrace("Fetch->")
			var req = NewRequestWithCallbackId("Fetch", request.CallbackId)
			req.MarshalData(game.Board())
			client.RequestChannel <- req
		} else if request.Type == "FetchPlayer" {
			client.UpdateTrace("FetchPlayer->")
			var req = NewRequestWithCallbackId("FetchPlayer", request.CallbackId)
			req.MarshalData(*player)
			client.RequestChannel <- req
		} else if request.Type == "Concede" && isPlayerTurn {
			client.UPTrace("Concede")
			player.Concede()
			game.BroadcastConcede(player)
			game.board.NextTurn()
			game.BroadcastRefresh()
		} else if request.Type == "Quit" && !player.HasPlaceabePieces(game.board) {
			client.UpdateTrace("Quit->")
			if request.Client.IsAuthenticated() {
				client.UpdateTrace("QuitingGame->")
				request.Client.State.Event("quit_game")
				client.PrintTrace()
			} else {
				client.UpdateTrace("QuitingDemo->")
				request.Client.State.Event("quit_demo")
				client.PrintTrace()
			}
		}
		if game.IsGameOver() {
			client.UPTrace("GameOverDetected")
			game.BroadcastGameOver()
			game.PersistGameHistory()
			game.DisconnectPlayers()
			GetServer().RemoveGame(game)
			return
		}
	}
	return
}

func (game *Game) BroadcastConcede(player *Player) {
	var req = NewRequest("Concede")
	req.MarshalData(*player)
	game.BroadCastRequest(req)
}

func (game *Game) BroadcastRefresh() {
	var req = NewRequest("Refresh")
	req.MarshalData(game.Board())
	game.BroadCastRequest(req)
}

func (game *Game) BroadcastGameOver() {
	var req = NewRequest("GameOver")
	game.BroadCastRequest(req)
}

func (game *Game) BroadCastRequest(request Request) {
	for index, _ := range game.Clients {
		if game.Clients[index].IsAi() {
			game.Clients[index].Ai.RequestChannel <- request
		} else {
			game.Clients[index].RequestChannel <- request
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
				game.Clients[index].State.Event("quit_game")
			} else {
				game.Clients[index].State.Event("quit_demo")
			}
		}
	}
}

func (game *Game) RemoveClient(client *Client) {
	for index, _ := range game.Clients {
		if game.Clients[index] == client {
			game.board.Players[index].Concede()
			delete(game.Clients, index)
			game.BroadcastConcede(game.board.Players[index])
			if game.board.PlayerTurn == game.board.Players[index] {
				game.board.NextTurn()
			}
			game.BroadcastRefresh()
		}
	}
}

func (game *Game) PersistGameHistory() {
	//TODO : remplacer bankdata par une list de struct representant un tour/ un coup
	blankdata := map[string]string{"data": "osef", "asd": "osef", "asdssed": "osef", "asded": "osef", "adedsd": "osef", "aeddsd": "osef", "adesd": "osef"}
	marshalledData, _ := json.Marshal(blankdata)

	//game example
	gj := GameJson{marshalledData}
	marshalledGJ, _ := json.Marshal(gj)
	_, response, _ := ApiRequest("POST", "manager/game", marshalledGJ)

	game_id, _ := strconv.Atoi(fmt.Sprintf("%v", response["id"]))

	//history example
	//TODO calculer le vrai tmps et score de chaque joueur
	for index := range game.board.Players {
		player := game.board.Players[index]
		if player.ApiId() != -1 {
			history := HistoryJson{game_id, player.ApiId(), index, player.Score, player.Time, 7}
			marshalledHistory, _ := json.Marshal(history)
			ApiRequest("POST", "manager/history", marshalledHistory)
		}
	}
}
