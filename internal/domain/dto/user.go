package dto

type User struct {
	ID       uint   `json:"id" example:"1"`
	Username string `json:"username" example:"Joe the technician"`
	Role     string `json:"role" example:"technician"`
	Password string `json:"password" example:"secure_password"`
}
