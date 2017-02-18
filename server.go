package main

import (
	. "./model"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	_ "io"
	"log"
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
	var playerId = 0
	var board Board
	board.InitBoard()
	board.InitPieces()
	board.InitPlayers()

	//envoi de la board à la connexion
	var req  = Request {"Fetch", "", nil}
	req.MarshalData(board)
	WriteTextMessage(conn, req.Marshal())

	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("read: ", err)
			return
		}
		if mt == websocket.TextMessage {
			request := Request{}
			json.Unmarshal(message, &request)
			if (request.Type == "PlacePiece"){
				piece := Piece{}
				json.Unmarshal(request.Data, &piece)
				fmt.Print("plaçage de Piece")
				if(piece.PlayerId == playerId){
					board.PlacePiece(piece)
					refreshBoard(conn, board)
				}
			}
		}	
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

func refreshBoard (conn *websocket.Conn, board Board){
	var req  = Request {"Refresh", "", nil}
	req.MarshalData(board)
	WriteTextMessage (conn, req.Data)
}

func WriteTextMessage (conn *websocket.Conn, data []byte){
	err := conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		log.Println("write:", err)
	}
}

func startGame() {

}


