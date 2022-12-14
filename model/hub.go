package model

import (
	"encoding/json"
	_ "github.com/Sirupsen/logrus"
)

type Hub struct {
	Id             int             `json:"id"`
	RequestChannel chan Request    `json:"-"`
	Clients        map[int]*Client `json:"-"`
	HolderType     string          `json:"-"`
	HolderId       int             `json:"-"`
}

type HubFactory struct {
	Id int
}

func NewHubFactory() *HubFactory {
	var factory = new(HubFactory)
	factory.Id = 1
	return factory
}

func (factory *HubFactory) NewHub() *Hub {
	var hub Hub
	hub.Id = factory.Id
	factory.Id++
	hub.RequestChannel = make(chan Request, 100)
	return &hub
}

func (hub *Hub) Start() {
	for {
		request, more := <-hub.RequestChannel
		client := request.Client
		if more {
			if hub.Clients[client.Id] == nil {
				return
			}
			message := Message{}
			json.Unmarshal(request.Data, &message)
			message.SenderId = client.Id
			message.SenderPseudo = client.Pseudo
			message.HolderType = hub.HolderType
			message.HolderId = hub.HolderId
			if message.Message == "" {
				continue
			}
			if request.Type == "BroadcastMessage" {
				var req = NewRequest("BroadcastMessage")
				req.MarshalData(message)
				for index, _ := range hub.Clients {
					if hub.Clients[index] != nil && hub.Clients[index].RequestChannel != nil {
						hub.Clients[index].RequestChannel <- req
					}
				}
			} else if request.Type == "SendMessage" {
				if message.RecipientId > 0 && hub.Clients[message.RecipientId] != nil && hub.Clients[message.RecipientId].RequestChannel != nil {
					message.RecipientPseudo = hub.Clients[message.RecipientId].Pseudo
					var req = NewRequest("SendMessage")
					req.MarshalData(message)
					hub.Clients[message.RecipientId].RequestChannel <- req
				}
			}
		} else {
			hub.RequestChannel = nil
			return
		}
	}
}

func (hub *Hub) broadcastRequest(request *Request) {
	for index, _ := range hub.Clients {
		hub.Clients[index].RequestChannel <- *request
	}
}

func (hub *Hub) Stop() {
	close(hub.RequestChannel)
}
