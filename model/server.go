package model

type Server struct {
	CurrentGames   []*Game
	ClientsInHome  []*Client
	ClientsInLobby []*Client
}

func (server *Server) Start() {

}
