package model

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	Conn        *websocket.Conn
	token       string
	State       string
	Ai          *AI
	CurrentGame *Game
}

func NewClient(conn *websocket.Conn) Client {
	var client Client
	client.Conn = conn
	client.token = ""
	client.State = "Start"
	client.Ai = nil
	return client
}

func NewAiClient() Client {
	var client Client
	client.Conn = nil
	client.token = ""
	client.State = "Start"
	client.Ai = &AI{}
	return client
}

func (client *Client) IsAi() bool {
	if client.Ai == nil {
		return false
	} else {
		return true
	}
}

func (client *Client) IsAuthenticated() bool {
	if client.token == "" {
		return false
	} else {
		//TODO check token validity
		return true
	}
}

func (client *Client) GameId() int {
  if client == client.CurrentGame.client0 {
    return 0
  } else if client == client.CurrentGame.client1 {
    return 1
  } else if client == client.CurrentGame.client2 {
    return 2
  } else if client == client.CurrentGame.client3 {
    return 3
  } else {
    return -1
  }
}
