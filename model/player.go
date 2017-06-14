package model

import (
	. "../utils"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

//Exemple strategy design pattern :
//https://github.com/yksz/go-design-patterns/blob/master/behavior/strategy.go

type Player struct {
	Id                int     `json:"id"`
	Name              string  `json:"name"`
	Color             string  `json:"color"`
	Pieces            []Piece `json:"pieces"`
	startingSquares   []*Square
	squares           []*Square
	hasPlaceabePieces bool
	ClientId          int           `json:"client_id"`
	Score             int           `json:"score"`
	Time              time.Duration `json:"time"`
	TurnStartTime     time.Time     `json:"-"`
}

func (player *Player) Init() {
	fmt.Println("playerId in playerInit :", player.Id)
	for index, _ := range player.Pieces {
		player.Pieces[index].PlayerId = &player.Id
	}
	fmt.Println("init player pieces")
}

func NewPlayer(id int, name string, color string, pieces []Piece, startingSquares []*Square) *Player {
	var player Player
	player.Id = id
	player.Name = name
	player.Color = color
	player.Pieces = pieces
	player.startingSquares = startingSquares
	player.squares = []*Square{}
	player.hasPlaceabePieces = true
	player.ClientId = -1
	player.StartTimer()
	return &player
}

func (player *Player) PlacePiece(piece Piece, board *Board, simulation bool) *Piece {
	if piece.Origin == nil {
		fmt.Println("Fatal Error piece to place has no Origin")
		return nil
	}
	if player.Pieces[piece.Id].Origin != nil {
		fmt.Println("Fatal Error piece has already been used")
		return nil
	}
	piece.PlayerId = &player.Id
	piece.Cubes = board.Pieces[piece.Id].Cubes
	var projectedCubes []Cube
	var placementAuthorized = false
	//fmt.Println("----- Plaçage d'une pièce -----")
	for _, cube := range piece.Cubes {
		var projectedCube = cube.Project(*piece.Origin, piece.Rotation, piece.Flipped) // on projete le cube dans l'espace = vrai position
		projectedCubes = append(projectedCubes, projectedCube)                         // on ajoute le cube à la liste des cube projeté => càd dire les vrais cases occupés par la pièces sur la board
		//fmt.Println(projectedCubes)
		//si le cube est en dehors de la board le placement est interdit
		if projectedCube.X < 0 || projectedCube.X > 19 || projectedCube.Y < 0 || projectedCube.Y > 19 {
			//fmt.Println("SIGSEV Placement Out of Board Exception")
			placementAuthorized = false
			return nil
		}
		//si le cube occupe un square occupé le placement est interdit
		if board.Squares[projectedCube.X][projectedCube.Y].PlayerId != nil {
			//fmt.Println("StackOverflow Board Exception le square est déjà occupé")
			placementAuthorized = false
			return nil
		}
		// si le cube en bas est dans la board et appartient au joueur le placement est interdit
		if projectedCube.Y+1 >= 0 && projectedCube.Y+1 < 20 {
			if board.Squares[projectedCube.X][projectedCube.Y+1].GetPlayerId() == player.Id {
				//fmt.Println("Placement Unauthorized Exceptio.cuz cube en bas appartient au joueur")
				placementAuthorized = false
				return nil
			}
		}
		// si le cube en haut est dans la board et appartient au joueur le placement est interdit
		if projectedCube.Y-1 >= 0 && projectedCube.Y-1 < 20 {
			if board.Squares[projectedCube.X][projectedCube.Y-1].GetPlayerId() == player.Id {
				//fmt.Println("Placement Unauthorized Exceptio.cuz cube en haut appartient au joueur")
				placementAuthorized = false
				return nil
			}
		}
		// si le cube à gauche est dans la board et appartient au joueur le placement est interdit
		if projectedCube.X-1 >= 0 && projectedCube.X-1 < 20 {
			if board.Squares[projectedCube.X-1][projectedCube.Y].GetPlayerId() == player.Id {
				//fmt.Println("Placement Unauthorized Exceptio.cuz cube à gauche appartient au joueur")
				placementAuthorized = false
				return nil
			}
		}
		// si le cube à droite est dans la board et appartient au joueur le placement est interdit
		if projectedCube.X+1 >= 0 && projectedCube.X+1 < 20 {
			if board.Squares[projectedCube.X+1][projectedCube.Y].GetPlayerId() == player.Id {
				//fmt.Println("Placement Unauthorized Exceptio.cuz cube à gauche appartient au joueur")
				placementAuthorized = false
				return nil
			}
		}
		// si le cube est le cube de départ du joueur le placement est autorisé
		if player.IsAStartingCube(projectedCube) {
			placementAuthorized = true
			//fmt.Println("Placement Authorized cuz Starting Cube  :", projectedCube)
		}
		// si le cube en bas à gauche est dans la board et appartient au joueur le placement est autorisé
		if projectedCube.X-1 >= 0 && projectedCube.X-1 < 20 && projectedCube.Y+1 >= 0 && projectedCube.Y+1 < 20 {
			if board.Squares[projectedCube.X-1][projectedCube.Y+1].GetPlayerId() == player.Id {
				//fmt.Println("Placement Authorized cuz cube en bas à gauche")
				placementAuthorized = true
			}
		}
		// si le cube en bas à droite est dans la board et appartient au joueur le placement est autorisé
		if projectedCube.X+1 >= 0 && projectedCube.X+1 < 20 && projectedCube.Y+1 >= 0 && projectedCube.Y+1 < 20 {
			if board.Squares[projectedCube.X+1][projectedCube.Y+1].GetPlayerId() == player.Id {
				//fmt.Println("Placement Authorized cuz cube en bas à droite")
				placementAuthorized = true
			}
		}
		// si le cube en haut à gauche est dans la board et appartient au joueur le placement est autorisé
		if projectedCube.X-1 >= 0 && projectedCube.X-1 < 20 && projectedCube.Y-1 >= 0 && projectedCube.Y-1 < 20 {
			if board.Squares[projectedCube.X-1][projectedCube.Y-1].GetPlayerId() == player.Id {
				//fmt.Println("Placement Authorized cuz cube en haut à gauche")
				placementAuthorized = true
			}
		}
		// si le cube en haut à droite est dans la board et appartient au joueur le placement est autorisé
		if projectedCube.X+1 >= 0 && projectedCube.X+1 < 20 && projectedCube.Y-1 >= 0 && projectedCube.Y-1 < 20 {
			if board.Squares[projectedCube.X+1][projectedCube.Y-1].GetPlayerId() == player.Id {
				//fmt.Println("Placement Authorized cuz cube en haut à droite")
				placementAuthorized = true
			}
		}
	}
	if !placementAuthorized {
		//fmt.Println("----- BADDIES Placement Unauthorized Exception -----")
		return nil
	} else {
		fmt.Println("----- Placement Authorized -----")
		if simulation {
			fmt.Println("returning true because simulation")
			return &piece
		}
		fmt.Println(piece)
		for _, cube := range projectedCubes {
			board.Squares[cube.X][cube.Y].PlayerId = piece.PlayerId
			player.squares = append(player.squares, board.Squares[cube.X][cube.Y])
		}
		player.Pieces[piece.Id].Origin = piece.Origin
		player.Pieces[piece.Id].Rotation = piece.Rotation
		player.Pieces[piece.Id].Flipped = piece.Flipped
		fmt.Println("returning true because piece placé")
		return &player.Pieces[piece.Id]
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

func (player *Player) PlaceRandomPieceWithIAEasy(board *Board, simulation bool) *Piece {
	//TODO attacher le rand à la game lors du refactoring
	rand.Seed(time.Now().UTC().UnixNano())
	//on récupère les pièces restantes à placer
	var remainingPieces = []*Piece{}
	for index, piece := range player.Pieces {
		if piece.Origin == nil {
			remainingPieces = append(remainingPieces, &player.Pieces[index])
		}
	}
	//on récupère un index de piece au hasard ssi il reste des pièces à placer
	var index int
	var targetSquares = []*Square{}
	if len(remainingPieces) > 0 {
		index = rand.Intn(len(remainingPieces))
	} else {
		//le joueur a placé toutes ses pièces !
		return nil
	}
	//si le joueur a encore toutes ses pièces le square cible est son square de départ
	if len(remainingPieces) == 21 {
		targetSquares = player.startingSquares
	} else {
		//sinon à partir des squares appartenant au joueur on récupère les squares ou l'IA peut poser une pièce
		//fmt.Println("generating authorize squares for player", player.Id)
		for _, square := range player.squares {
			//fmt.Println("checking first player Square:", square.X, ",", square.Y)
			targetSquares = append(targetSquares, square.getDiagonalAuthorizedSquares(board)...)
		}
	}
	//on essaye de placer toutes les pièces
tryagain:
	//fmt.Println("remainingPieces: ", len(remainingPieces))
	//fmt.Println("playerSquares: ", player.squares)
	//fmt.Println("targetSquares: ", targetSquares)
	index = rand.Intn(len(remainingPieces))
	piece := remainingPieces[index]
	if player.TryPlacePieceOnSquares(board, piece, targetSquares, simulation) {
		return piece
	} else if len(remainingPieces) > 1 {
		//on enlève la pièce du slience
		remainingPieces = append(remainingPieces[:index], remainingPieces[index+1:]...)
		goto tryagain
	}

	//TODO Remove targetSquares duplicates in slice !
	//le joueur ne peut placer aucune pièce !
	return nil
}

func (player *Player) PlaceRandomPieceWithIAMedium(board *Board, simulation bool) *Piece {
	rand.Seed(time.Now().UTC().UnixNano())
	//classification des pièces par taille en nombre de cubes
	remainingPieces := make(map[int][]*Piece)
	remainingPieces[1] = []*Piece{}
	remainingPieces[2] = []*Piece{}
	remainingPieces[3] = []*Piece{}
	remainingPieces[4] = []*Piece{}
	remainingPieces[5] = []*Piece{}
	for index, _ := range player.Pieces {
		if player.Pieces[index].Origin == nil {
			remainingPieces[len(player.Pieces[index].Cubes)] = append(remainingPieces[len(player.Pieces[index].Cubes)], &player.Pieces[index])
		}
	}
	var index int
	var targetSquares = []*Square{}
	currentSize := 6
	if sizeOfremainingPieces(remainingPieces) <= 0 {
		return nil
	}
	if sizeOfremainingPieces(remainingPieces) == 21 {
		targetSquares = player.startingSquares
	} else {
		for _, square := range player.squares {
			targetSquares = append(targetSquares, square.getDiagonalAuthorizedSquares(board)...)
		}
	}
tryagain:
	currentSize = currentSize - 1
	if currentSize == 0 {
		return nil
	} else if len(remainingPieces[currentSize]) == 0 {
		goto tryagain
	}
	index = rand.Intn(len(remainingPieces[currentSize]))
	piece := remainingPieces[currentSize][index]
	if player.TryPlacePieceOnSquares(board, piece, targetSquares, simulation) {
		return piece
	} else if sizeOfremainingPieces(remainingPieces) > 1 {
		remainingPieces[currentSize] = append(remainingPieces[currentSize][:index], remainingPieces[currentSize][index+1:]...)
		goto tryagain
	}
	return nil
}

func sizeOfremainingPieces(remainingPieces map[int][]*Piece) int {
	size := 0
	for key := range remainingPieces {
		size += len(remainingPieces[key])
	}
	return size
}

func (player *Player) TryPlacePieceOnSquares(board *Board, piece *Piece, squares []*Square, simulation bool) bool {
	for _, square := range squares {
		//essayer les 8 rotation/coté possible
		piece.Flipped = false
		piece.Rotation = "N"
		//essayeer tous les positionnements de la pièce avec cette rotation/coté sur le square
		if player.TryPlacePieceOnSquareWithOrientation(board, *piece, square, simulation) {
			//on renvoit true si la pièce a été placé
			return true
		}
		piece.Rotation = "E"
		if player.TryPlacePieceOnSquareWithOrientation(board, *piece, square, simulation) {
			return true
		}
		piece.Rotation = "S"
		if player.TryPlacePieceOnSquareWithOrientation(board, *piece, square, simulation) {
			return true
		}
		piece.Rotation = "W"
		if player.TryPlacePieceOnSquareWithOrientation(board, *piece, square, simulation) {
			return true
		}
		piece.Flipped = true
		if player.TryPlacePieceOnSquareWithOrientation(board, *piece, square, simulation) {
			return true
		}
		piece.Rotation = "N"
		if player.TryPlacePieceOnSquareWithOrientation(board, *piece, square, simulation) {
			return true
		}
		piece.Rotation = "E"
		if player.TryPlacePieceOnSquareWithOrientation(board, *piece, square, simulation) {
			return true
		}
		piece.Rotation = "S"
		if player.TryPlacePieceOnSquareWithOrientation(board, *piece, square, simulation) {
			return true
		}
		piece.Rotation = "W"
		if player.TryPlacePieceOnSquareWithOrientation(board, *piece, square, simulation) {
			return true
		}
	}
	//false si la pièce n'a pas été placé
	return false
}

//essaye tous les positionnements de la pièce avec cette rotation/coté sur le square
func (player *Player) TryPlacePieceOnSquareWithOrientation(board *Board, piece Piece, square *Square, simulation bool) bool {
	for _, cube := range piece.Cubes {
		if AllowedCoordinates(square.X-cube.X, square.Y-cube.Y) {
			piece.Origin = board.Squares[square.X-cube.X][square.Y-cube.Y]
			if player.PlacePiece(piece, board, simulation) != nil {
				return true
			} else {
				piece.Origin = nil
			}
		}

	}
	return false
}

func (player *Player) HasPlaceabePieces(board *Board) bool {
	if player.hasPlaceabePieces {
		res := player.PlaceRandomPieceWithIAEasy(board, true)
		if res != nil {
			return true
		} else {
			player.hasPlaceabePieces = false
			return false
		}
	} else {
		return false
	}
}

func (player *Player) Concede() {
	player.hasPlaceabePieces = false
}

func (player *Player) SetApiId(apiId int) {
	player.ClientId = apiId
}

func (player *Player) ApiId() int {
	return player.ClientId
}

func (player *Player) StartTimer() {
	player.TurnStartTime = time.Now()
}
func (player *Player) GetTurnTime() time.Duration {
	return time.Since(player.TurnStartTime)
}
func (player *Player) UpdateScore(move Move) {
	points := len(move.Piece.Cubes)
	if player.PlacedPieceCount() == 21 {
		points += 15
		if len(move.Piece.Cubes) == 1 {
			points += 20
		}
	}
	player.Score += points
}
func (player *Player) PlacedPieceCount() int {
	count := 0
	for index, _ := range player.Pieces {
		if player.Pieces[index].Origin != nil {
			count++
		}
	}
	return count
}
