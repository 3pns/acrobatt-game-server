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
	id int
	cubes []Cube
	origin *Square
	rotation string
	flipped bool
	player Player
}

type Cube struct {
	x int
	y int
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
	piece.id = 0
	//piece.origin = board.squares[10][10]
	piece.rotation = "N"
	piece.flipped = false
	piece.player = player
	fmt.Println(piece)
	fmt.Println(piece.origin)
	
	fmt.Println("----- BOARD TO JSON -----")
	b, err := json.Marshal(&board.Squares[0][0])
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
			
			if board.Squares[i][j].Empty == false {
				fmt.Print("0")
				
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