package main

import (
	//. "../model"
	_"encoding/json"
	_"net"
	"fmt"
	_"bufio"
	_"os"
	"github.com/gorilla/websocket"
	"flag"
	"net/url"
)

//client Websocket
func main() {
	
	//conn, _ := net.Dial("tcp", "127.0.0.1:8081") // connect to this socket local
	//conn, _ := net.Dial("tcp", "94.23.249.62:8081") // production server
	var addr = flag.String("addr", "127.0.0.1:8081", "http service address")
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/"}
	fmt.Print("connecting to %s", u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Print("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})
	go func() {
		defer c.Close()
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				fmt.Print("read:", err)
				return
			}
			fmt.Print("recv: %s", message)
		}
	}()



	/*for {
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
	}*/
}
