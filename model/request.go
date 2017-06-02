package model

import (
	. "../jsonapi"
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"strconv"
)

type Request struct {
	Type       string  `json:"type"`
	DataType   string  `json:"dataType"`
	Data       []byte  `json:"data"`
	CallbackId string  `json:"callbackId"`
	Client     *Client `json:"-"`
}

func NewRequest(requestType string) Request {
	var req = Request{requestType, "", nil, "", nil}
	return req
}

func NewRequestWithCallbackId(requestType string, callbackId string) Request {
	var req = Request{requestType, "", nil, callbackId, nil}
	return req
}

func (request *Request) MarshalData(t interface{}) {

	board, ok := t.(Board)
	if ok {
		b, err := json.Marshal(board)
		if err != nil {
			log.Warn(err)
		}
		request.DataType = "Board"
		request.Data = b
		return
	}
	player, ok := t.(Player)
	if ok {
		b, err := json.Marshal(player)
		if err != nil {
			log.Warn(err)
		}
		request.DataType = "Player"
		request.Data = b
		return
	}
	piece, ok := t.(Piece)
	if ok {
		b, err := json.Marshal(piece)
		if err != nil {
			log.Warn(err)
		}
		request.DataType = "Piece"
		request.Data = b
		return
	}
	games, ok := t.(GameSlice)
	if ok {
		b, err := json.Marshal(games)
		if err != nil {
			log.Warn(err)
		}
		request.DataType = "ListGame"
		request.Data = b
		return
	}
	lobbies, ok := t.(LobbySlice)
	if ok {
		b, err := json.Marshal(lobbies)
		if err != nil {
			log.Warn(err)
		}
		request.DataType = "ListLobby"
		request.Data = b
		return
	}
	clients, ok := t.(ClientSlice)
	if ok {
		b, err := json.Marshal(clients)
		if err != nil {
			log.Warn(err)
		}
		request.DataType = "ListClient"
		request.Data = b
		return
	}
	lobby, ok := t.(Lobby)
	if ok {
		b, err := json.Marshal(lobby)
		if err != nil {
			log.Warn(err)
		}
		request.DataType = "Lobby"
		request.Data = b
		return
	}
	client, ok := t.(Client)
	if ok {
		client.MyState = client.State.Current()
		if client.CurrentLobby != nil {
			client.MyLobbyId = client.CurrentLobby.Id
		}
		if client.CurrentGame != nil {
			client.MyGameId = client.CurrentGame.Id
		}
		b, err := json.Marshal(client)
		if err != nil {
			log.Warn(err)
		}
		request.DataType = "Client"
		request.Data = b
		return
	}
	log.Warn("MarshalData Failed")
}

func (request *Request) Marshal() []byte {
	marshaleldrequest, err := json.Marshal(request)
	if err != nil {
		log.Warn(err)
	}
	return marshaleldrequest
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
	client.UpdateTrace("->dispatching")
	if request.Type == "FetchClient" {
		var req = NewRequestWithCallbackId("FetchClient", request.CallbackId)
		req.MarshalData(*request.Client)
		client.UpdateTrace("->FetchClient->")
		client.RequestChannel <- req
	} else if request.Type == "DestroyClient" {
		client.UPTrace("->destroying_client["+strconv.Itoa(client.Id)+"]")
		GetServer().CleanClient(client)
		client.Stop()
		client.Conn.Close()
	} else if request.Type == "Disconnect" {
		client.State.Event("disconnect")
	}  else if client.State.Current() == "game" && client.CurrentGame != nil {
		client.UpdateTrace("->toCurrentGameRequestChannel->")
		client.CurrentGame.RequestChannel <- *request
	} else if client.State.Current() == "home" {
		client.UpdateTrace("->toServer->")
		GetServer().Process(*request)
	} else if client.State.Current() == "start" && request.Type == "CreateDemo" {
		client.UpdateTrace("->create_demo->")
		client.PrintTrace()
		client.State.Event("create_demo")
	} else if client.State.Current() == "start" && request.Type == "Authenticate" {
		client.UpdateTrace("->Authenticating->")
		auth := AuthenticateJson{}
		json.Unmarshal(request.Data, &auth)
		client.Authenticate(auth)
	} else if client.State.Current() == "lobby" && client.CurrentLobby != nil {
		client.UpdateTrace("->toCurrentLobbyRequestChannel->")
		client.CurrentLobby.RequestChannel <- *request
	} else {
		client.UPTrace("->dispatching_failed")
	}

}
