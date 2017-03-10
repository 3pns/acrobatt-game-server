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
	client1 *Client
	client2 *Client
	client3 *Client
	client4 *Client
	board   *Board
}

func NewGame(client1 *Client, client2 *Client, client3 *Client, client4 *Client) *Game {
	var game = Game{client1, client2, client3, client4, nil}
	client1.CurrentGame = &game
	client2.CurrentGame = &game
	client3.CurrentGame = &game
	client4.CurrentGame = &game
	return &game
}

func (game *Game) Board() Board{
	return *game.board
}

func (game *Game) Start() {
	go startGame(game)
}

func startGame(game *Game) {
	var board Board
	board.InitBoard()
	board.InitPieces()
	board.InitPlayers()
	game.board = &board
	var turn = 1
	var player *Player = board.Players[0]

	//TODO envoi de la board à la connexion aux 4 joueurs
	var conn = game.client1.Conn

	var req = Request{"Fetch", "", nil, ""}
	req.MarshalData(board)
	WriteTextMessage(conn, req.Marshal())
	var placed = [4]bool{true, true, true, true}
	for {
		if turn > 21 || (!placed[0] && !placed[1] && !placed[2] && !placed[3]) {
			fmt.Println("GAME OVER AT THE BEGGINING")
			refreshBoard(conn, board)
			gameOver(conn)
			board.PrintBoard()
			return
		}

		if player.HasPlaceabePieces(&board) {
			mt, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("read: ", err)
				return
			}
			if mt == websocket.TextMessage {
				request := Request{}
				json.Unmarshal(message, &request)
				if request.Type == "PlacePiece" {
					fmt.Println("Message de type PlacePiece detected !")
					piece := Piece{}
					json.Unmarshal(request.Data, &piece)
					fmt.Println("plaçage de Piece")
					fmt.Println(piece.String())
					placed[0] = player.PlacePiece(piece, &board, false)
					if placed[0] {
						var req = Request{"PlacementConfirmed", "", nil, request.CallbackId}

						placed[1] = board.Players[(player.Id+1)%4].PlaceRandomPieceWithIAEasy(&board, false)
						placed[2] = board.Players[(player.Id+2)%4].PlaceRandomPieceWithIAEasy(&board, false)
						placed[3] = board.Players[(player.Id+3)%4].PlaceRandomPieceWithIAEasy(&board, false)
						WriteTextMessage(conn, req.Marshal())
						//TODO A REFACTORISER AVEC VRAI ARCHI
						refreshBoard(conn, board)
					} else {
						var req = Request{"PlacementRefused", "", nil, request.CallbackId}
						WriteTextMessage(conn, req.Marshal())
					}

				} else if request.Type == "PlaceRandom" {
					fmt.Println("Message de type PlaceRandom detected !")
					//TODO A REFACTORISER AVEC VRAI ARCHI
					placed[0] = player.PlaceRandomPieceWithIAEasy(&board, false)
					placed[1] = board.Players[(player.Id+1)%4].PlaceRandomPieceWithIAEasy(&board, false)
					placed[2] = board.Players[(player.Id+2)%4].PlaceRandomPieceWithIAEasy(&board, false)
					placed[3] = board.Players[(player.Id+3)%4].PlaceRandomPieceWithIAEasy(&board, false)

					var req = Request{"PlacementConfirmed", "", nil, request.CallbackId}
					WriteTextMessage(conn, req.Marshal())
					refreshBoard(conn, board)
					fmt.Println("fin du tour numéro : %i", turn)
					if !placed[0] && !placed[1] && !placed[2] && !placed[3] {
						fmt.Println("GAME OVER CUZ NO MORE PLACEABLE PIECE!!!")
						gameOver(conn)
						return
					}

				} else if request.Type == "Fetch" {
					var req = Request{"Fetch", "", nil, request.CallbackId}
					req.MarshalData(board)
					WriteTextMessage(conn, req.Marshal())
				} else if request.Type == "FetchPlayer" {
					var req = Request{"FetchPlayer", "Player", nil, request.CallbackId}
					req.MarshalData(*player)
					WriteTextMessage(conn, req.Marshal())
				}

			}
		} else {
			//skip turn
			fmt.Println("SKIPPING TURN")
			placed[0] = false
			placed[1] = board.Players[(player.Id+1)%4].PlaceRandomPieceWithIAEasy(&board, false)
			placed[2] = board.Players[(player.Id+2)%4].PlaceRandomPieceWithIAEasy(&board, false)
			placed[3] = board.Players[(player.Id+3)%4].PlaceRandomPieceWithIAEasy(&board, false)
			if !placed[0] && !placed[1] && !placed[2] && !placed[3] {
				fmt.Println("GAME OVER CUZ NO MORE PLACEABLE PIECE!!!")
				refreshBoard(conn, board)
				gameOver(conn)
				board.PrintBoard()
				return
			}
		}
		turn++
		board.PrintBoard()
		placed[0] = true
		placed[1] = true
		placed[2] = true
		placed[3] = true

		/*
		   messageType, r, err := conn.NextReader()
		   fmt.Println("Message Type Received:", string(messageType))
		   fmt.Println("Message Received:", r)
		   if err != nil {
		       return
		   }
		   w, err := conn.NextWriter(messageType)
		   if err != nil {
		       return
		   }
		   if _, err := io.Copy(w, r); err != nil {
		       return
		   }
		   if err := w.Close(); err != nil {
		       return
		   }*/
	}
}

func refreshBoard(conn *websocket.Conn, board Board) {
	var req = Request{"Refresh", "", nil, ""}
	req.MarshalData(board)
	WriteTextMessage(conn, req.Marshal())
}

func gameOver(conn *websocket.Conn) {
	var req = Request{"GameOver", "", nil, ""}
	WriteTextMessage(conn, req.Marshal())
}
