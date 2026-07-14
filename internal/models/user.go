package models

type User struct {
	ID        string `json:"id" db:"id"`
	Username  string `json:"username" db:"username"`
	Email     string `json:"email" db:"email"`
	Password  string `json:"password" db:"password_hash"`
	CreatedAt string `json:"created_at" db:"created_at"`
	UpdatedAt string `json:"updated_at" db:"updated_at"`
}

type UserInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required, min=8"`
}

type AuthPayload struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}
