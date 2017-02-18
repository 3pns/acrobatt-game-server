package main

import (
	. "./model"
	_"encoding/json"
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
	var board Board
	board.InitBoard()
	board.InitPieces()
	board.InitPlayers()

	//envoi de la board à la connexion
	var req  = ClientRequest {"Fetch", "", nil}
	req.MarshalData(board)
	err := conn.WriteMessage(websocket.TextMessage, req.Marshal())
	if err != nil {
		log.Println("write:", err)
	}

	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("read: ", err)
			return
		}
		if mt == websocket.TextMessage {
			fmt.Println("message de type TextMessage détécté")
			myJson := string(message)
			fmt.Println("message reçu: ", myJson)
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

func startGame() {

}


