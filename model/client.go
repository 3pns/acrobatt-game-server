package model

import (
	. "../jsonapi"
	. "../utils"
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	"github.com/looplab/fsm"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	Id             int             `json:"id"`
	MyState        string          `json:"state"`
	MyLobbyId      int             `json:"lobby_id"`
	MyGameId       int             `json:"game_id"`
	Pseudo         string          `json:"pseudo"`
	Ping           int64           `json:"ping"`
	retry          int             `json:"-"`
	Conn           *websocket.Conn `json:"-"`
	State          *fsm.FSM        `json:"-"`
	Ai             *AI             `json:"ai"`
	CurrentGame    *Game           `json:"-"`
	CurrentLobby   *Lobby          `json:"-"`
	RequestChannel chan Request    `json:"-"`
	trace          string          `json:"-"`
	listening      bool            `json:"-"`
	terminating    bool            `json:"-"`
	quitReader     chan int        `json:"-"`
	quitWriter     chan int        `json:"-"`
	quitPing       chan int        `json:"-"`
	dcByOtherClient chan int        `json:"-"`
}

type ClientSlice struct {
	Clients []*Client `json:"clients"`
}

type ClientFactory struct {
	id int
}

func NewClientFactory() *ClientFactory {
	var factory = new(ClientFactory)
	factory.id = 1
	return factory
}

func (factory *ClientFactory) NewClient(conn *websocket.Conn) *Client {
	var client Client
	client.Id = -1 //factory.id
	//factory.id++
	client.Ai = nil
	client.Conn = conn
	client.RequestChannel = make(chan Request, 100)

	client.State = fsm.NewFSM(
		"start",
		fsm.Events{
			{Name: "create_demo", Src: []string{"start"}, Dst: "game"},
			{Name: "quit_demo", Src: []string{"game"}, Dst: "start"},
			{Name: "authenticate", Src: []string{"start"}, Dst: "home"},
			{Name: "disconnect", Src: []string{"home", "lobby", "game"}, Dst: "start"},
			{Name: "join_lobby", Src: []string{"home"}, Dst: "lobby"},
			{Name: "create_lobby", Src: []string{"home"}, Dst: "lobby"},
			{Name: "quit_lobby", Src: []string{"lobby"}, Dst: "home"},
			{Name: "join_game", Src: []string{"home", "lobby"}, Dst: "game"},
			{Name: "quit_game", Src: []string{"game"}, Dst: "home"},
		},
		fsm.Callbacks{
			"create_demo": func(e *fsm.Event) { StartDemo(&client) },
			"quit_demo": func(e *fsm.Event) {
				client.CurrentGame = nil
				client.UpdateTrace("quiting demo : " + e.FSM.Current())
			},
			"authenticate": func(e *fsm.Event) {
				client.UpdateTrace("authenticating : " + e.FSM.Current())
				GetServer().AddClient(&client)
			},
			"disconnect": func(e *fsm.Event) {
				client.UPTrace("disconnecting : " + e.FSM.Current())
				GetServer().CleanClient(GetServer().clients[client.Id])
				//GetServer().RemoveClient(&client)
				//client.Stop()
			},
			"join_lobby": func(e *fsm.Event) { client.UpdateTrace("joining lobby : " + e.FSM.Current()) },
			"create_lobby": func(e *fsm.Event) {
				client.UpdateTrace("creating lobby : " + e.FSM.Current())
			},
			"quit_lobby": func(e *fsm.Event) {
				client.CurrentLobby = nil
				client.UpdateTrace("quiting lobby : " + e.FSM.Current())
			},
			"join_game": func(e *fsm.Event) {
				client.CurrentLobby = nil
				client.UpdateTrace("joining game : " + e.FSM.Current())
			},
			"quit_game": func(e *fsm.Event) {
				client.CurrentGame = nil
				client.UpdateTrace("quiting game : " + e.FSM.Current())
			},
		},
	)

	return &client
}

func NewAiClient() *Client {
	var client Client
	client.Conn = nil
	client.State = nil
	ai := NewIA(&client)
	client.Ai = &ai
	return &client
}

func (client *Client) StartTrace() {
	var buffer bytes.Buffer
	buffer.WriteString("Client[")
	buffer.WriteString(strconv.Itoa(client.Id))
	buffer.WriteString("]")
	client.trace = buffer.String()
}
func (client *Client) UpdateTrace(trace string) {
	var buffer bytes.Buffer
	buffer.WriteString(client.trace)
	buffer.WriteString(trace)
	client.trace = buffer.String()
}
func (client *Client) PrintTrace() {
	log.Info(client.trace)
	client.trace = ""
}

func (client *Client) UPTrace(trace string) {
	client.UpdateTrace(trace)
	client.PrintTrace()
}
func (client *Client) Trace() string {
	return client.trace
}
func (client *Client) Tracing() bool {
	if client.trace != "" {
		return true
	}
	return false
}

func (client *Client) IsAi() bool {
	if client.Ai == nil {
		return false
	} else {
		return true
	}
}

func (client *Client) IsAuthenticated() bool {
	if client.Id == -1 {
		return false
	} else {
		return true
	}
}

func (client *Client) GameId() int {
	if client.CurrentGame == nil {
		return -1
	}
	for key := range client.CurrentGame.Clients {
		if client.CurrentGame.Clients[key] == client {
			return key
		}
	}
	return -1
}

func (client *Client) ObserverId() int {
	if client.CurrentGame == nil {
		return -1
	}
	for index, _ := range client.CurrentGame.Observers {
		if client.CurrentGame.Observers[index] == client {
			return index
		}
	}
	return -1
}

func (client *Client) Authenticate(auth AuthenticateJson) bool {
	marshalledAuth, err := json.Marshal(auth)
	if err != nil {
		log.Error("%s\n", err)
		return false
	}
	resp, response, err := ApiRequest("POST", "manager/authenticate_player", marshalledAuth)
	if err != nil {
		return false
	}

	if resp.StatusCode == 200 {
		client.Id = auth.PlayerId
		client.Pseudo = fmt.Sprintf("%s", response["pseudo"])
		if GetServer().reconnectClient(client) {
			log.Info("succefully reconnected Client[" + strconv.Itoa(GetServer().clients[client.Id].Id) + "] to state " + GetServer().clients[client.Id].State.Current())
		} else {
			client.State.Event("authenticate")
			client.UPTrace("success")
		}
		return true
	} else {
		client.UPTrace("failure:" + strconv.Itoa(auth.PlayerId) + ":" + auth.AccessToken + ":" + auth.Client)
		return false
	}
}

func (client *Client) Start() {
	if client.IsAi() {
		go client.Ai.Start()
		return
	}
	client.listening = true
	client.quitReader = make(chan int)
	//defer client.Conn.Close()
	for {
		mt, message, err := client.Conn.ReadMessage()
		if err != nil {
			log.Warn("read: ", err)
			client.listening = false
			// 1000 : websocket: close 1000 (normal)
			// 1001 : websocket: close 1001 (going away)
			// 1006 : websocket: close 1006 (abnormal closure): unexpected EOF
			// Si client going away ou abnormal closure on stop le reader, si timeout ça relit un message tous les 1 secondes
			if strings.Contains(err.Error(), "1001") || strings.Contains(err.Error(), "1006") || strings.Contains(err.Error(), "timeout") {
				return
			}
			if strings.Contains(err.Error(), "1000") {
				//TODO
			}
			//TODO on atterit ici et sa affiche websocket: close 1005 (no status)  lorsqu'un utilisateur ferme la fenetre ou a temporairement plus de réseau
			time.Sleep(1 * time.Second)
			if client.terminating {
				return
			}
		}
		if mt == websocket.TextMessage {
			request := Request{}
			json.Unmarshal(message, &request)
			request.Client = client
			client.StartTrace()
			request.Dispatch()
		}
	}
}

func (client *Client) StartWriter() {
	//request := Request{}
	retryLimit := 30
	pingInterval := time.Second * 1
	c := time.Tick(pingInterval)
	var sendTime time.Time
	client.quitWriter = make(chan int)
	client.dcByOtherClient = make(chan int)
	client.quitPing = make(chan int)
	//ping the client
	pingChan := make(chan int)
	go func() {
		for _ = range c {
			select {
			case _ = <-client.quitPing:
				return
			case <-time.After(pingInterval):
				// Si le client n'est plus le même on le deco
				if client.terminating || client.retry > retryLimit || (GetServer().clients[client.Id] != nil && GetServer().clients[client.Id].Conn != client.Conn) {
					return
				}
				pingChan <- 1
			}
		}
	}()

	client.Conn.SetPongHandler(func(test string) error {
		recieveTime := time.Now().Add(time.Second * 20)
		//log.Info("ms:", int64(recieveTime.Sub(sendTime)/time.Millisecond)) // ms: 100
		client.Ping = int64(recieveTime.Sub(sendTime) / time.Millisecond)
		if !client.listening {
			go client.Start()
		}
		client.Conn.SetReadDeadline(time.Now().Add(time.Second * 30))
		client.Conn.SetWriteDeadline(time.Now().Add(time.Second * 30))
		return client.Conn.SetReadDeadline(recieveTime)
	})

	for {
		if client.terminating {
			return
		}
		select {
			case request, more := <-client.RequestChannel:
				if more {
					if client.Tracing() {
						client.UpdateTrace("->Writer->Sending")
					}
					if client.retry > 0 {
						break
					}
					err := client.Conn.WriteMessage(websocket.TextMessage, request.Marshal())
					if err != nil {
						log.Warn("Client[" + strconv.Itoa(client.Id) + "] " + err.Error())
					} else if client.Tracing() {
						client.UpdateTrace("->Sent")
						client.PrintTrace()
					}
				} else {
					fmt.Println("terminating client")
					return
				}
			case _ = <-pingChan:
				sendTime = time.Now().Add(time.Second * 20)
				if err := client.Conn.WriteControl(websocket.PingMessage, []byte("test"), sendTime); err != nil {
					log.Warn("pinging error")
					client.retry += 1
					client.Ping = 1000
					log.Warn(client.retry)
					log.Warn(err)
					if client.retry > retryLimit {
						log.Info("Client[" + strconv.Itoa(client.Id) + "] is being removed from Server")
						GetServer().CleanClient(client)
					}
				} else if client.Ping < 1000 {
					client.retry = 0
					if !client.listening {
						go client.Start()
					}
				}
			case _ = <-client.quitWriter:
				return
			case _ = <-client.dcByOtherClient:
				req := NewRequest("disconnectedByOtherClient")
				client.Conn.WriteMessage(websocket.TextMessage, req.Marshal())
				return
		}
		if client.retry > retryLimit {
			return
		}
	}
}

func (client *Client) JoinLobby(lobby *Lobby) {
	client.CurrentLobby = lobby
	client.State.Event("join_lobby")
}

func (client *Client) LeaveLobby() {
	client.State.Event("quit_lobby")
}

func (client *Client) Stop() {
	client.terminating = true
	client.quitPing <- 1
	client.quitWriter <- 1
	client.listening = false
}

func (client *Client) StopAndSendDCByOtherClientRequest() {
	client.terminating = true
	client.quitPing <- 1
	client.dcByOtherClient <- 1
	client.listening = false
}

func (client *Client) stopReader() {
	//duration := time.Duration(10)*time.Second // Pause for 10 seconds
	client.terminating = true
}

func (client *Client) Restart() {
	client.terminating = false
	client.retry = 0
	go client.Start()
	go client.StartWriter()
}
