package main

import (
	_ "net"
	"fmt"
	_ "bufio"
	_ "strings"
	"encoding/json"
)


type Board struct {
	Squares [20][20] Square `json:"squares"`
	pieces [] *Piece
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
	Player Player `json:"player"`
}

type Cube struct {
	X int
	Y int
}

type Player struct {
	name string
	color string
}

	
func main() {
	fmt.Println("----- Test -----")

	//joueur
	player := Player{"Bertrand", "yellow"}
	fmt.Println(player)

	//plateau
	var board Board
	initBoard(&board)
	
	//pièces
	fmt.Println("----- Piece -----")
	var piece Piece
	piece.Id = 0
	piece.Rotation = "N"
	piece.Flipped = false
	piece.Player = player
	board.Squares[10][15].Empty = false
	piece.Origin = &board.Squares[10][15]
	piece.Cubes = append(piece.Cubes, Cube{0, 0})

	fmt.Println("----- PRINT TO JSON -----")
	//b, err := json.Marshal(Square{0, 0, true})
	b, err := json.Marshal(piece)
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
			
			if board.Squares[i][j].Empty == true {
				printBlack("▇ ")
				
			} else {
				//printRed("▪") //petit carré
				//printRed("⎕") // carré vide
				//printBlack("■ ") // carré moyen
				if board.Squares[i][j].Empty {
					printBlack("▇ ")
					} else {
						printRed("▇ ")
					}
				 // 7/8 pavé
				//printBlack("▆ ") // 3/4
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
			(*board).Squares[i][j] = Square{i, j, true}
		}
	}
	fmt.Println("board initialized with success !\n")
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