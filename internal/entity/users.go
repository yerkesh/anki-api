package entity

type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required"`
	Role     string `json:"role"`
}
