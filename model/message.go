package model

type Message struct {
	Id             int             `json:"id"`
	ClientId             int        `json:"client_id"`
	HubId             int             `json:"hub_id"`
	HolderType             string             `json:"holder_type"`
	HolderId             int             `json:"holder_id"`
	RecipientId             int             `json:"recipient_id"`
	Message        string 	`json:"message"`
}