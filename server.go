package main

import (
	. "./model"
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"io"
)
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
func main() {

	fmt.Println("Launching server on port 8081...")

	// listen on all interfaces
	ln, _ := net.Listen("tcp", ":8081")
	// accept connection on port
	waitingState:
	fmt.Print("waiting cest ici que sa wait")
	conn, _ := ln.Accept()
	go startConnect (conn)
	// run loop forever (or until ctrl-c)

	//on retourne à l'état d'attente d'une connexion
	goto waitingState
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
	b, err := json.Marshal(board)
	if err != nil {
		fmt.Println(err)
	}
	conn.Write(b)
	for {
		// will listen for message to process ending in newline (\n)
		message, err := bufio.NewReader(conn).ReadString('\n')

		//detection de la fin de la connexion
		if err != nil {
		   if err == io.EOF {
		     break
		   }
		}
		// output message received
		fmt.Print("Message Received:", string(message))
		// sample process for string received
		newmessage := strings.ToUpper(message)
		// send new string back to client
		conn.Write([]byte(newmessage + "\n"))
		if (string(message) == "QUIT") {
			break
		}
	}
}