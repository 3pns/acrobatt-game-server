package main

import (
	. "./model"
	. "./utils"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	_ "io"
	"net/http"
	_ "strings"
)

// standard types
//https://github.com/gorilla/websocket/blob/master/conn.go

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	fmt.Println("Launching server on port 8081...")
	var addr = flag.String("addr", ":8081", "http service address")
	http.HandleFunc("/", handleNewConnection)
	http.ListenAndServe(*addr, nil)
}

func handleNewConnection(w http.ResponseWriter, r *http.Request) {
	fmt.Print("New Connection Established:")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		//log.Println(err)
		return
	}
	go startSocket(conn, w, r)
}

func startSocket(conn *websocket.Conn, w http.ResponseWriter, r *http.Request) {
	//création de board à factoriser dans une autre socket ...
	var board Board
	board.InitBoard()
	board.InitPieces()
	board.InitPlayers()
	var turn = 1
	var player *Player = board.Players[0]
	//envoi de la board à la connexion

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

func startGame() {

}
