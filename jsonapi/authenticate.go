package jsonapi

type AuthenticateJson struct {
  PlayerId int `json:"player_id"`
  AccessToken string `json:"access_token"`
  Client string `json:"client"`
}