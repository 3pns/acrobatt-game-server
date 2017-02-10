package main

import (
	_ "net"
	"fmt"
	_ "bufio"
	_ "strings"
	"encoding/json"
	"strconv"
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
	board.placePiece(&board.Pieces[4], Square{10, 3, true})

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
	piece.Rotation = "S"
	piece.Flipped = false
	piece.Cubes = append(piece.Cubes, Cube{0, 0})
	board.Pieces = append(board.Pieces, piece)

	var piece1 Piece
	piece1.Id = 1
	piece1.Rotation = "S"
	piece1.Flipped = false
	piece1.Cubes = append(piece1.Cubes, Cube{0, 0})
	piece1.Cubes = append(piece1.Cubes, Cube{0, 1})
	board.Pieces = append(board.Pieces, piece1)

	var piece2 Piece
	piece2.Id = 2
	piece2.Rotation = "S"
	piece2.Flipped = false
	piece2.Cubes = append(piece2.Cubes, Cube{0, 0})
	piece2.Cubes = append(piece2.Cubes, Cube{0, 1})
	piece2.Cubes = append(piece2.Cubes, Cube{1, 0})
	board.Pieces = append(board.Pieces, piece2)

	var piece3 Piece
	piece3.Id =3
	piece3.Rotation = "S"
	piece3.Flipped = false
	piece3.Cubes = append(piece3.Cubes, Cube{0, 0})
	piece3.Cubes = append(piece3.Cubes, Cube{0, 1})
	piece3.Cubes = append(piece3.Cubes, Cube{0, 2})
	board.Pieces = append(board.Pieces, piece3)

	var piece4 Piece
	piece4.Id =3
	piece4.Rotation = "S"
	piece4.Flipped = false
	piece4.Cubes = append(piece4.Cubes, Cube{0, 0})
	piece4.Cubes = append(piece4.Cubes, Cube{1, 0})
	piece4.Cubes = append(piece4.Cubes, Cube{0, 1})
	piece4.Cubes = append(piece4.Cubes, Cube{0, 2})
	board.Pieces = append(board.Pieces, piece4)

	fmt.Println("pieces generated with success !\n")
}

func (board Board) placePiece(piece *Piece, square Square) {
	board.Squares[square.X][square.Y].Empty = false
	piece.Origin = board.Squares[square.X][square.Y]
	piece.Flipped = true
	piece.Rotation = "N"

	//postionnement de la pièce dans l'espace


	//1 - vérifier si on a le droit de placer la pièce
	//2 - placer la pièce
	fmt.Println("----- Plaçage d'une pièce -----")
	for _,cube := range piece.Cubes {
		var xFactor = 1
		var yFactor = 1
		var xBoardValue = cube.X
		var yBoardValue = cube.Y
		if piece.Flipped {
			xFactor = -xFactor
		}
		fmt.Println("Avant : x :" + strconv.Itoa(xBoardValue) + ", y: " + strconv.Itoa(yBoardValue))
		if piece.Rotation == "W" {
			xBoardValue = -cube.Y
			yBoardValue = -cube.X
		} else if piece.Rotation == "N" {
			xBoardValue = -cube.X
			yBoardValue = -cube.Y
		} else if piece.Rotation == "E" {
			xBoardValue = cube.Y
			yBoardValue = cube.X
		}
			fmt.Println("Apres : x :" + strconv.Itoa(xBoardValue) + ", y: " + strconv.Itoa(yBoardValue))
			board.Squares[square.X + xFactor * xBoardValue][square.Y + yFactor * yBoardValue].Empty = false
		}
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