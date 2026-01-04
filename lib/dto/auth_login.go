package dto

type LoginRequest struct {
    UID    string `json:"uid" binding:"required"`
    TokenHash string `json:"token" binding:"required,min=4,max=64"`
}

type LoginResponse struct {
    AccessToken  map[string]string `json:"access_token"`
    RefreshToken map[string]string `json:"refresh_token"`
}

