package model

type Piece struct {
  Id int `json:"id"`
  Cubes []Cube `json:"cubes"`
  Origin *Square `json:"origin"`
  Rotation string `json:"rotation"`
  Flipped bool `json:"flipped"`
  Player *Player `json:"player"`
}