package domain

type User struct {
	Login    string `json:"login" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=6"`
}
