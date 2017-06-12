package model

type Message struct {
	Id             int             `json:"id"`
	SenderId             int        `json:"senderId"`
  SenderPseudo             string        `json:"senderPseudo"`
	HubId             int             `json:"hubId"`
	HolderType             string             `json:"holderType"`
	HolderId             int             `json:"holderId"`
	RecipientId             int             `json:"recipientId"`
  RecipientPseudo             string             `json:"recipientPseudo"`
	Message        string 	`json:"message"`
}