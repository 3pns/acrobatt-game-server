package model

import (
	. "../utils"
	"fmt"
	"strconv"
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

func (square *Square) getDiagonalAuthorizedSquares(board *Board) []*Square {
	authorizedSquares := []*Square{}
	//si le square en haut à gauche existe et est vide
	if AllowedCoordinates(square.X-1, square.Y-1) && board.Squares[square.X-1][square.Y-1].IsEmpty() {
		//si ce square n'a pas de case adjacent appartenant au joueur
		if !board.Squares[square.X-1][square.Y-1].hasAdjacentSquaresWithPlayerId(board, square.GetPlayerId()) {
			authorizedSquares = append(authorizedSquares, board.Squares[square.X-1][square.Y-1])
		}
	}
	//si le square en haut à droite existe et est vide
	if AllowedCoordinates(square.X+1, square.Y-1) && board.Squares[square.X+1][square.Y-1].IsEmpty() {
		//si ce square n'a pas de case adjacent appartenant au joueur
		if !board.Squares[square.X+1][square.Y-1].hasAdjacentSquaresWithPlayerId(board, square.GetPlayerId()) {
			authorizedSquares = append(authorizedSquares, board.Squares[square.X+1][square.Y-1])
		}
	}
	//si le square en bas à gauche existe et est vide
	if AllowedCoordinates(square.X-1, square.Y+1) && board.Squares[square.X-1][square.Y+1].IsEmpty() {
		//si ce square n'a pas de case adjacent appartenant au joueur
		if !board.Squares[square.X-1][square.Y+1].hasAdjacentSquaresWithPlayerId(board, square.GetPlayerId()) {
			authorizedSquares = append(authorizedSquares, board.Squares[square.X-1][square.Y+1])
		}
	}
	//si le square en bas à droite existe et est vide
	if AllowedCoordinates(square.X+1, square.Y+1) && board.Squares[square.X+1][square.Y+1].IsEmpty() {
		//si ce square n'a pas de case adjacent appartenant au joueur
		if !board.Squares[square.X+1][square.Y+1].hasAdjacentSquaresWithPlayerId(board, square.GetPlayerId()) {
			authorizedSquares = append(authorizedSquares, board.Squares[square.X+1][square.Y+1])
		}
	}
	return authorizedSquares
}

func (square *Square) hasAdjacentSquaresWithPlayerId(board *Board, playerId int) bool {
	//si le square en haut appartient au joueur
	fmt.Print("square(", square.X, ",", square.Y, ") : ")
	if AllowedCoordinates(square.X, square.Y-1) && board.Squares[square.X][square.Y-1].GetPlayerId() == playerId {
		fmt.Println("batard : ", square.X, ",", strconv.Itoa(square.Y-1))
		fmt.Println("pas bon a cause de square du bas ")
		return true
	} else if AllowedCoordinates(square.X, square.Y+1) && board.Squares[square.X][square.Y+1].GetPlayerId() == playerId {
		fmt.Println("pas bon a cause de square du haut ")
		return true
	} else if AllowedCoordinates(square.X-1, square.Y) && board.Squares[square.X-1][square.Y].GetPlayerId() == playerId {
		fmt.Println("pas bon a cause de square de gauche ")
		return true
	} else if AllowedCoordinates(square.X+1, square.Y) && board.Squares[square.X+1][square.Y].GetPlayerId() == playerId {
		fmt.Println("pas bon a cause de square de droite ")
		return true
	}
	fmt.Println("pas de case adjacente appartenant au joueur ")
	return false
}

func (square *Square) IsEmpty() bool {
	if square.PlayerId == nil {
		return true
	} else {
		return false
	}
}
