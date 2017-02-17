package main

import (
	. "./model"
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"net"
	_ "strings"
	"io"
	"github.com/gorilla/websocket"
	"flag"
	"log"
)

// standard types
//https://github.com/gorilla/websocket/blob/master/conn.go

/*
func main() {
	fmt.Println("----- Test -----")

	//plateau
	var board Board
	board.InitBoard()
	board.InitPieces()

	//joueur
	player := Player{0, "Joueur", "blue", board.Pieces}
	ai1 := Player{1, "AI-1", "green", board.Pieces}
	ai2 := Player{2, "AI-2", "yellow", board.Pieces}
	ai3 := Player{3, "AI-3", "red", board.Pieces}
	player.Init()
	ai1.Init()
	ai2.Init()
	ai3.Init()
	board.Players = []*Player{&player, &ai1, &ai2, &ai3}
	fmt.Println(player)

	board.PlacePiece(&board.Pieces[len(board.Pieces)-1], board.Squares[10][10])

	fmt.Println("----- PRINT TO JSON -----")
	b, err := json.Marshal(board)
	if err != nil {
		fmt.Println(err)
	}
	myJson := string(b) // converti byte en string
	fmt.Println(myJson)

	board.PrintBoard()
	fmt.Println("\n----- Game Over -----")
}*/

// code server Websocket
var upgrader = websocket.Upgrader{}
func main() {

	fmt.Println("Launching server on port 8081...")

	//Bytestream Listen
	// listen on all interfaces
	/*
	ln, _ := net.Listen("tcp", ":8081")
	// accept connection on port
	waitingState:
	fmt.Print("waiting cest ici que sa wait")
	conn, _ := ln.Accept()
	go startConnect (conn)*/

	//WebSocket Listen
	var addr = flag.String("addr", "127.0.0.1:8081", "http service address")
	http.HandleFunc("/", handleNewConnection)
	http.ListenAndServe(*addr, nil)


	// run loop forever (or until ctrl-c)

	//on retourne à l'état d'attente d'une connexion
	//goto waitingState
}

func handleNewConnection(w http.ResponseWriter, r *http.Request) {
	fmt.Print("New Connection Established:")
  conn, err := upgrader.Upgrade(w, r, nil)
  if err != nil {
      //log.Println(err)
      return
  }
  go startSocket (conn, w, r)
}

func startConnect (conn net.Conn) {
	var board Board
	board.InitBoard()
	board.InitPieces()
	player := Player{0, "Joueur", "blue", board.Pieces}
	ai1 := Player{1, "AI-1", "green", board.Pieces}
	ai2 := Player{2, "AI-2", "yellow", board.Pieces}
	ai3 := Player{3, "AI-3", "red", board.Pieces}
	player.Init()
	ai1.Init()
	ai2.Init()
	ai3.Init()
	board.Players = []*Player{&player, &ai1, &ai2, &ai3}

	//envoi de la board
	board.Refresh(conn)
	for {
		// will listen for message to process ending in newline (\n)
		message, err := bufio.NewReader(conn).ReadString('\n')

		//detection de la fin de la connexion
		if err != nil {
		   if err == io.EOF {
		     break
		   }
		}
		fmt.Print("Message Received:", string(message))
		board.Refresh(conn)
		if (string(message) == "QUIT") {
			break
		}
	}

}

func startSocket (conn *websocket.Conn, w http.ResponseWriter, r *http.Request){
	var board Board
	board.InitBoard()
	board.InitPieces()
	player := Player{0, "Joueur", "blue", board.Pieces}
	ai1 := Player{1, "AI-1", "green", board.Pieces}
	ai2 := Player{2, "AI-2", "yellow", board.Pieces}
	ai3 := Player{3, "AI-3", "red", board.Pieces}
	player.Init()
	ai1.Init()
	ai2.Init()
	ai3.Init()
	board.Players = []*Player{&player, &ai1, &ai2, &ai3}

	b, err := json.Marshal(board)
	if err != nil {
		fmt.Println(err)
	}

	err = conn.WriteMessage(1, b)
	if err != nil {
		log.Println("write:", err)
	}

	for {
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
	    }
	}
}
/*var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}*/

/*
func handler(w http.ResponseWriter, r *http.Request){
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        //log.Println(err)
        return
    }
		for {
		    messageType, r, err := conn.NextReader()
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
		    }
}
}*/