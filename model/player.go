package model

import (
	"fmt"
)

type Player struct {
	Id     int     `json:"id"`
	Name   string  `json:"name"`
	Color  string  `json:"color"`
	Pieces []Piece `json:"pieces"`
}

func (player *Player) Init() {
	for index, _ := range player.Pieces {
		player.Pieces[index].PlayerId = player.Id
	}
	fmt.Println("init player pieces")
}
