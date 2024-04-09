// models/response.go

package models

type Response struct {
	Ok   bool `json:"ok"`
	Data Data `json:"data"`
}

type Data struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
