package model

import (
	"../utils"
	"encoding/json"
	"fmt"
	"net"
	_"strconv"
)

type Board struct {
	Squares [20][20]*Square `json:"squares"`
	Pieces  []Piece         `json:"pieces"`
	Players []*Player       `json:"players"`
}

func (board *Board) PlacePiece(piece Piece) {
	if piece.Origin == nil {
		fmt.Println("Fatal Error piece has no Origin")
		return
	}

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
	var hasAtLeastACubeAtStartOrDiagonal = false
	fmt.Println("----- Plaçage d'une pièce -----")
	for _, cube := range piece.Cubes {
		var projectedCube = cube.Project(*piece.Origin, piece.Rotation, piece.Flipped) // on projete le cube dans l'espace = vrai position
		projectedCubes = append(projectedCubes, projectedCube)

		hasAtLeastACubeAtStartOrDiagonal = true
		if !placementAuthorized{
			return
		}
		
	}
	if placementAuthorized && hasAtLeastACubeAtStartOrDiagonal {
		for _, cube := range projectedCubes {
			board.Squares[cube.X][cube.Y].PlayerId = piece.PlayerId
		}
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

func (board *Board) InitPlayers() {
	//Joueur 0
	//copie des pieces modèle dans un nouveau slice pour le joueur
	var player0Pieces = make([]Piece, len(board.Pieces))
	copy(player0Pieces, board.Pieces)
	//génération des cases de départs
	var player0StartCubes = []Cube{}
	for i:=0; i<10; i++ {
		player0StartCubes = append(player0StartCubes, Cube{i,0})
	}
	for i:=0; i<10; i++ {
		player0StartCubes = append(player0StartCubes, Cube{0,i})
	}
	//création du joueur
	player0 := Player{0, "Joueur", "blue", player0Pieces, player0StartCubes}

	//Joueur 1
	var player1Pieces = make([]Piece, len(board.Pieces))
	copy(player1Pieces, board.Pieces)

	var player1StartCubes = []Cube{}
	for i:=0; i<10; i++ {
		player1StartCubes = append(player1StartCubes, Cube{i,19})
	}
	for i:=10; i<20; i++ {
		player1StartCubes = append(player1StartCubes, Cube{0,i})
	}
	player1 := Player{1, "AI-1", "green", player1Pieces, player1StartCubes}

	//Joueur 2
	var player2Pieces = make([]Piece, len(board.Pieces))
	copy(player2Pieces, board.Pieces)

	var player2StartCubes = []Cube{}
	for i:=10; i<20; i++ {
		player2StartCubes = append(player2StartCubes, Cube{i,0})
	}
	for i:=0; i<10; i++ {
		player2StartCubes = append(player2StartCubes, Cube{19,i})
	}
	player2 := Player{2, "AI-2", "yellow", player2Pieces, player2StartCubes}

	//Joueur 3
	var player3Pieces = make([]Piece, len(board.Pieces))
	copy(player3Pieces, board.Pieces)
	
	var player3StartCubes = []Cube{}
	for i:=10; i<20; i++ {
		player3StartCubes = append(player3StartCubes, Cube{i,19})
	}
	for i:=10; i<20; i++ {
		player3StartCubes = append(player3StartCubes, Cube{19,i})
	}
	player3 := Player{3, "AI-3", "red", player3Pieces, player3StartCubes}

	player0.Init()
	player1.Init()
	player2.Init()
	player3.Init()

	for index, _ := range player0.Pieces {
		fmt.Println("index :",index," player ID :",*player0.Pieces[index].PlayerId)
	}

	board.Players = []*Player{&player0, &player1, &player2, &player3}
}

func (board *Board) PrintBoard() {
	for i := 0; i < 20; i++ {
		utils.SetWhiteBackground()
		fmt.Print(" ")
		for j := 0; j < 20; j++ {
			if board.Squares[j][i].PlayerId == nil {
				utils.PrintBlack("▇ ")

			} else if *board.Squares[j][i].PlayerId == 0 {
				utils.PrintBlue("▇ ")

			} else if *board.Squares[j][i].PlayerId == 1 {
				utils.PrintGreen("▇ ")

			} else if *board.Squares[j][i].PlayerId == 2 {
				utils.PrintYellow("▇ ")

			}else {
				utils.PrintRed("▇ ")
			}
		}
		fmt.Print(" ")
		utils.SetBlackBackground()
		fmt.Println("")
	}
	utils.PrintReset()
}

func (board *Board) Refresh(conn net.Conn) {
	b, err := json.Marshal(board)
	if err != nil {
		fmt.Println(err)
	}
	conn.Write(b)
	conn.Write([]byte("\n"))
}
