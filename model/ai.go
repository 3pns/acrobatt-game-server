package model

type AI struct {
  RequestChannel chan Request
  Difficulty string
  client *Client
  Player *Player
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
  player := ai.Player
  for {
    request = <- ai.RequestChannel
    board := client.CurrentGame.Board()
    isPlayerTurn := player.Id == board.PlayerTurn.Id
    if request.Type == "Refresh" && isPlayerTurn{
      var req  = Request {"PlaceRandom", "", nil, "", client}
      client.CurrentGame.RequestChannel <- req
    }
  }
}