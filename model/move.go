package model

type Move struct {
	Turn     int    `json:"turn"`
	PlayerId int    `json:"player_id"`
	ClientId int    `json:"client_id"`
	Piece    *Piece `json:"piece"`
	Duration int    `json:"duration"`
}
