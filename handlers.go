package main

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/NOVAPokemon/utils"
	"github.com/NOVAPokemon/utils/api"
	"github.com/NOVAPokemon/utils/clients"
	transactionDB "github.com/NOVAPokemon/utils/database/transactions"
	"github.com/NOVAPokemon/utils/tokens"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

const offersFile = "microtransaction_offers.json"

var (
	// variables
	offersMap       map[string]utils.TransactionTemplate
	marshaledOffers []byte
)

var httpClient = &http.Client{}

func init() {
	var err error
	offersMap, marshaledOffers, err = loadOffers()
	if err != nil {
		log.Fatal(err)
	}
}

func getTransactionOffers(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write(marshaledOffers)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapGetTransactionsError(err), http.StatusInternalServerError)
	}
}

func makeTransaction(w http.ResponseWriter, r *http.Request) {
	offerId := mux.Vars(r)[api.OfferIdPathVar]
	log.Infof("Got transaction request for offer: %s", offerId)

	offer, ok := offersMap[offerId]
	if !ok {
		err := wrapMakeTransactionError(newOfferNotFoundError(offerId))
		utils.LogAndSendHTTPError(&w, err, http.StatusNotFound)
		return
	}

	authToken, err := tokens.ExtractAndVerifyAuthToken(r.Header)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapMakeTransactionError(err), http.StatusUnauthorized)
		return
	}

	trainerStatsToken, err := tokens.ExtractAndVerifyTrainerStatsToken(r.Header)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapMakeTransactionError(err), http.StatusUnauthorized)
		return
	}

	trainersClient := clients.NewTrainersClient(httpClient)
	valid, err := trainersClient.VerifyTrainerStats(authToken.Username, trainerStatsToken.TrainerHash, r.Header.Get(tokens.AuthTokenHeaderName))
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapMakeTransactionError(err), http.StatusUnauthorized)
		return
	}

	if !*valid {
		err = wrapMakeTransactionError(tokens.ErrorInvalidStatsToken)
		utils.LogAndSendHTTPError(&w, err, http.StatusUnauthorized)
		return
	}

	makeTransactionWithBankEntity(offer)

	transactionRecord := utils.TransactionRecord{
		TemplateName: offer.Name,
		User:         authToken.Username,
	}

	id, err := transactionDB.AddTransaction(transactionRecord)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapMakeTransactionError(err), http.StatusInternalServerError)
		return
	}

	newTrainerStats := utils.TrainerStats{
		Level: trainerStatsToken.TrainerStats.Level,
		Coins: trainerStatsToken.TrainerStats.Coins + offer.Coins,
	}

	newStats, err := trainersClient.UpdateTrainerStats(authToken.Username, newTrainerStats,
		r.Header.Get(tokens.AuthTokenHeaderName))
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapMakeTransactionError(err), http.StatusInternalServerError)
		return
	}

	log.Infof("Previous Coins: %d", trainerStatsToken.TrainerStats.Coins)
	log.Infof("Updated Coins: %d", newStats.Coins)

	w.Header().Set(tokens.StatsTokenHeaderName, trainersClient.TrainerStatsToken)
	toSend, err := id.MarshalJSON()
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapMakeTransactionError(err), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(toSend)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapMakeTransactionError(err), http.StatusInternalServerError)
	}
}

func getPerformedTransactions(w http.ResponseWriter, r *http.Request) {
	authToken, err := tokens.ExtractAndVerifyAuthToken(r.Header)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapGetPerfomedTransactionsError(err), http.StatusUnauthorized)
		return
	}

	performedTransactions, err := transactionDB.GetTransactionsFromUser(authToken.Username)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapGetPerfomedTransactionsError(err), http.StatusInternalServerError)
		return
	}

	toSend, err := json.Marshal(performedTransactions)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapGetPerfomedTransactionsError(err), http.StatusUnauthorized)
		return
	}

	_, err = w.Write(toSend)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapGetPerfomedTransactionsError(err), http.StatusUnauthorized)
	}
}

func loadOffers() (map[string]utils.TransactionTemplate, []byte, error) {
	data, err := ioutil.ReadFile(offersFile)
	if err != nil {
		return nil, nil, wrapLoadOffersError(err)
	}

	var offersArr []utils.TransactionTemplate
	err = json.Unmarshal(data, &offersArr)
	if err != nil {
		return nil, nil, wrapLoadOffersError(err)
	}

	var offersMapAux = make(map[string]utils.TransactionTemplate, len(offersArr))
	for _, offer := range offersArr {
		offersMapAux[offer.Name] = offer
	}

	log.Infof("Loaded %d offers.", len(offersArr))

	marshaledOffers, err = json.Marshal(offersArr)
	if err != nil {
		return nil, nil, wrapLoadOffersError(err)
	}

	return offersMapAux, marshaledOffers, nil
}

func makeTransactionWithBankEntity(offer utils.TransactionTemplate) {
	log.Infof("Making transaction %s", offer.Name)

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(1500)
	time.Sleep(time.Duration(n) * time.Millisecond)

	return
}
