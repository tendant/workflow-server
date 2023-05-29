package main

import "github.com/tendant/workflow-server/app"

func main() {
	newApp := app.Default()

	newApp.Run()
}
