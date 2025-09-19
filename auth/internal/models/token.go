package models

type RefreshToken struct {
	Id      int    `json:"id"`
	UserId  int    `json:"user_id"`
	Token   string `json:"token"`
	Revoked bool   `json:"revoked"`
}
