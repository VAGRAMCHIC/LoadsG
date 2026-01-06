package handler


import "loadsg/lib/service"

type Handler struct{
	User *UserHandler
	Auth *AuthHandler
	Load *LoadManagerHandler
}

func NewHandler(
	userService service.UserService,
	authService service.AuthService,
	loadManagerService service.LoadManagerService,
) *Handler {
	return &Handler{
		User: NewUserHandler(userService),
		Auth: NewAuthHandler(authService),
		Load: NewLoadManagerHandler(loadManagerService),
	}
}

