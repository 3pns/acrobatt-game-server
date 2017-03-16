package model

type AI struct {
  RequestChannel chan Request
  Difficulty string
  client *Client
}

func NewIA(client *Client) AI{
  var ai AI
  ai.RequestChannel = make(chan Request, 100)
  ai.Difficulty = "easy"
  ai.client = client
  return ai
}

func (client *Client) start () {
  request := Request{}
  board := client.CurrentGame.Board()
  GameRequestChannel := client.CurrentGame.RequestChannel
  player := board.Players[request.Client.GameId()]
  for {
    request = <- client.Ai.RequestChannel
    isPlayerTurn := player == board.PlayerTurn
    if request.Type == "Refresh" && isPlayerTurn{
      var req  = Request {"PlaceRandom", "", nil, "", nil}
      GameRequestChannel <- req
    }
  }
}