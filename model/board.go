package model

import (
	"../utils"
	"encoding/json"
	"fmt"
	"net"
	"sort"
	_ "strconv"
)

type Board struct {
	Squares    [20][20]*Square `json:"squares"`
	Pieces     []Piece         `json:"pieces"`
	Players    []*Player       `json:"players"`
	PlayerTurn *Player         `json:"playerTurn"`
	Turn       int             `json:"turn"`
	GameId     int             `json:"gameId"`
}

func (board *Board) InitBoard() {
	fmt.Println("initializing board")
	for i := 0; i < 20; i++ {
		for j := 0; j < 20; j++ {
			board.Squares[i][j] = &Square{i, j, nil}
		}
	}
	board.Turn = 1
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
	var player0Pieces = make([]Piece, len(board.Pieces))
	copy(player0Pieces, board.Pieces)

	var player0StartSquares = []*Square{board.Squares[0][0]}
	player0 := NewPlayer(0, "Joueur", "blue", player0Pieces, player0StartSquares)

	//Joueur 1
	var player1Pieces = make([]Piece, len(board.Pieces))
	copy(player1Pieces, board.Pieces)

	var player1StartSquares = []*Square{board.Squares[0][19]}
	player1 := NewPlayer(1, "AI-1", "green", player1Pieces, player1StartSquares)

	//Joueur 2
	var player2Pieces = make([]Piece, len(board.Pieces))
	copy(player2Pieces, board.Pieces)

	var player2StartSquares = []*Square{board.Squares[19][0]}
	player2 := NewPlayer(2, "AI-2", "yellow", player2Pieces, player2StartSquares)

	//Joueur 3
	var player3Pieces = make([]Piece, len(board.Pieces))
	copy(player3Pieces, board.Pieces)

	var player3StartSquares = []*Square{board.Squares[19][19]}
	player3 := NewPlayer(3, "AI-3", "red", player3Pieces, player3StartSquares)

	player0.Init()
	player1.Init()
	player2.Init()
	player3.Init()

	board.Players = []*Player{player0, player1, player2, player3}
	board.PlayerTurn = player0
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

func (board *Board) Refresh(conn net.Conn) {
	b, err := json.Marshal(board)
	if err != nil {
		fmt.Println(err)
	}
	conn.Write(b)
	conn.Write([]byte("\n"))
}

func (board *Board) SquareExistsAndBelongsTo(x int, y int, player Player) bool {
	if x >= 0 && x < 20 && y >= 0 && y < 20 {
		if board.Squares[x][y].GetPlayerId() == player.Id {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func (board *Board) NextTurn() {
	if board.Players[len(board.Players)-1] == board.PlayerTurn {
		board.Turn++
	}
	for i := 1; i <= len(board.Players); i++ {
		if board.Players[((board.PlayerTurn.Id + i) % len(board.Players))].HasPlaceabePieces(board) {
			board.PlayerTurn = board.Players[((board.PlayerTurn.Id + i) % len(board.Players))]
			board.PlayerTurn.StartTimer()
			return
		}
	}
	return
}

func (board *Board) GetRankByPlayer(player *Player) int {
	ranking := make(map[*Player]int)
	for index, _ := range board.Players {
		ranking[board.Players[index]] = board.Players[index].Score
	}
	type kv struct {
		Key   *Player
		Value int
	}
	var ss []kv
	for k, v := range ranking {
		ss = append(ss, kv{k, v})
	}
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	for index, kv := range ss {
		if kv.Key == player {
			return index + 1
		}
	}
	return 777
}
