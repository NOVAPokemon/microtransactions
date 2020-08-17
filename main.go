package main

import (
	"github.com/NOVAPokemon/utils"
	transactionDB "github.com/NOVAPokemon/utils/database/transactions"
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
		locationTag := utils.GetLocationTag(utils.DefaultLocationTagsFilename, serverName)
		commsManager = utils.CreateDefaultDelayedManager(locationTag, false)
	}

	transactionDB.InitTransactionsDBClient(*flags.ArchimedesEnabled)
	utils.StartServer(serviceName, host, port, routes, commsManager)
}
