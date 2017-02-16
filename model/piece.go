package model

type Piece struct {
	Id       int     `json:"id"`
	Cubes    []Cube  `json:"cubes"`
	Origin   *Square `json:"origin"`
	Rotation string  `json:"rotation"`
	Flipped  bool    `json:"flipped"`
	PlayerId int     `json:"playerId"`
}

type PieceFactory struct {
	Id       int
	Rotation string
	Flipped  bool
}

func NewPieceFactory() *PieceFactory {
	var factory = new(PieceFactory)
	factory.Id = 0
	factory.Rotation = "N"
	factory.Flipped = false
	return factory
}

func (factory *PieceFactory) NewPiece() Piece {
	var piece Piece
	piece.Id = factory.Id
	piece.Rotation = factory.Rotation
	piece.Flipped = factory.Flipped
	factory.Id++
	return piece
}
