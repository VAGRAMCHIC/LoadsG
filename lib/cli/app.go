package cli

import "loadsg/lib/client"

type App struct {
	Client *client.Client
}

func NewApp(c *client.Client) *App {
	return &App{Client: c}
}