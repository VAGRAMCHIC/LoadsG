package dto

type LoginRequest struct {
    UID    string `json:"uid" binding:"required"`
    PasswordHash string `json:"password" binding:"required,min=4,max=64"`
}

type LoginResponse struct {
    AccessToken string `json:"access_token"`
}

