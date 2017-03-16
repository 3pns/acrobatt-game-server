package main

import (
	. "./model"
	_ "./utils"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	_ "io"
	"net/http"
	_ "strings"
)

// standard types
//https://github.com/gorilla/websocket/blob/master/conn.go

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	fmt.Println("Launching server on port 8081...")
	var addr = flag.String("addr", ":8081", "http service address")
	http.HandleFunc("/", handleNewConnection)
	http.ListenAndServe(*addr, nil)
}

func handleNewConnection(w http.ResponseWriter, r *http.Request) {
	fmt.Print("New Connection Established:")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		//log.Println(err)
		return
	}
	//lancement d'une partie
	var client0 = NewClient(conn)
	var client1 = NewAiClient()
	var client2 = NewAiClient()
	var client3 = NewAiClient()
	clients := []*Client{client0, client1, client2, client3}
	var game = NewGame(clients)
	game.Start()
	go client0.Start()
	go client1.Start()
	go client2.Start()
	go client3.Start()
	fmt.Println("GO !!!")
}
