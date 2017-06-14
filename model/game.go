package model

import (
	. "../jsonapi"
	. "../utils"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	_ "github.com/Sirupsen/logrus"

)

type Game struct {
	Id             int             `json:"id"`
	Clients        map[int]*Client `json:"clients"`
	Observers        map[int]*Client `json:"observers"`
	HubClients        map[int]*Client `json:"-"`
	board          *Board          `json:"-"`
	RequestChannel chan Request    `json:"-"`
	Moves map[int]Move `json:"-"`
	hub *Hub `json:"-"`
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
	var game Game 
	game.Id = factory.Id
	game.Clients = make(map[int]*Client)
	game.Observers = make(map[int]*Client)
	game.board = nil
	game.RequestChannel = make(chan Request, 100)
	game.Moves = make(map[int]Move)
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
	game.board.GameId = game.Id
	game.StartHub()

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
	//join game observers
	for index := range game.Observers {
		game.Observers[index].State.Event("join_game")
	}
	game.BroadcastRefresh()
	game.BroadcastClients()
	request := Request{}
	for {
		request = <-game.RequestChannel
		var player *Player
		client := request.Client
		client.UpdateTrace("Game[" + strconv.Itoa(game.Id) + "]->")

		//Pour les Observateurs
		if request.Client.ObserverId() >= 0 {
			if request.Type == "Quit" {
				game.RemoveClient(client)
				client.UpdateTrace("QuitingGame->")
				request.Client.State.Event("quit_game")
				client.PrintTrace()
			}
		}

		// Pour les Joueurs et Observateurs
		if request.Client.GameId() >= 0 || request.Client.ObserverId() >= 0 {
			if request.Type == "Fetch" {
				client.UpdateTrace("Fetch->")
				var req = NewRequestWithCallbackId("Fetch", request.CallbackId)
				req.MarshalData(game.Board())
				client.RequestChannel <- req
			} else if request.Type == "BroadcastMessage" {
				game.hub.RequestChannel <-request
			}
		}

		//Pour les Joeurs
		if request.Client.GameId() >= 0 {
			player = board.Players[request.Client.GameId()]
			isPlayerTurn := player == board.PlayerTurn
			 if request.Type == "FetchPlayer" {
				client.UpdateTrace("FetchPlayer->")
				var req = NewRequestWithCallbackId("FetchPlayer", request.CallbackId)
				req.MarshalData(*player)
				client.RequestChannel <- req
			} else if request.Type == "PlacePiece" && isPlayerTurn {
				client.UpdateTrace("PlacePiece->")
				piece := Piece{}
				json.Unmarshal(request.Data, &piece)
				fmt.Println(piece)
				placedPiece := player.PlacePiece(piece, &board, false)
				if placedPiece != nil {
					var req = NewRequestWithCallbackId("PlacementConfirmed", request.CallbackId)
					client.UpdateTrace("PlacementConfirmed->")
					client.RequestChannel <- req
					game.board.PlayerTurn.Time += game.board.PlayerTurn.GetTurnTime()
					move := Move{game.board.Turn, game.board.PlayerTurn.Id, game.board.PlayerTurn.ApiId(), placedPiece, int(game.board.PlayerTurn.GetTurnTime() / time.Millisecond)}
					game.Moves[game.board.PlayerTurn.Id] = move
					game.board.PlayerTurn.UpdateScore(move)
					game.board.NextTurn()
					//game.board.PrintBoard()
					game.BroadcastRefresh()
				} else {
					var req = NewRequestWithCallbackId("PlacementRefused", request.CallbackId)
					client.UpdateTrace("PlacementRefused->")
					client.RequestChannel <- req
				}
			} else if request.Type == "PlaceRandom" && isPlayerTurn {
				client.UPTrace("PlaceRandom")
				var piece *Piece
				if !client.IsAi() || client.IsAi() && client.Ai.Difficulty == "easy"{
					piece = player.PlaceRandomPieceWithIAEasy(&board, false)
				} else if client.IsAi() && client.Ai.Difficulty == "medium"{
					piece = player.PlaceRandomPieceWithIAMedium(&board, false)
				}
				if piece != nil {
					move := Move{game.board.Turn, game.board.PlayerTurn.Id, game.board.PlayerTurn.ApiId(), piece, int(game.board.PlayerTurn.GetTurnTime() / time.Millisecond)}
					game.Moves[game.board.PlayerTurn.Id] = move
					game.board.PlayerTurn.UpdateScore(move)
				}
				game.board.PlayerTurn.Time += game.board.PlayerTurn.GetTurnTime()
				game.board.NextTurn()
				//game.board.PrintBoard()
				game.BroadcastRefresh()
			} else if request.Type == "Concede" {
				client.UPTrace("Concede")
				player.Concede()
				game.BroadcastConcede(player)
				if isPlayerTurn {
					game.board.PlayerTurn.Time += game.board.PlayerTurn.GetTurnTime()
					game.board.NextTurn()
				}
				game.BroadcastRefresh()
			} else if request.Type == "Quit" && !player.HasPlaceabePieces(game.board) {
				client.UpdateTrace("Quit->")
				if request.Client.IsAuthenticated() {
					client.CurrentGame = nil
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
				game.board.PlayerTurn.Time += board.PlayerTurn.GetTurnTime()
				game.BroadcastGameOver()
				game.PersistGameHistory()
				game.hub.Stop()
				game.DisconnectPlayers()
				game.DisconnectObservers()
				GetServer().RemoveGame(game)
				return
			}
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

func (game *Game) BroadcastClients() {
	var req = NewRequest("BroadcastClients")
	req.MarshalData(game)
	game.BroadCastRequest(req)
}

func (game *Game) BroadCastRequest(request Request) {
	for index, _ := range game.Clients {
		if game.Clients[index].IsAi() {
			game.Clients[index].Ai.RequestChannel <- request
		} else {
			if game.Clients[index].CurrentGame == game{
				game.Clients[index].RequestChannel <- request
			}
		}
	}
	for index, _ := range game.Observers {
		game.Observers[index].RequestChannel <- request
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

func (game *Game) DisconnectObservers() {
	for index, _ := range game.Observers {
				game.Observers[index].State.Event("quit_game")
				delete(game.Observers, index)
	}
}

func (game *Game) JoinAsObserver(client *Client) {
	game.Observers[client.Id] = client
	game.HubClients[client.Id] = client
	client.State.Event("join_game")
	client.CurrentGame = game
	game.BroadcastClients()
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
	if game.Observers[client.Id] != nil {
		delete(game.Observers, client.Id)
	}
	if game.HubClients[client.Id] != nil {
		delete(game.HubClients, client.Id)
	}
	game.BroadcastClients()
}

func (game *Game) PersistGameHistory() {
	//game 
	marshalledData, _ := json.Marshal(game.Moves)
	gj := GameJson{marshalledData}
	marshalledGJ, _ := json.Marshal(gj)
	_, response, _ := ApiRequest("POST", "manager/game", marshalledGJ)

	game_id, _ := strconv.Atoi(fmt.Sprintf("%v", response["id"]))

	for index := range game.board.Players {
		player := game.board.Players[index]
		rank := game.board.GetRankByPlayer(player)
		fmt.Println("inb4 save")
		fmt.Println(player.Time)
		var history HistoryJson
		if player.ApiId() == -1 {
			history = HistoryJson{game_id, -(player.Id + 1), player.Id, player.Score, int(player.Time / time.Millisecond), rank}
		} else {
			history = HistoryJson{game_id, player.ApiId(), player.Id, player.Score, int(player.Time / time.Millisecond), rank}
		}
		
		marshalledHistory, _ := json.Marshal(history)
		ApiRequest("POST", "manager/history", marshalledHistory)
	}
}

func (game *Game) SwapClients(oldClient *Client, newClient *Client) bool {
	for index, _ := range game.Clients {
		if game.Clients[index] == oldClient {
			game.Clients[index] = newClient
		}
	}
	return true
}

func (game *Game) StartHub() {
	game.HubClients = make(map[int]*Client)
	//ajout des joueurs
	for index, _ := range game.Clients {
		if game.Clients[index] != nil {
			game.HubClients[game.Clients[index].Id] = game.Clients[index]
		}
	}
	//ajout des observateurs
	for index, _ := range game.Observers {
		if game.Observers[index] != nil {
			game.HubClients[index] = game.Observers[index]
		}
	}
	hub := GetServer().GetHubFactory().NewHub()
	hub.Clients = game.HubClients
	hub.HolderType = "game"
	hub.HolderId = game.Id
	go hub.Start()
	game.hub = hub
}
