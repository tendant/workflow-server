package main

import (
	"embed"
	"os"

	"github.com/tendant/workflow-server/app"
	"github.com/tendant/workflow-server/dsl"
	"github.com/tendant/workflow-server/handler"
	"go.temporal.io/sdk/client"
)

//go:embed static
var ef embed.FS

func main() {
	app := app.Default()
	slog := app.Slog

	// Create the client object just once per process
	c, err := client.Dial(client.Options{})
	if err != nil {
		slog.Error("unable to create Temporal client", "err", err)
		os.Exit(1)
	}
	defer c.Close()

	handle := handler.Handle{
		Slog:   slog,
		Client: c,
		Ef:     ef,
		TxnMap: make(map[int]dsl.WorkflowRunAct),
	}

	handler.Routes(app.R, handle)

	app.Run()
}
