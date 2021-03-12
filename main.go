package main

import (
	"github.com/NOVAPokemon/utils"
)

const (
	host        = utils.ServeHost
	port        = utils.MicrotransactionsPort
	serviceName = "MICROTRANSACTIONS"
)

func main() {
	flags := utils.ParseFlags(serverName)

	if !*flags.LogToStdout {
		utils.SetLogFile(serverName)
	}

	if !*flags.DelayedComms {
		commsManager = utils.CreateDefaultCommunicationManager()
	} else {
		commsManager = utils.CreateDefaultDelayedManager(false)
	}

	utils.StartServer(serviceName, host, port, routes, commsManager)
}
