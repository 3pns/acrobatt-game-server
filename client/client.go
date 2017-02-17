package main

import (
	. "../model"
	"encoding/json"
	"net"
	"fmt"
	"bufio"
	"os"
)

//client Websocket
func main() {
	
	conn, _ := net.Dial("tcp", "127.0.0.1:8081") // connect to this socket local
	//conn, _ := net.Dial("tcp", "94.23.249.62:8081") // production server

	for {
		// read in input from stdin
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("JSON to send: ")
		text, _ := reader.ReadString('\n')
		// send to socket
		fmt.Fprintf(conn, text+"\n")
		// listen for reply
		jsonMessage, _ := bufio.NewReader(conn).ReadString('\n')
		board := Board{}
		json.Unmarshal([]byte(jsonMessage), &board)
		fmt.Println("###### DATA #####")
		fmt.Println(board.Squares[2][2])
		board.PrintBoard()
	}
}
