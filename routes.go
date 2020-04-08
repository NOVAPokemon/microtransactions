package main

import (
	"github.com/NOVAPokemon/utils"
	"github.com/NOVAPokemon/utils/api"
)

const GET = "GET"
const POST = "POST"

const MakeTransactionName = "MAKE_TRANSACTION"
const GetTransactionOffersName = "GET_TRANSACTION_OFFERs"
const GetPerformedTransactionsName = "GET_PERFORMED_TRANSACTIONS"

var routes = utils.Routes{

	utils.Route{
		Name:        GetTransactionOffersName,
		Method:      GET,
		Pattern:     api.GetTransactionOffersRoute,
		HandlerFunc: GetTransactionOffers,
	},

	utils.Route{
		Name:        GetPerformedTransactionsName,
		Method:      GET,
		Pattern:     api.GetPerformedTransactionsRoute,
		HandlerFunc: GetPerformedTransactions,
	},

	utils.Route{
		Name:        MakeTransactionName,
		Method:      POST,
		Pattern:     api.MakeTransactionRoute,
		HandlerFunc: MakeTransaction,
	},
}
