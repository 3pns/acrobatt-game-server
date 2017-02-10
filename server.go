package server

import (
	_ "net"
	"fmt"
	_ "bufio"
	_ "strings"
	"encoding/json"
)

import . "./model"
	
func main() {
	fmt.Println("----- Test -----")

	//plateau
	var board Board
	board.InitBoard()
	board.InitPieces()

	//joueur
	player := Player{0, "Bertrand", "yellow", board.Pieces}

	//causing stackoverflow
	//player.initPieces()
	fmt.Println(player)
	
	//pièces
	fmt.Println("----- Piece -----")
	fmt.Println(board.Pieces[0])
	board.PlacePiece(&board.Pieces[4], Square{10, 3, true})

	fmt.Println("----- PRINT TO JSON -----")
	//b, err := json.Marshal(Square{0, 0, true})
	b, err := json.Marshal(player)
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