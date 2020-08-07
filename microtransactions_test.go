package main

import (
	"math/rand"
	http "github.com/bruno-anjos/archimedesHTTPClient"
	"os"
	"testing"
	"time"

	"github.com/NOVAPokemon/utils"
	"github.com/NOVAPokemon/utils/clients"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var (
	authClientTest         = clients.NewAuthClient(commsManager)
	trainersClientTest     = clients.NewTrainersClient(&http.Client{}, commsManager)
	transactionsClientTest = clients.NewMicrotransactionsClient(commsManager)
)

func TestMain(m *testing.M) {

	username := RandomString(10)
	password := RandomString(10)
	err := authClientTest.Register(username, password)

	if err != nil {
		logrus.Error(err)
		return
	}

	err = authClientTest.LoginWithUsernameAndPassword(username, password)

	if err != nil {
		logrus.Error(err)
		return
	}

	err = trainersClientTest.GetAllTrainerTokens(username, authClientTest.AuthToken)

	if err != nil {
		logrus.Error(err)
		return
	}

	res := m.Run()
	os.Exit(res)
}

// Location should be added
func TestGetOffers(t *testing.T) {

	transactions, err := transactionsClientTest.GetOffers()

	if err != nil {
		logrus.Error(err)
		t.Error(err)
		t.FailNow()
	}

	logrus.Info(transactions)
}

// Location should be added
func TestMakeTransaction(t *testing.T) {

	offers, err := transactionsClientTest.GetOffers()

	if err != nil {
		logrus.Error(err)
		t.Error(err)
		t.FailNow()
	}

	id, updatedTkn, err := transactionsClientTest.PerformTransaction(offers[len(offers)-1].Name, authClientTest.AuthToken, trainersClientTest.TrainerStatsToken)

	if err != nil {
		logrus.Error(err)
		t.Error(err)
		t.FailNow()
	}

	t.Logf("Made transaction: %s", id)

	performedTransactions, err := transactionsClientTest.GetTransactionRecords(authClientTest.AuthToken)

	if err != nil {
		logrus.Error(err)
		t.FailNow()
		return
	}
	assert.NotEmpty(t, performedTransactions)

	contains := false

	for _, transaction := range performedTransactions {
		if transaction.Id == *id {
			contains = true
			break
		}
	}

	assert.True(t, contains)

	if updatedTkn == trainersClientTest.TrainerStatsToken {
		t.Log("Token sent: " + trainersClientTest.TrainerStatsToken)
		t.Log("Token rcvd: " + updatedTkn)
		t.Error("Stats token did not update")
		t.Fail()
	}
}

func TestGetPerformedTransactions(t *testing.T) {

	performedTransactions, err := transactionsClientTest.GetTransactionRecords(authClientTest.AuthToken)

	if err != nil {
		logrus.Error(err)
		t.FailNow()
		return
	}
	assert.IsType(t, []utils.TransactionRecord{}, performedTransactions)
}

func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	rand.Seed(time.Now().Unix())

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}
