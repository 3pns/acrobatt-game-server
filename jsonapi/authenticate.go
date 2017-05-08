package jsonapi

type AuthenticateJson struct {
  PlayerId int `json:"player_id"`
  AccessToken string `json:"access-token"`
  Client string `json:"client"`
}