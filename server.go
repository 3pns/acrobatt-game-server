package main

import (
	_ "net"
	"fmt"
	_ "bufio"
	_ "strings"
)


type Board struct {
	squareTab [20][20] *Square
}

type Square struct {
	x int
	y int
	empty bool
}

type Piece struct {
	id int
	origin Square
	rotation string
	flipped bool
	player Player
}

type Player struct {
	name string
	color string
}

	
func main() {
	fmt.Println("#### Test ####")

	player := Player{"Bertrand", "yellow"}
	fmt.Println(player)

	var board Board
	var tab [20][20] *Square
	mySquare := Square{1, 2, false}
	
	
	tab[0][0] = &mySquare
	board.squareTab = tab
	initBoard(&board)
	printBoard(&board)
	fmt.Println("")
	fmt.Print(board.squareTab)
}

func printBoard(board *Board){
	for i := 0; i < 20; i++ {
		for j := 0; j < 20; j++ {
			
			if board.squareTab[i][j] == nil {
				fmt.Print("0")
				
			} else {
				fmt.Print("#")
			}
		}
		fmt.Println("")
	}
}

func initBoard(board *Board){
	fmt.Println("initializing board")
	for i := 0; i < 20; i++ {
		for j := 0; j < 20; j++ {
			(*board).squareTab[i][j] = &Square{i, j, true}
		}
	}
	fmt.Println("board initialized with success !")
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