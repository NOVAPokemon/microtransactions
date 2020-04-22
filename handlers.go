package main

import (
	"encoding/json"
	"errors"
	"github.com/NOVAPokemon/utils"
	"github.com/NOVAPokemon/utils/api"
	"github.com/NOVAPokemon/utils/clients"
	transactionDB "github.com/NOVAPokemon/utils/database/transactions"
	"github.com/NOVAPokemon/utils/tokens"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

const OffersFile = "microtransaction_offers.json"

var (
	// errors
	OfferNotFound = errors.New("offer not found")

	// variables
	offersMap       map[string]utils.TransactionTemplate
	marshaledOffers []byte
)

var httpClient = &http.Client{}

func init() {
	offersMap, marshaledOffers = loadOffers()
}

func GetTransactionOffers(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write(marshaledOffers)
}

func MakeTransaction(w http.ResponseWriter, r *http.Request) {

	log.Infof("Got transaction request for offer: %s", mux.Vars(r)[api.OfferIdPathVar])
	offer, ok := offersMap[mux.Vars(r)[api.OfferIdPathVar]]

	if !ok {
		log.Infof("offer %s not found", mux.Vars(r)[api.OfferIdPathVar])
		http.Error(w, OfferNotFound.Error(), http.StatusNotFound)
		return
	}

	authToken, err := tokens.ExtractAndVerifyAuthToken(r.Header)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	trainerStatsToken, err := tokens.ExtractAndVerifyTrainerStatsToken(r.Header)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	trainersClient := clients.NewTrainersClient(httpClient)
	valid, err := trainersClient.VerifyTrainerStats(authToken.Username, trainerStatsToken.TrainerHash, r.Header.Get(tokens.AuthTokenHeaderName))
	if err != nil || !*valid {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	err = makeTransactionWithBankEntity(offer)

	if err != nil {
		log.Error(err)
		http.Error(w, "Error occurred making transaction", http.StatusInternalServerError)
		return
	}

	transactionRecord := utils.TransactionRecord{
		TemplateName: offer.Name,
		User:         authToken.Username,
	}

	id, err := transactionDB.AddTransaction(transactionRecord)

	if err != nil {
		log.Error(err)
		http.Error(w, "Error occurred storing transaction", http.StatusInternalServerError)
		return
	}

	newTrainerStats := utils.TrainerStats{
		Level: trainerStatsToken.TrainerStats.Level,
		Coins: trainerStatsToken.TrainerStats.Coins + offer.Coins,
	}

	newStats, err := trainersClient.UpdateTrainerStats(authToken.Username, newTrainerStats, r.Header.Get(tokens.AuthTokenHeaderName))

	if err != nil {
		http.Error(w, "Error fetching trainer stats", http.StatusInternalServerError)
		return
	}

	log.Infof("Previous Coins: %d", trainerStatsToken.TrainerStats.Coins)
	log.Infof("Updated Coins: %d", newStats.Coins)

	w.Header().Set(tokens.StatsTokenHeaderName, trainersClient.TrainerStatsToken)
	toSend, _ := id.MarshalJSON()
	_, _ = w.Write(toSend)
}

func GetPerformedTransactions(w http.ResponseWriter, r *http.Request) {

	authToken, err := tokens.ExtractAndVerifyAuthToken(r.Header)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	performedTransactions, err := transactionDB.GetTransactionsFromUser(authToken.Username)

	toSend, err := json.Marshal(performedTransactions)

	if err != nil {
		log.Errorf("error occurred decoding :%s", performedTransactions)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	_, _ = w.Write(toSend)

}

func loadOffers() (map[string]utils.TransactionTemplate, []byte) {

	data, err := ioutil.ReadFile(OffersFile)
	if err != nil {
		log.Errorf("Error loading offers file ")
		log.Fatal(err)
		panic(err)
	}

	var offersArr []utils.TransactionTemplate
	err = json.Unmarshal(data, &offersArr)

	var offersMap = make(map[string]utils.TransactionTemplate, len(offersArr))
	for _, offer := range offersArr {
		offersMap[offer.Name] = offer
	}

	if err != nil {
		log.Errorf("Error unmarshalling offer names")
		log.Fatal(err)
		panic(err)
	}

	log.Infof("Loaded %d offers.", len(offersArr))

	marshaledOffers, _ = json.Marshal(offersArr)

	return offersMap, marshaledOffers
}

func makeTransactionWithBankEntity(offer utils.TransactionTemplate) error {
	log.Infof("Making transaction %s", offer.Name)
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(1500)
	time.Sleep(time.Duration(n) * time.Millisecond)
	return nil
}
