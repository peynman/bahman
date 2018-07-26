package main

import (
	"github.com/peyman-abdi/bahman/app/modules/kernel"
	"github.com/peyman-abdi/bahman/app/cli/console"
)

func main() {
	services := kernel.SetupCLIServices()

	consoleObj := console.New("console", services)
	if err := consoleObj.Run(); err != nil {
		panic(err)
	}
}
