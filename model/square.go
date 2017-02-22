package model

import(
	. "../utils"
)

type Square struct {
	X        int  `json:"X"`
	Y        int  `json:"Y"`
	PlayerId *int `json:"playerId"`
}

func (square *Square) GetPlayerId() int {
	if square.PlayerId == nil {
		return -1
	} else {
		return *square.PlayerId
	}
}

func (square *Square) Equal(cube Cube) bool {
	if square.X == cube.X && square.Y == cube.Y {
		return true
	} else {
		return false
	}
}

func (square *Square) getDiagonalAuthorizedSquares (board *Board) []*Square{
	authorizedSquares := []*Square {}
	//si le square en haut à gauche existe et est vide
	if AllowedCoordinates(square.X-1, square.Y-1) && board.Squares[square.X-1][square.Y-1].IsEmpty() {
		//si ce square n'a pas de case adjacent appartenant au joueur
		if !board.Squares[square.X-1][square.Y-1].hasAdjacentSquaresWithPlayerId(board, square.GetPlayerId()){
			authorizedSquares = append(authorizedSquares, board.Squares[square.X-1][square.Y-1])
		}
	}
	//si le square en haut à droite existe et est vide
	if AllowedCoordinates(square.X+1, square.Y-1) && board.Squares[square.X+1][square.Y-1].IsEmpty() {
		//si ce square n'a pas de case adjacent appartenant au joueur
		if !board.Squares[square.X+1][square.Y-1].hasAdjacentSquaresWithPlayerId(board, square.GetPlayerId()){
			authorizedSquares = append(authorizedSquares, board.Squares[square.X+1][square.Y-1])
		}
	}
	//si le square en bas à gauche existe et est vide
	if AllowedCoordinates(square.X-1, square.Y+1) && board.Squares[square.X-1][square.Y+1].IsEmpty() {
		//si ce square n'a pas de case adjacent appartenant au joueur
		if !board.Squares[square.X-1][square.Y+1].hasAdjacentSquaresWithPlayerId(board, square.GetPlayerId()){
			authorizedSquares = append(authorizedSquares, board.Squares[square.X-1][square.Y+1])
		}
	}
	//si le square en bas à droite existe et est vide
	if AllowedCoordinates(square.X+1, square.Y+1) && board.Squares[square.X+1][square.Y+1].IsEmpty() {
		//si ce square n'a pas de case adjacent appartenant au joueur
		if !board.Squares[square.X+1][square.Y+1].hasAdjacentSquaresWithPlayerId(board, square.GetPlayerId()){
			authorizedSquares = append(authorizedSquares, board.Squares[square.X+1][square.Y+1])
		}
	}
	return authorizedSquares
}

func (square *Square) hasAdjacentSquaresWithPlayerId (board *Board, playerId int) bool {
	//si le square en haut appartient au joueur
	if AllowedCoordinates(square.X, square.Y-1) && board.Squares[square.X][square.Y-1].GetPlayerId() == square.GetPlayerId(){
		return true
	} else if AllowedCoordinates(square.X, square.Y+1) && board.Squares[square.X][square.Y+1].GetPlayerId() == square.GetPlayerId(){
		return true
	} else if AllowedCoordinates(square.X-1, square.Y) && board.Squares[square.X-1][square.Y].GetPlayerId() == square.GetPlayerId(){
		return true
	} else if AllowedCoordinates(square.X+1, square.Y) && board.Squares[square.X+1][square.Y].GetPlayerId() == square.GetPlayerId(){
		return true
	}
	return false
}



func (square *Square) IsEmpty() bool {
	if square.PlayerId == nil {
		return true
	} else {
		return false
	}
}