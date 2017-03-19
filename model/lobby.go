package model

type Lobby struct {
	Clients []*Client
	Master *Client
}

type LobbySlice struct {
	lobbies []*Lobby
}

func NewLobby(client *Client) *Lobby {
	var lobby Lobby
	lobby.Clients = []*Client{}
	lobby.Clients = append(lobby.Clients, client)
	lobby.Master = client
	client.CurrentLobby = &lobby
	return &lobby
}
