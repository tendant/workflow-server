package main

import (
	"github.com/tendant/workflow-server/app"
	"github.com/tendant/workflow-server/handler"
)

func main() {
	app := app.Default()

	handle := handler.Handle{
		Slog: app.Slog,
	}

	handler.Routes(app.R, handle)

	app.Run()
}
