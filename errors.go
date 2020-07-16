package main

import (
	"fmt"

	"github.com/NOVAPokemon/utils"
	"github.com/pkg/errors"
)

const (
	errorOfferNotFoundFormat = "offer %s not found"
	errorLoadOffers          = "error loading offers"
)

// Handler wrappers
func wrapGetTransactionsError(err error) error {
	return errors.Wrap(err, fmt.Sprintf(utils.ErrorInHandlerFormat, getTransactionOffersName))
}

func wrapMakeTransactionError(err error) error {
	return errors.Wrap(err, fmt.Sprintf(utils.ErrorInHandlerFormat, makeTransactionName))
}

func wrapGetPerfomedTransactionsError(err error) error {
	return errors.Wrap(err, fmt.Sprintf(utils.ErrorInHandlerFormat, getPerformedTransactionsName))
}

// Other wrappers
func wrapLoadOffersError(err error) error {
	return errors.Wrap(err, errorLoadOffers)
}

// Error builders
func newOfferNotFoundError(offerId string) error {
	return errors.New(fmt.Sprintf(errorOfferNotFoundFormat, offerId))
}
