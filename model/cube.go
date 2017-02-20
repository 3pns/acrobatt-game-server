package model

type Cube struct {
	X int `json:"X"`
	Y int `json:"Y"`
}

func (cube *Cube) Project(origin Square, rotation string, flipped bool) Cube {
	var xFactor = 1
	var yFactor = 1
	var xBoardValue = cube.X
	var yBoardValue = cube.Y
	if flipped {
		xFactor = -xFactor
	}
	//fmt.Println("Avant : x :" + strconv.Itoa(xBoardValue) + ", y: " + strconv.Itoa(yBoardValue))
	if rotation == "W" {
		xBoardValue = -cube.Y
		yBoardValue = -cube.X
	} else if rotation == "N" {
		xBoardValue = -cube.X
		yBoardValue = -cube.Y
	} else if rotation == "E" {
		xBoardValue = cube.Y
		yBoardValue = cube.X
	}
	//fmt.Println("Apres : x :" + strconv.Itoa(xBoardValue) + ", y: " + strconv.Itoa(yBoardValue))
	//board.Squares[piece.Origin.X+xFactor*xBoardValue][piece.Origin.Y+yFactor*yBoardValue].PlayerId = piece.PlayerId
	return Cube{origin.X+xFactor*xBoardValue, origin.Y+yFactor*yBoardValue}
}