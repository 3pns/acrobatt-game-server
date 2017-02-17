package model

type Square struct {
	X        int  `json:"X"`
	Y        int  `json:"Y"`
	PlayerId *int `json:"playerId"`
}
