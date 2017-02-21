package model

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
