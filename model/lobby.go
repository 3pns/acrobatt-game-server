package model

type Lobby struct {
	Id int
	Name string
	Clients []*Client
	Master *Client
}

type LobbyFactory struct {
	Id       int
}

func NewLobbyFactory() *LobbyFactory {
	var factory = new(LobbyFactory)
	factory.Id = 0
	return factory
}

type LobbySlice struct {
	Lobbies []*Lobby `json:"lobbies"`
}

func (factory *LobbyFactory)NewLobby(client *Client) *Lobby {
	var lobby Lobby
	lobby.Id = factory.Id
	factory.Id++
	lobby.Name = "TEST"
	lobby.Clients = []*Client{}
	lobby.Clients = append(lobby.Clients, client)
	lobby.Master = client
	client.CurrentLobby = &lobby
	return &lobby
}
