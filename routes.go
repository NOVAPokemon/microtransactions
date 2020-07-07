package main

import (
	"strings"

	"github.com/NOVAPokemon/utils"
	"github.com/NOVAPokemon/utils/api"
)

const (
	makeTransactionName          = "MAKE_TRANSACTION"
	getTransactionOffersName     = "GET_TRANSACTION_OFFERS"
	getPerformedTransactionsName = "GET_PERFORMED_TRANSACTIONS"
)

const (
	get  = "GET"
	post = "POST"
)

var routes = utils.Routes{
	api.GenStatusRoute(strings.ToLower(serviceName)),
	utils.Route{
		Name:        getTransactionOffersName,
		Method:      get,
		Pattern:     api.GetTransactionOffersRoute,
		HandlerFunc: getTransactionOffers,
	},

	utils.Route{
		Name:        getPerformedTransactionsName,
		Method:      get,
		Pattern:     api.GetPerformedTransactionsRoute,
		HandlerFunc: getPerformedTransactions,
	},

	utils.Route{
		Name:        makeTransactionName,
		Method:      post,
		Pattern:     api.MakeTransactionRoute,
		HandlerFunc: makeTransaction,
	},
}
