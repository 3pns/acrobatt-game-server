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

func (ai *AI) Start () {
  request := Request{}
  client := ai.client
  GameRequestChannel := client.CurrentGame.RequestChannel
  player := client.CurrentGame.Board().Players[request.Client.GameId()]
  for {
    request = <- client.Ai.RequestChannel
    board := client.CurrentGame.Board()
    isPlayerTurn := player == board.PlayerTurn
    if request.Type == "Refresh" && isPlayerTurn{
      var req  = Request {"PlaceRandom", "", nil, "", nil}
      GameRequestChannel <- req
    }
  }
}