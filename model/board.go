package model

import (
	"../utils"
	"fmt"
	"strconv"
)

type Board struct {
	Squares [20][20]*Square `json:"squares"`
	Pieces  []Piece         `json:"pieces"`
	Players []*Player       `json:"players"`
}

func (board *Board) PlacePiece(piece *Piece, square *Square) {
	board.Squares[square.X][square.Y].PlayerId = &piece.PlayerId
	piece.Origin = board.Squares[square.X][square.Y]

	//1 - vérifier si on a le droit de placer la pièce
	//piece.Rotation = "E"
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
		board.Squares[square.X+xFactor*xBoardValue][square.Y+yFactor*yBoardValue].PlayerId = &piece.PlayerId
	}
}

func (board *Board) InitBoard() {
	fmt.Println("initializing board")
	for i := 0; i < 20; i++ {
		for j := 0; j < 20; j++ {
			board.Squares[i][j] = &Square{i, j, nil}
		}
	}
	fmt.Println("board initialized with success !\n")
}

func (board *Board) InitPieces() {
	fmt.Println("generating pieces")

	var factory = NewPieceFactory()

	//1 cube
	var piece = factory.NewPiece()
	piece.Cubes = []Cube{Cube{0, 0}}
	board.Pieces = append(board.Pieces, piece)

	//2 cubes
	var piece1 = factory.NewPiece()
	piece1.Cubes = []Cube{Cube{0, 0}, Cube{0, 1}}
	board.Pieces = append(board.Pieces, piece1)

	//3 cubes
	var piece2 = factory.NewPiece()
	piece2.Cubes = []Cube{Cube{0, 0}, Cube{0, 1}, Cube{1, 0}}
	board.Pieces = append(board.Pieces, piece2)

	var piece3 = factory.NewPiece()
	piece3.Cubes = []Cube{Cube{0, 0}, Cube{0, 1}, Cube{0, 2}}
	board.Pieces = append(board.Pieces, piece3)

	//4 cubes
	var piece4 = factory.NewPiece()
	piece4.Cubes = []Cube{Cube{0, 0}, Cube{0, 1}, Cube{0, 2}, Cube{0, 3}}
	board.Pieces = append(board.Pieces, piece4)

	var piece5 = factory.NewPiece()
	piece5.Cubes = []Cube{Cube{0, 0}, Cube{1, 0}, Cube{0, 1}, Cube{0, 2}}  
	board.Pieces = append(board.Pieces, piece5)

	var piece6 = factory.NewPiece()
	piece6.Cubes = []Cube{Cube{0, 0}, Cube{1, 0}, Cube{2, 0}, Cube{1, 1}}  
	board.Pieces = append(board.Pieces, piece6)

	var piece7 = factory.NewPiece()
	piece7.Cubes = []Cube{Cube{0, 0}, Cube{1, 0}, Cube{0, 1}, Cube{1, 1}}  
	board.Pieces = append(board.Pieces, piece7)

	var piece8 = factory.NewPiece()
	piece8.Cubes = []Cube{Cube{0, 0}, Cube{1, 0}, Cube{1, 1}, Cube{2, 1}}  
	board.Pieces = append(board.Pieces, piece8)

	// 5 cubes
	var piece9 = factory.NewPiece()
	piece9.Cubes = []Cube{Cube{0, 0}, Cube{0, 1}, Cube{0, 2}, Cube{0, 3}, Cube{0, 4}}  
	board.Pieces = append(board.Pieces, piece9)

	var piece10 = factory.NewPiece()
	piece10.Cubes = []Cube{Cube{0, 0}, Cube{0, 1}, Cube{0, 2}, Cube{0, 3}, Cube{-1, 3}}  
	board.Pieces = append(board.Pieces, piece10)

	var piece11 = factory.NewPiece()
	piece11.Cubes = []Cube{Cube{0, 0}, Cube{0, 1}, Cube{0, 2}, Cube{-1, 2}, Cube{-1, 3}}  
	board.Pieces = append(board.Pieces, piece11)

	var piece12 = factory.NewPiece()
	piece12.Cubes = []Cube{Cube{0, 0}, Cube{0, 1}, Cube{0, 2}, Cube{-1, 1}, Cube{-1, 2}}  
	board.Pieces = append(board.Pieces, piece12)

	var piece13 = factory.NewPiece()
	piece13.Cubes = []Cube{Cube{0, 0}, Cube{1, 0}, Cube{1, 1}, Cube{1, 2}, Cube{0, 2}}  
	board.Pieces = append(board.Pieces, piece13)

	var piece14 = factory.NewPiece()
	piece14.Cubes = []Cube{Cube{0, 0}, Cube{0, 1}, Cube{0, 2}, Cube{0, 3}, Cube{1, 1}}  
	board.Pieces = append(board.Pieces, piece14)

	var piece15 = factory.NewPiece()
	piece15.Cubes = []Cube{Cube{0, 0}, Cube{0, 1}, Cube{0, 2}, Cube{1, 2}, Cube{-1, 2}}  
	board.Pieces = append(board.Pieces, piece15)

	var piece16 = factory.NewPiece()
	piece16.Cubes = []Cube{Cube{0, 0}, Cube{0, 1}, Cube{0, 2}, Cube{1, 2}, Cube{2, 2}}  
	board.Pieces = append(board.Pieces, piece16)

	var piece17 = factory.NewPiece()
	piece17.Cubes = []Cube{Cube{0, 0}, Cube{1, 0}, Cube{1, 1}, Cube{2, 1}, Cube{2, 2}}  
	board.Pieces = append(board.Pieces, piece17)

	var piece18 = factory.NewPiece()
	piece18.Cubes = []Cube{Cube{0, 0}, Cube{0, 1}, Cube{1, 1}, Cube{2, 1}, Cube{2, 2}}  
	board.Pieces = append(board.Pieces, piece18)

	var piece19 = factory.NewPiece()
	piece19.Cubes = []Cube{Cube{0, 0}, Cube{0, 1}, Cube{1, 1}, Cube{2, 1}, Cube{1, 2}}  
	board.Pieces = append(board.Pieces, piece19)

	var piece20 = factory.NewPiece()
	piece20.Cubes = []Cube{Cube{0, 0}, Cube{0, 1}, Cube{1, 0}, Cube{-1, 0}, Cube{0, -1}}  
	board.Pieces = append(board.Pieces, piece20)

	fmt.Println("pieces generated with success !\n")
}

func (board *Board) PrintBoard() {
	for i := 0; i < 20; i++ {
		utils.SetWhiteBackground()
		fmt.Print(" ")
		for j := 0; j < 20; j++ {

			if board.Squares[j][i].PlayerId == nil {
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
