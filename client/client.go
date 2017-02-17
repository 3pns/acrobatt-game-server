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
	"log"
	"time"
	"os"
	"os/signal"
)

//client Websocket
func main() {
	
	//conn, _ := net.Dial("tcp", "127.0.0.1:8081") // connect to this socket local
	//conn, _ := net.Dial("tcp", "94.23.249.62:8081") // production server

	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	var addr = flag.String("addr", "127.0.0.1:8081", "http service address")
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/"}
	fmt.Println("connecting to ", u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Print("dial:", err)
	}
	defer c.Close()
	fmt.Println("TEST1")
	done := make(chan struct{})
	fmt.Println("TEST2")

	//on read les messages dans une goroutine
	go func() {
		fmt.Println("TEST4")
		defer c.Close()
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				fmt.Println("read: ", err)
				return
			}
			myJson := string(message)
			fmt.Println("recv: ", myJson)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")
			// To cleanly close a connection, a client should send a close
			// frame and wait for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			c.Close()
			return
		}
	}
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
