package model

type Square struct {
	X        int  `json:"x"`
	Y        int  `json:"y"`
	PlayerId *int `json:"playerId"`
}
