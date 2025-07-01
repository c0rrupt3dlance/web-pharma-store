package models

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name" binding:"required"`
	Surname  string `json:"surname" binding:"required"`
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	IsAdmin  bool   `json:"is_admin"`
}
