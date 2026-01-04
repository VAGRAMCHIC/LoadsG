package dto

type CreateUserRequest struct {
    Comment    string `json:"comment"`
    Token string `json:"token" binding:"required,min=8,max=64"`
}

type CreateUserResponse struct {
    UID     string `json:"uid"`
		Comment string `json:"comment"`
}


