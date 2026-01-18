package request

import "github.com/google/uuid"

type Register struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"email,required"`
	Password string `json:"password" binding:"required"`
}

type Login struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ConfirmUserEmail struct {
	Username string `json:"username"`
	Code     string `json:"code"`
}

type GetUserByID struct {
	ID uuid.UUID `json:"id" binding:"required,number"`
}

type GetUserByUsername struct {
	Username string `json:"username" binding:"required"`
}

type UpdateUserInfo struct {
	NewEmail    string `json:"email" binding:"email"`
	NewPassword string `json:"password"`
}

type UpdateUserStatus struct {
	NewStatus string `json:"status"`
}
