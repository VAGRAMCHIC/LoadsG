package handler


import "loadsg/lib/service"

type Handler struct{
	User *UserHandler
	Auth *AuthHandler
}

func NewHandler(
	userService service.UserService,
	authService service.AuthService,
) *Handler {
	return &Handler{
		User: NewUserHandler(userService),
		Auth: NewAuthHandler(authService),
	}
}

