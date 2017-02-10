package model

import (
	"../utils"
	"fmt"
	"strconv"
)

type Board struct {
	Squares [20][20]*Square `json:"squares"`
	Pieces  []Piece         `json:"pieces"`
	Players []*Player       `json:"pieces"`
}

func (board *Board) PlacePiece(piece *Piece, square *Square) {
	board.Squares[square.X][square.Y].Empty = false
	piece.Origin = board.Squares[square.X][square.Y]

	//1 - vérifier si on a le droit de placer la pièce
	//2 - placer la pièce
	fmt.Println("----- Plaçage d'une pièce -----")
	for _, cube := range piece.Cubes {
		var xFactor = 1
		var yFactor = 1
		var xBoardValue = cube.X
		var yBoardValue = cube.Y
		if piece.Flipped {
			xFactor = -xFactor
		}
		fmt.Println("Avant : x :" + strconv.Itoa(xBoardValue) + ", y: " + strconv.Itoa(yBoardValue))
		if piece.Rotation == "W" {
			xBoardValue = -cube.Y
			yBoardValue = -cube.X
		} else if piece.Rotation == "N" {
			xBoardValue = -cube.X
			yBoardValue = -cube.Y
		} else if piece.Rotation == "E" {
			xBoardValue = cube.Y
			yBoardValue = cube.X
		}
		fmt.Println("Apres : x :" + strconv.Itoa(xBoardValue) + ", y: " + strconv.Itoa(yBoardValue))
		board.Squares[square.X+xFactor*xBoardValue][square.Y+yFactor*yBoardValue].Empty = false
	}
}

func (board *Board) InitBoard() {
	fmt.Println("initializing board")
	for i := 0; i < 20; i++ {
		for j := 0; j < 20; j++ {
			board.Squares[i][j] = &Square{i, j, true}
		}
	}
	fmt.Println("board initialized with success !\n")
}

func (board *Board) InitPieces() {
	fmt.Println("generating pieces")

  var factory = NewPieceFactory()

	var piece = factory.NewPiece()
	piece.Cubes = []Cube{Cube{0,0}}
	board.Pieces = append(board.Pieces, piece)

	var piece1 = factory.NewPiece()
	piece1.Cubes = []Cube{Cube{0,0}, Cube{0,1}}
	board.Pieces = append(board.Pieces, piece1)

	var piece2 = factory.NewPiece()
	piece2.Cubes = []Cube{Cube{0,0}, Cube{0,1}, Cube{1,0}}
	board.Pieces = append(board.Pieces, piece2)

	var piece3 = factory.NewPiece()
	piece3.Cubes = []Cube{Cube{0,0}, Cube{0,1}, Cube{0,2}}
	board.Pieces = append(board.Pieces, piece3)

	var piece4 = factory.NewPiece()
	piece4.Cubes = []Cube{Cube{0,0}, Cube{1,0}, Cube{0,1}, Cube{0,2}}
	board.Pieces = append(board.Pieces, piece4)

	fmt.Println("pieces generated with success !\n")
}

func (board *Board) PrintBoard() {
	for i := 0; i < 20; i++ {
		utils.SetWhiteBackground()
		fmt.Print(" ")
		for j := 0; j < 20; j++ {

			if board.Squares[j][i].Empty == true {
				utils.PrintBlack("▇ ")

			} else {
				utils.PrintRed("▇ ")
			}
		}
		fmt.Print(" ")
		utils.SetBlackBackground()
		fmt.Println("")
	}
	utils.PrintReset()
}
