package model

import (
  "fmt"
)

type Player struct {
  Id int `json:"id"`
  Name string `json:"name"`
  Color string `json:"color"`
  Pieces [] Piece `json:"pieces"`
}

func (player *Player) initPieces() {
  for index,piece := range player.Pieces {
    player.Pieces[index].Player = player
    piece.Player = player
    // index is the index where we are
    // element is the element from someSlice for where we are
  }
    fmt.Println("init player pieces")
}