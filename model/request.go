package model

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Request struct {
	Type       string `json:"type"`
	DataType   string `json:"dataType"`
	Data       []byte `json:"data"`
	CallbackId string `json:"callbackId"`
	Client     *Client `json:"-"`
	Kill bool `json:"-"`
}

func NewRequest (requestType string) Request{
	var req = Request{requestType, "", nil, "", nil, false}
	return req
}

func NewRequestWithCallbackId (requestType string, callbackId string) Request{
	var req = Request{requestType, "", nil, callbackId, nil, false}
	return req
}

//request used to kill goroutines
func NewKillRequest () Request{
	var req = Request{"KILL", "KILL", nil, "666", nil, true}
	return req
}

func (request *Request) MarshalData(t interface{}) {

	board, ok := t.(Board)
	if ok {
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
		b, err := json.Marshal(piece)
		if err != nil {
			fmt.Println(err)
		}
		request.DataType = "Piece"
		request.Data = b
		return
	}
	games, ok := t.(GameSlice)
	if ok {
		b, err := json.Marshal(games)
		if err != nil {
			fmt.Println(err)
		}
		request.DataType = "ListGame"
		request.Data = b
		return
	}
	lobbies, ok := t.(LobbySlice)
	if ok {
		b, err := json.Marshal(lobbies)
		if err != nil {
			fmt.Println(err)
		}
		request.DataType = "ListLobby"
		request.Data = b
		return
	}
	lobby, ok := t.(Lobby)
	if ok {
		b, err := json.Marshal(lobby)
		if err != nil {
			fmt.Println(err)
		}
		request.DataType = "Lobby"
		request.Data = b
		return
	}
	client, ok := t.(Client)
	if ok {
		b, err := json.Marshal(client)
		if err != nil {
			fmt.Println(err)
		}
		request.DataType = "Client"
		request.Data = b
		return
	}
	fmt.Print("->MarshalDataFailed")
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

func (request *Request) DataToInt() int {
	if request.DataType == "int" {
		res, _ := strconv.Atoi(string(request.Data))
		return res
	}
	return -1
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
	fmt.Print("->dispatching")
	if request.Type == "FetchClient" {
		var req = NewRequestWithCallbackId ("FetchClient", request.CallbackId)
		req.MarshalData(*request.Client)
		client.RequestChannel <- *request
	} else if client.State.Current() == "game" && client.CurrentGame != nil {
		fmt.Print("->toCurrentGameRequestChannel->")
		client.CurrentGame.RequestChannel <- *request
	} else if client.State.Current() == "home" {
		fmt.Print("->toServer->")
		GetServer().Process(*request)
	} else if client.State.Current() == "start" && request.Type == "CreateDemo" {
		fmt.Print("->create_demo->")
		client.State.Event("create_demo")
	} else if client.State.Current() == "start" && request.Type == "Authenticate" {
		fmt.Print("->Authenticating->")
		client.Authenticate(request.DataToString())
	} else if client.State.Current() == "lobby" && client.CurrentLobby != nil {
		fmt.Print("->toCurrentLobbyRequestChannel->")
		client.CurrentLobby.RequestChannel <- *request
	} else {
		fmt.Println("->dispatching_failed")
	}
	
}
