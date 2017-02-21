package main

import (
	. "../model"
	"bufio"
	"encoding/json"
	"fmt"
	_ "net"
	_ "strings"
	"flag"
	"github.com/gorilla/websocket"
	"net/url"
	. "../utils"
	"os"
)

//client de test websocket pour aider au dévelopeme,t

//client de test offline pour aider au dévelopement
func main() {

	//connexion au serveur
	var addr = flag.String("addr", "127.0.0.1:8081", "http service address")
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/"}
	fmt.Println("connecting to ", u.String())
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Print("dial:", err)
	}
	defer conn.Close()

	board := Board{}
	mt, message, err := conn.ReadMessage()
	if err != nil {
		fmt.Println("read: ", err)
		return
	}
	if mt == websocket.TextMessage {
		clientRequest := Request{}
		json.Unmarshal(message, &clientRequest)
		if (clientRequest.DataType == "Board"){
			fmt.Println("Data de type Board détécté !!! ")
			json.Unmarshal(clientRequest.Data, &board)
		}
	}

	//var player = board.Players[0]
	board.Pieces[10].Origin = board.Squares[0][1]
	board.Pieces[10].Rotation = "E"

	var req  = Request {"PlacePiece", "Piece", nil}
	req.MarshalData(board.Pieces[10])
	WriteTextMessage(conn, req.Marshal())
	
	bufio.NewReader(os.Stdin).ReadBytes('\n') 
}
/*
func main2() {
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
	board.Pieces[10].Origin = board.Squares[2][2]
	board.Pieces[10].Rotation = "S"
	player.PlacePiece(board.Pieces[10], &board)
*/
	/*fmt.Println("----- PRINT TO JSON -----")
	b, err := json.Marshal(board)
	if err != nil {
		fmt.Println(err)
	}
	myJson := string(b) // converti byte en string
	fmt.Println(myJson)

	board.PrintBoard()
	fmt.Println("\n----- Game Over -----")
}
*/
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
