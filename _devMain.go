package main

import (
	. "./model"
	_ "bufio"
	"encoding/json"
	"fmt"
	_ "net"
	_ "strings"
)

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
}

/*

// code server Websocket
func main() {

	fmt.Println("Launching server on port 8081...")

	// listen on all interfaces
	ln, _ := net.Listen("tcp", ":8081")

	// accept connection on port
	conn, _ := ln.Accept()
	// run loop forever (or until ctrl-c)
	for {
		// will listen for message to process ending in newline (\n)
		message, _ := bufio.NewReader(conn).ReadString('\n')
		// output message received
		fmt.Print("Message Received:", string(message))
		// sample process for string received
		newmessage := strings.ToUpper(message)
		// send new string back to client
		conn.Write([]byte(newmessage + "\n"))

		if (string(message) == "QUIT") {
			return
		}
	}
}*/
