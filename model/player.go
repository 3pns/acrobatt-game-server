package model

import (
	"fmt"
)

type Player struct {
	Id     int     `json:"id"`
	Name   string  `json:"name"`
	Color  string  `json:"color"`
	Pieces []Piece `json:"pieces"`
	StartingCubes	[]Cube `json:"-"`
}

func (player *Player) Init() {
	fmt.Println("playerId in playerInit :",player.Id)
	for index, _ := range player.Pieces {
		player.Pieces[index].PlayerId = &player.Id
		fmt.Println("index :",index," player ID :",*player.Pieces[index].PlayerId)
	}
	fmt.Println("init player pieces")
}
