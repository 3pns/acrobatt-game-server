package main

import (
	. "../model"
	"encoding/json"
	_"net"
	"fmt"
	"bufio"
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
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	var addr = flag.String("addr", "127.0.0.1:8081", "http service address")
	//var addr = flag.String("addr", "94.23.249.62:8081", "http service address")
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/"}
	fmt.Println("connecting to ", u.String())
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Print("dial:", err)
	}
	defer conn.Close()
	done := make(chan struct{})

	//on read les messages dans une goroutine
	go func() {
		defer conn.Close()
		defer close(done)
		for {
			mt, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("read: ", err)
				return
			}
			if mt == websocket.TextMessage {

				clientRequest := ClientRequest{}
				json.Unmarshal(message, &clientRequest)
				if (clientRequest.DataType == "Board"){
					fmt.Println("Data de type Board détécté: ")
					board := Board{}
					json.Unmarshal(clientRequest.Data, &board)
					board.PrintBoard()
				}
				if (clientRequest.DataType == "Pieces"){
					fmt.Println("Data de type Board détécté: ")
					board := Board{}
					json.Unmarshal(clientRequest.Data, &board)
					board.PrintBoard()
				}
				if (clientRequest.DataType == "Board"){
					fmt.Println("Data de type Board détécté: ")
					board := Board{}
					json.Unmarshal(clientRequest.Data, &board)
					board.PrintBoard()
				}	
				fmt.Print("Enter text: ")
			}
		}
	}()

	//var req ClientRequest ("placePiece", "Piece", )
	//conn.WriteMessage(websocket.TextMessage, []byte(req))

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		fmt.Print("Enter text: ")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		sz := len(text)//on enlève le dernier \n
		text = text[:sz-1]
		if (text == "exit"){
			return
		} else {
			err = conn.WriteMessage(websocket.TextMessage, []byte(text))
			if err != nil {
				log.Println("write:", err)
				return
			}
		}
	}
/*
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
*/
}
