package main

import (
	. "./model"
	. "./utils"
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
	var ai1 = NewAiClient()
	var ai2 = NewAiClient()
	var ai3 = NewAiClient()
	var game = NewGame(&client, &ai1, &ai2, &ai3)
	game.Start()

	//lancement du joueur
	go startSocket(&client, w, r)
}

func startSocket(client *Client, w http.ResponseWriter, r *http.Request) {
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
			if request.Type == "PlacePiece" {
				fmt.Println("Message de type PlacePiece detected !")
				//TODO placePiece sur board dans client.Game.Board
			} else if request.Type == "PlaceRandom" {
				fmt.Println("Message de type PlaceRandom detected !")
				//TODO placeRandom sur board dans game du client
			} else if request.Type == "Fetch" {
				fmt.Println("Message de type Fetch detected !")
				var req = Request{"Fetch", "", nil, request.CallbackId}
				req.MarshalData(client.CurrentGame.Board())
				WriteTextMessage(conn, req.Marshal())
			} else if request.Type == "FetchPlayer" {
				fmt.Println("Message de type FetchPlayer detected !")
				//TODO FetchPlayer depuis client.Game.??????
				/*var req = Request{"FetchPlayer", "Player", nil, request.CallbackId}
				req.MarshalData(*client.CurrentGame.Board().Players[client.GameId()])
				WriteTextMessage(conn, req.Marshal())*/
			}
		}
	}
}
