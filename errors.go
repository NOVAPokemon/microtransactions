package main

import (
	"fmt"
	"github.com/NOVAPokemon/utils"
	"github.com/pkg/errors"
)

const (
	errorOfferNotFoundFormat = "offer %s not found"

	errorLoadOffers = "error loading offers"
)

// Handler wrappers
func wrapGetTransactionsError(err error) error {
	return errors.Wrap(err, fmt.Sprintf(utils.ErrorInHandlerFormat, GetTransactionOffersName))
}

func wrapMakeTransactionError(err error) error {
	return errors.Wrap(err, fmt.Sprintf(utils.ErrorInHandlerFormat, MakeTransactionName))
}

func wrapGetPerfomedTransactionsError(err error) error {
	return errors.Wrap(err, fmt.Sprintf(utils.ErrorInHandlerFormat, GetPerformedTransactionsName))
}

// Other wrappers
func wrapLoadOffersError(err error) error {
	return errors.Wrap(err, errorLoadOffers)
}

// Error builders
func newOfferNotFoundError(offerId string) error {
	return errors.New(fmt.Sprintf(errorOfferNotFoundFormat, offerId))
}
