package model

type Message struct {
	Id             int             `json:"id"`
	ClientId             int        `json:"clientId"`
	HubId             int             `json:"hubId"`
	HolderType             string             `json:"holderType"`
	HolderId             int             `json:"holderId"`
	RecipientId             int             `json:"recipientId"`
	Message        string 	`json:"message"`
}