package model

import (
	"fmt"
	"encoding/json"
)

type Player struct {
	Id     int     `json:"id"`
	Name   string  `json:"name"`
	Color  string  `json:"color"`
	Pieces []Piece `json:"pieces"`
	startingCubes	[]Cube
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
		fmt.Println("Fatal Error piece to place has no Origin")
		return
	}
	if player.Pieces[piece.Id].Origin != nil {
		fmt.Println("Fatal Error piece has already been used")
		return
	}
	player.Pieces[piece.Id].Origin = piece.Origin
	player.Pieces[piece.Id].Rotation = piece.Rotation
	piece = player.Pieces[piece.Id]

	fmt.Println(piece.String())
	
	fmt.Println("##### INAFTER #####")

	//1 - vérifier si on a le droit de placer la pièce
	//1.1 un cube est toujours dans la board
	//1.2 un cube n'est pas adjacent à un autre cube de la même couleur d'une autre pièce
	//2.2 un des cubes est dans la zone de départ ET/OU en diagonale d'un cube de la même couleur
	//piece.Rotation = "E"
	//2 - placer la pièce
	var projectedCubes []Cube
	var placementAuthorized = false
	var cubeOutOfBoard = false
	//var hasAtLeastACubeAtStartOrDiagonal = false
	fmt.Println("----- Plaçage d'une pièce -----")
	for _, cube := range piece.Cubes {
		var projectedCube = cube.Project(*piece.Origin, piece.Rotation, piece.Flipped) // on projete le cube dans l'espace = vrai position
		projectedCubes = append(projectedCubes, projectedCube) // on ajoute le cube à la liste des cube projeté => càd dire les vrais cases occupés par la pièces sur la board

		if projectedCube.X < 0 || projectedCube.X > 19 || projectedCube.Y < 0 || projectedCube.Y > 19 {
			fmt.Println("SIGSEV Placement Out of Board Exception")
			cubeOutOfBoard = true
			return
		}
		if player.IsAStartingCube(projectedCube){
			placementAuthorized = true
			fmt.Println("Placement Authorized cuz Starting Cube  :", projectedCube)
		}
		// si le cube en bas à gauche est dans la board et appartient au joueur le placement est autorisé
		if projectedCube.X-1 > 0 && projectedCube.X-1 < 20 && projectedCube.Y+1 > 0 && projectedCube.Y+1 < 20 {
			if board.Squares[projectedCube.X-1][projectedCube.Y+1].GetPlayerId() == player.Id{
				placementAuthorized = true
			}
		}
	}
	if !placementAuthorized {
		fmt.Println("----- BADDIES Placement Unauthorized Exception -----")
		return
	}
	fmt.Println("ALLO FRANCIS")
	fmt.Println("etat des variables, placement :", placementAuthorized, "cube out of board:", !cubeOutOfBoard)
	if placementAuthorized && !cubeOutOfBoard {
		fmt.Println("----- Placement Authorized -----")
		for _, cube := range projectedCubes {
			board.Squares[cube.X][cube.Y].PlayerId = piece.PlayerId
		}
		//board.Squares[piece.Origin.X][piece.Origin.Y].PlayerId = &board.Players[*piece.PlayerId].Id
	}
}

func (player *Player) IsAStartingCube(cube Cube) bool {
	for _,startingCube := range player.startingCubes {
		if startingCube.Equal(cube) {
			return true
		}
	}
	return false
}

func (player Player) String() string {
	b, err := json.Marshal(player)
	if err != nil {
		fmt.Println(err)
	}
	return string(b)
}

func (player Player) PrintStartingCubes() string {
	b, err := json.Marshal(player.startingCubes)
	if err != nil {
		fmt.Println(err)
	}
	return string(b)
}

func (player *Player) AppendStartingCube(cube Cube) {
	player.startingCubes = append(player.startingCubes, cube)
}