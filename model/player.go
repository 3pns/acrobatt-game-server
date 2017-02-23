package model

import (
	"encoding/json"
	"fmt"
	"math/rand"
	. "../utils"
)

type Player struct {
	Id            int     `json:"id"`
	Name          string  `json:"name"`
	Color         string  `json:"color"`
	Pieces        []Piece `json:"pieces"`
	startingSquares []*Square
	squares       []*Square
}

func (player *Player) Init() {
	fmt.Println("playerId in playerInit :", player.Id)
	for index, _ := range player.Pieces {
		player.Pieces[index].PlayerId = &player.Id
	}
	fmt.Println("init player pieces")
}

func (player *Player) PlacePiece(piece Piece, board *Board, simulation bool) bool {
	if piece.Origin == nil {
		fmt.Println("Fatal Error piece to place has no Origin")
		return false
	}
	if player.Pieces[piece.Id].Origin != nil {
		fmt.Println("Fatal Error piece has already been used")
		return false
	}
	piece.PlayerId = &player.Id
	//1 - vérifier si on a le droit de placer la pièce
	//1.1 un cube est toujours dans la board
	//1.2 un cube n'est pas adjacent à un autre cube de la même couleur d'une autre pièce
	//2.2 un des cubes est dans la zone de départ ET/OU en diagonale d'un cube de la même couleur
	//piece.Rotation = "E"
	//2 - placer la pièce
	var projectedCubes []Cube
	var placementAuthorized = false
	//var hasAtLeastACubeAtStartOrDiagonal = false
	fmt.Println("----- Plaçage d'une pièce -----")
	for _, cube := range piece.Cubes {
		var projectedCube = cube.Project(*piece.Origin, piece.Rotation, piece.Flipped) // on projete le cube dans l'espace = vrai position
		projectedCubes = append(projectedCubes, projectedCube)                         // on ajoute le cube à la liste des cube projeté => càd dire les vrais cases occupés par la pièces sur la board
		fmt.Println(projectedCubes)
		//si le cube est en dehors de la board le placement est interdit
		if projectedCube.X < 0 || projectedCube.X > 19 || projectedCube.Y < 0 || projectedCube.Y > 19 {
			fmt.Println("SIGSEV Placement Out of Board Exception")
			placementAuthorized = false
			return false
		}
		//si le cube occupe un square occupé le placement est interdit
		if board.Squares[projectedCube.X][projectedCube.Y].PlayerId != nil {
			fmt.Println("StackOverflow Board Exception le square est déjà occupé")
			placementAuthorized = false
			return false
		}
		// si le cube en bas est dans la board et appartient au joueur le placement est interdit
		if projectedCube.Y+1 > 0 && projectedCube.Y+1 < 20 {
			if board.Squares[projectedCube.X][projectedCube.Y+1].GetPlayerId() == player.Id {
				fmt.Println("Placement Unauthorized Exceptio.cuz cube en bas appartient au joueur")
				placementAuthorized = false
				return false
			}
		}
		// si le cube en haut est dans la board et appartient au joueur le placement est interdit
		if projectedCube.Y-1 > 0 && projectedCube.Y-1 < 20 {
			if board.Squares[projectedCube.X][projectedCube.Y-1].GetPlayerId() == player.Id {
				fmt.Println("Placement Unauthorized Exceptio.cuz cube en haut appartient au joueur")
				placementAuthorized = false
				return false
			}
		}
		// si le cube à gauche est dans la board et appartient au joueur le placement est interdit
		if projectedCube.X-1 > 0 && projectedCube.X-1 < 20 {
			if board.Squares[projectedCube.X-1][projectedCube.Y].GetPlayerId() == player.Id {
				fmt.Println("Placement Unauthorized Exceptio.cuz cube à gauche appartient au joueur")
				placementAuthorized = false
				return false
			}
		}
		// si le cube à droite est dans la board et appartient au joueur le placement est interdit
		if projectedCube.X+1 > 0 && projectedCube.X+1 < 20 {
			if board.Squares[projectedCube.X+1][projectedCube.Y].GetPlayerId() == player.Id {
				fmt.Println("Placement Unauthorized Exceptio.cuz cube à gauche appartient au joueur")
				placementAuthorized = false
				return false
			}
		}
		// si le cube est le cube de départ du joueur le placement est autorisé
		if player.IsAStartingCube(projectedCube) {
			placementAuthorized = true
			fmt.Println("Placement Authorized cuz Starting Cube  :", projectedCube)
		}
		// si le cube en bas à gauche est dans la board et appartient au joueur le placement est autorisé
		if projectedCube.X-1 >= 0 && projectedCube.X-1 < 20 && projectedCube.Y+1 >= 0 && projectedCube.Y+1 < 20 {
			if board.Squares[projectedCube.X-1][projectedCube.Y+1].GetPlayerId() == player.Id {
				fmt.Println("Placement Authorized cuz cube en bas à gauche")
				placementAuthorized = true
			}
		}
		// si le cube en bas à droite est dans la board et appartient au joueur le placement est autorisé
		if projectedCube.X+1 >= 0 && projectedCube.X+1 < 20 && projectedCube.Y+1 >= 0 && projectedCube.Y+1 < 20 {
			if board.Squares[projectedCube.X+1][projectedCube.Y+1].GetPlayerId() == player.Id {
				fmt.Println("Placement Authorized cuz cube en bas à droite")
				placementAuthorized = true
			}
		}
		// si le cube en haut à gauche est dans la board et appartient au joueur le placement est autorisé
		if projectedCube.X-1 >= 0 && projectedCube.X-1 < 20 && projectedCube.Y-1 >= 0 && projectedCube.Y-1 < 20 {
			if board.Squares[projectedCube.X-1][projectedCube.Y-1].GetPlayerId() == player.Id {
				fmt.Println("Placement Authorized cuz cube en haut à gauche")
				placementAuthorized = true
			}
		}
		// si le cube en haut à droite est dans la board et appartient au joueur le placement est autorisé
		if projectedCube.X+1 >= 0 && projectedCube.X+1 < 20 && projectedCube.Y-1 >= 0 && projectedCube.Y-1 < 20 {
			if board.Squares[projectedCube.X+1][projectedCube.Y-1].GetPlayerId() == player.Id {
				fmt.Println("Placement Authorized cuz cube en haut à droite")
				placementAuthorized = true
			}
		}
	}
	if !placementAuthorized {
		fmt.Println("----- BADDIES Placement Unauthorized Exception -----")
		return false
	} else {
		fmt.Println("----- Placement Authorized -----")
		if simulation {
			return true
		}
		for _, cube := range projectedCubes {
			board.Squares[cube.X][cube.Y].PlayerId = piece.PlayerId
			player.squares = append(player.squares, board.Squares[cube.X][cube.Y])
		}
		player.Pieces[piece.Id].Origin = piece.Origin
		player.Pieces[piece.Id].Rotation = piece.Rotation
		player.Pieces[piece.Id].Flipped = piece.Flipped
		return true
	}
}

func (player *Player) IsAStartingCube(cube Cube) bool {
	for _, startingSquare := range player.startingSquares {
		if startingSquare.Equal(cube) {
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

func (player Player) PrintStartingSquares() string {
	b, err := json.Marshal(player.startingSquares)
	if err != nil {
		fmt.Println(err)
	}
	return string(b)
}

func (player *Player) PlaceRandomPieceWithIAEasy(board *Board, simulation bool) bool {
	//on récupère les pièces restantes à placer
	var remainingPieces = [] *Piece{}
	for index, piece := range player.Pieces {
		if piece.Origin == nil {
			remainingPieces = append(remainingPieces, &player.Pieces[index])
		}
	}
	//on récupère un index de piece au hasard ssi il reste des pièces à placer
	var index int
	var targetSquares = [] *Square {}
	if len(remainingPieces) > 0{
		index = rand.Intn(len(remainingPieces))	
	} else {
		//le joueur a placé toutes ses pièces !
		return false
	}
	//si le joueur a encore toutes ses pièces le square cible est son square de départ
	if len(remainingPieces) == 21{
		targetSquares =  player.startingSquares
	} else {
		//sinon à partir des squares appartenant au joueur on récupère les squares ou l'IA peut poser une pièce
		fmt.Println("generating authorize squares for player",player.Id)
		for _, square := range player.squares {
			fmt.Println("checking first player Square:",square.X,",",square.Y)
			targetSquares = append(targetSquares, square.getDiagonalAuthorizedSquares(board)...)
		} 
	}
	//on essaye de placer toutes les pièces
	tryagain:
	fmt.Println("remainingPieces: ", len(remainingPieces))
	fmt.Println("playerSquares: ", player.squares)
	fmt.Println("targetSquares: ", targetSquares)
	index = rand.Intn(len(remainingPieces))
	piece := remainingPieces[index]
	if player.TryPlacePieceOnSquares(board, piece, targetSquares, simulation){
		return true
	} else if len(remainingPieces) > 1 {
		//on enlève la pièce du slience
		remainingPieces = append(remainingPieces[:index], remainingPieces[index+1:]...)
		goto tryagain
	}
	
	//TODO Remove targetSquares duplicates in slice !
	//le joueur ne peut placer aucune pièce !
	return false	
}

func (player *Player) TryPlacePieceOnSquares(board *Board, piece *Piece, squares [] *Square, simulation bool) bool {
	for _, square := range squares {
		//essayer les 8 rotation/coté possible
		piece.Flipped = false
		piece.Rotation = "N"
		//essayeer tous les positionnements de la pièce avec cette rotation/coté sur le square
		if player.TryPlacePieceOnSquareWithOrientation(board, *piece, square, simulation ){
			//on renvoit true si la pièce a été placé
			return true
		}
		piece.Rotation = "E"
		if player.TryPlacePieceOnSquareWithOrientation(board, *piece, square, simulation ){
			return true
		}
		piece.Rotation = "S"
		if player.TryPlacePieceOnSquareWithOrientation(board, *piece, square, simulation ){
			return true
		}
		piece.Rotation = "W"
		if player.TryPlacePieceOnSquareWithOrientation(board, *piece, square, simulation ){
			return true
		}
		piece.Flipped = true
		if player.TryPlacePieceOnSquareWithOrientation(board, *piece, square, simulation ){
			return true
		}
		piece.Rotation = "N"
		if player.TryPlacePieceOnSquareWithOrientation(board, *piece, square, simulation ){
			return true
		}
		piece.Rotation = "E"
		if player.TryPlacePieceOnSquareWithOrientation(board, *piece, square, simulation ){
			return true
		}
		piece.Rotation = "S"
		if player.TryPlacePieceOnSquareWithOrientation(board, *piece, square, simulation ){
			return true
		}
		piece.Rotation = "W"
		if player.TryPlacePieceOnSquareWithOrientation(board, *piece, square, simulation ){
			return true
		}
	}
	//false si la pièce n'a pas été placé
	return false
}

//essaye tous les positionnements de la pièce avec cette rotation/coté sur le square
func (player *Player) TryPlacePieceOnSquareWithOrientation(board *Board, piece Piece, square *Square, simulation bool) bool {
	for _,cube := range	piece.Cubes {
		if AllowedCoordinates(square.X -cube.X, square.Y -cube.Y){
			piece.Origin = board.Squares[square.X -cube.X][square.Y -cube.Y]
			if player.PlacePiece(piece, board, simulation){
				return true
			} else {
				piece.Origin = nil
			}
		}
		
	}
	return false
}

func (player *Player) HasPlaceabePieces(board *Board) bool{
	return player.PlaceRandomPieceWithIAEasy(board, true)
}