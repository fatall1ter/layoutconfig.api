package model

//easyjson:json
type AuthenticationRequestJSON struct {
	// Private token for authentication
	Token string `json:"token"`
}
