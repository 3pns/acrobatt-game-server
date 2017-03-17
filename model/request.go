package model

import (
	"encoding/json"
	_ "flag"
	"fmt"
	_ "github.com/gorilla/websocket"
	_ "io"
	_ "log"
	_ "net/http"
	_ "strings"
)

type Request struct {
	Type       string `json:"type"`
	DataType   string `json:"dataType"`
	Data       []byte `json:"data"`
	CallbackId string `json:"callbackId"`
	Client     *Client
}

func (request *Request) MarshalData(t interface{}) {

	board, ok := t.(Board)
	if ok {
		fmt.Println("Marshalling Board")
		b, err := json.Marshal(board)
		if err != nil {
			fmt.Println(err)
		}
		request.DataType = "Board"
		request.Data = b
		return
	}
	player, ok := t.(Player)
	if ok {
		fmt.Println("Marshalling Player")
		b, err := json.Marshal(player)
		if err != nil {
			fmt.Println(err)
		}
		request.DataType = "Player"
		request.Data = b
		return
	}
	piece, ok := t.(Piece)
	if ok {
		fmt.Println("Marshalling Piece")
		b, err := json.Marshal(piece)
		if err != nil {
			fmt.Println(err)
		}
		request.DataType = "Piece"
		request.Data = b
		return
	}

}

func (request *Request) Marshal() []byte {
	marshaleldrequest, err := json.Marshal(request)
	if err != nil {
		fmt.Println(err)
	}
	return marshaleldrequest
}

func (request *Request) Unmarshal() {
	fmt.Print("Unmarshalling")
}

func (request *Request) DataToString() string {
	if request.DataType == "string" {
		return string(request.Data)
	}
	return ""
}

func (request *Request) HasClient() bool {
	if request.Client != nil {
		return true
	} else {
		return false
	}
}

func (request *Request) Dispatch() {
	var client = request.Client
	if client.State.Current() == "game" && (request.Type == "PlacePiece" || request.Type == "PlaceRandom" || request.Type == "Fetch" || request.Type == "FetchPlayer" || request.Type == "Concede") {
		client.CurrentGame.RequestChannel <- *request
	} else if client.State.Current() == "start" && request.Type == "CreateDemo" {
		client.State.Event("create_demo")
	} else if client.State.Current() == "start" && request.Type == "Authenticate" {
		client.Authenticate(request.DataToString())
	}
}
