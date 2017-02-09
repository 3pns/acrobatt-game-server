package main

import (
	_ "net"
	"fmt"
	_ "bufio"
	_ "strings"
	"encoding/json"
)


type Board struct {
	Squares [20][20] *Square `json:"squares"`
	Pieces [] Piece `json:"pieces"`
	Players [] *Player `json:"pieces"`
}

type Square struct {
	X int `json:"x"`
	Y int `json:"y"`
	Empty bool `json:"empty"`
}

type Piece struct {
	Id int `json:"id"`
	Cubes []Cube `json:"cubes"`
	Origin *Square `json:"origin"`
	Rotation string `json:"rotation"`
	Flipped bool `json:"flipped"`
	Player *Player `json:"player"`
}

type Cube struct {
	X int
	Y int
}

type Player struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Color string `json:"color"`
	Pieces [] Piece `json:"pieces"`
}

	
func main() {
	fmt.Println("----- Test -----")

	//plateau
	var board Board
	initBoard(&board)
	initPieces(&board)

	//joueur
	player := Player{0, "Bertrand", "yellow", board.Pieces}

	//causing stackoverflow
	//player.initPieces()
	fmt.Println(player)
	
	//pièces
	fmt.Println("----- Piece -----")
	fmt.Println(board.Pieces[0])
	board.placePiece(&board.Pieces[0], Square{10, 3, true})

	fmt.Println("----- PRINT TO JSON -----")
	//b, err := json.Marshal(Square{0, 0, true})
	b, err := json.Marshal(player)
  if err != nil {
      fmt.Println(err)
  }
  myJson := string(b) // converti byte en string
	fmt.Println(myJson)

	printBoard(&board)
	fmt.Println("\n----- Game Over -----")
}

func printBoard(board *Board){
	for i := 0; i < 20; i++ {
		setWhiteBackground()
		fmt.Print(" ")
		for j := 0; j < 20; j++ {
			
			if board.Squares[j][i].Empty == true {
				printBlack("▇ ")
				
			} else {
				printRed("▇ ")
			}
		}
		fmt.Print(" ")
		setBlackBackground ()
		fmt.Println("")
	}
	printReset()
}

func initBoard(board *Board){
	fmt.Println("initializing board")
	for i := 0; i < 20; i++ {
		for j := 0; j < 20; j++ {
			(*board).Squares[i][j] = &Square{i, j, true}
		}
	}
	fmt.Println("board initialized with success !\n")
}

func initPieces(board *Board){
	fmt.Println("generating pieces")

	var piece Piece
	piece.Id = 0
	piece.Rotation = "N"
	piece.Flipped = false
	piece.Cubes = append(piece.Cubes, Cube{0, 0})
	board.Pieces = append(board.Pieces, piece)

	fmt.Println("pieces generated with success !\n")
}

func (board Board) placePiece(piece *Piece, square Square) {
	board.Squares[square.X][square.Y].Empty = false
	piece.Origin = board.Squares[square.X][square.Y]
}

func (player *Player) initPieces() {
for index,piece := range player.Pieces {
	player.Pieces[index].Player = player
	piece.Player = player
  // index is the index where we are
  // element is the element from someSlice for where we are
}
	fmt.Println("init player pieces")
}

func setWhiteBackground (){
	fmt.Print("\x1B[47m")
}
func setBlackBackground (){
	fmt.Print("\x1B[40m")
}

func printReset(){
	fmt.Print("\x1B[0m")
}

func printBlack(str string){
	fmt.Print("\x1B[30m" + str)
}

func printRed(str string){
	fmt.Print("\x1B[31m" + str)
}

func printGreen(str string){
	fmt.Print("\x1B[32m" + str)
}

func printYellow(str string){
	fmt.Print("\x1B[33m" + str)
}

func printBlue(str string){
	fmt.Print("\x1B[34m" + str)
}

func printWhite(str string){
	fmt.Print("\x1B[37m" + str)
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