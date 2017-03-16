package main

import (
	. "./model"
	_ "./utils"
	"encoding/json"
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
	var client = NewClient(conn)
	var aiClient1 = NewAiClient()
	var aiClient2 = NewAiClient()
	var aiClient3 = NewAiClient()
	clients := []*Client{&client, aiClient1, aiClient2, aiClient3}
	var game = NewGame(clients)
	game.Start()
	fmt.Println("Starting AI 1")
	go aiClient1.Ai.Start()
	fmt.Println("Starting AI 2")
	go aiClient2.Ai.Start()
	fmt.Println("Starting AI 3")
	go aiClient3.Ai.Start()
	fmt.Println("GO !!!")

	//lancement du joueur
	go startRequestDispatcher(&client, w, r)
}

func startRequestDispatcher(client *Client, w http.ResponseWriter, r *http.Request) {
	var conn = client.Conn
	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("read: ", err)
			return
		}
		if mt == websocket.TextMessage {
			request := Request{}
			json.Unmarshal(message, &request)
			request.Client = client
			if request.Type == "PlacePiece" || request.Type == "PlaceRandom" || request.Type == "Fetch" || request.Type == "FetchPlayer" {
				client.CurrentGame.RequestChannel <- request
			}
		}
	}
}
