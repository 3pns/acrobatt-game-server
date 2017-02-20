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
	}
	fmt.Println("init player pieces")
}

func (player *Player) PlacePiece(piece Piece, board *Board) {
	if piece.Origin == nil {
		fmt.Println("Fatal Error piece has no Origin")
		return
	}
	player.Pieces[piece.Id].Origin = piece.Origin
	piece = player.Pieces[piece.Id]

	fmt.Println(piece.String())
	board.Squares[piece.Origin.X][piece.Origin.Y].PlayerId = &board.Players[*piece.PlayerId].Id
	fmt.Println("##### INBETWEEN #####")
	board.Players[*piece.PlayerId].Pieces[piece.Id].Origin = board.Squares[piece.Origin.X][piece.Origin.Y]
	fmt.Println("##### INAFTER #####")

	//1 - vérifier si on a le droit de placer la pièce
	//1.1 un cube est toujours dans la board
	//1.2 un cube n'est pas adjacent à un autre cube de la même couleur d'une autre pièce
	//2.2 un des cubes est dans la zone de départ ET/OU en diagonale d'un cube de la même couleur
	//piece.Rotation = "E"
	//2 - placer la pièce
	var projectedCubes []Cube
	var placementAuthorized = true
	var cubeOutOfBoard = false
	//var hasAtLeastACubeAtStartOrDiagonal = false
	fmt.Println("----- Plaçage d'une pièce -----")
	for _, cube := range piece.Cubes {
		var projectedCube = cube.Project(*piece.Origin, piece.Rotation, piece.Flipped) // on projete le cube dans l'espace = vrai position
		projectedCubes = append(projectedCubes, projectedCube) // on ajoute le cube à la liste des cube projeté => càd dire les vrais cases occupés par la pièces sur la board

		if projectedCube.X < 0 || projectedCube.X > 19 || projectedCube.Y < 0 || projectedCube.Y > 19 {
			cubeOutOfBoard = true
			return
		}

		fmt.Println(projectedCube, " : ")
		for _,startingCubes := range player.StartingCubes {
			fmt.Print(startingCubes)
			if startingCubes == projectedCube {
				//hasAtLeastACubeAtStartOrDiagonal = true
			}
		}
		if !placementAuthorized {
			return
		}
		
	}
	if placementAuthorized && !cubeOutOfBoard {
		fmt.Println("----- Placement Authorized -----")
		for _, cube := range projectedCubes {
			board.Squares[cube.X][cube.Y].PlayerId = piece.PlayerId
		}
	}
	board.PrintBoard()

}