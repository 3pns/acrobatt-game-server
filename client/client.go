package main

import (
	. "../model"
	. "../utils"
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
	var menu = "Choose: (1)PlacePiece (2)Refresh (3)Fetch (4)FetchPlayers (exit) Close the game  "
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
	cboard := make(chan Board, 1)
	cplayer := make(chan Player, 1)
	var myBoard Board
	var myPlayer Player
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

				clientRequest := Request{}
				json.Unmarshal(message, &clientRequest)
				if (clientRequest.DataType == "Board"){
					fmt.Println("Data de type Board détécté !!! ")
					board := Board{}
					json.Unmarshal(clientRequest.Data, &board)
					cboard <- board
				}else if (clientRequest.DataType == "Piece"){
					fmt.Println("Data de type Pieces détécté !!! ")
					piece := Piece{}
					json.Unmarshal(clientRequest.Data, &piece)
				}else if (clientRequest.DataType == "Player"){
					fmt.Println("Data de type Player détécté !!! ")
					player := Player{}
					json.Unmarshal(clientRequest.Data, &player)
					fmt.Println(player)
					cplayer <- player
				}	
				fmt.Print(menu)
			}
		}
	}()

	//var req ClientRequest ("placePiece", "Piece", )
	//conn.WriteMessage(websocket.TextMessage, []byte(req))

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		fmt.Print(menu)
		var text = getInput()
		if (text == "exit"){
			return
		} else if text == "1" {
			myPlayer.Pieces[17].Origin = &Square {10, 10, nil}
			var req  = Request {"PlacePiece", "Piece", nil}
			fmt.Println(getJson(myPlayer.Pieces))
			//var piececopy = myPlayer.Pieces[17]
			req.MarshalData(myPlayer.Pieces[17])
			WriteTextMessage(conn, req.Marshal())
		} else if text == "2" {
		    select {
			    case newBoard, ok := <-cboard:
			    	//nouvelle donnée dans le buffer
			        if ok {
			            myBoard = newBoard
			            myBoard.PrintBoard()
			        } else {
			            fmt.Println("Channel closed!")
			        }
			    default:
			        myBoard.PrintBoard()
    		}
		    select {
			    case newPlayer, ok := <-cplayer:
			    	//nouvelle donnée dans le buffer
			        if ok {
			            myPlayer = newPlayer
			        } else {
			            fmt.Println("Channel closed!")
			        }
			    default:
			        fmt.Println(myPlayer)
    		}
		}
		if text == "3" {
			var req  = Request {"Fetch", "", nil}
			//fmt.Println(getJson(req))
			WriteTextMessage(conn, req.Marshal())
		}
		if text == "4" {
			var req  = Request {"FetchPlayer", "", nil}
			fmt.Println(getJson(req))
			WriteTextMessage(conn, req.Marshal())
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

func getInput () string{
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	sz := len(text)//on enlève le dernier \n
	text = text[:sz-1]
	return text
}

func getJson (t interface{}) string{
	b, err := json.Marshal(t)
	if err != nil {
		fmt.Print("getJson Marshell Error :", err)
	}
	return string (b)
}