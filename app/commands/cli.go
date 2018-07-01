package main

import (
	"github.com/peyman-abdi/avalanche/app/modules/kernel"
	"github.com/peyman-abdi/avalanche/app/commands/console"
)

func main() {
	services := kernel.SetupCLIKernel()

	console := new(console.ConsoleAppImpl)
	console.Services = services

	console.SetupWithConfPath("console")

	if err := console.Run(); err != nil {
		panic(err)
	}
}
