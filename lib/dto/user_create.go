package dto

type CreateUserRequest struct {
    UID    string `json:"uid" binding:"required"`
    Password string `json:"password" binding:"required,min=8,max=64"`
}

type CreateUserResponse struct {
    UID    string `json:"id"`
}


