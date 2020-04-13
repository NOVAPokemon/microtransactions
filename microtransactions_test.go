package main

import (
	"fmt"
	"github.com/NOVAPokemon/utils"
	"github.com/NOVAPokemon/utils/clients"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

var authClientTest = clients.NewAuthClient(fmt.Sprintf("%s:%d", utils.Host, utils.AuthenticationPort))
var trainersClientTest = clients.NewTrainersClient(fmt.Sprintf("%s:%d", utils.Host, utils.TrainersPort))
var transactionsClientTest = clients.NewMicrotransactionsClient(fmt.Sprintf("%s:%d", utils.Host, utils.MicrotransactionsPort))

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

	username := RandomString(10)
	password := RandomString(10)
	err := authClientTest.Register(username, password)

	if err != nil {
		logrus.Error(err)
		t.FailNow()
		return
	}

	err = authClientTest.LoginWithUsernameAndPassword(username, password)

	if err != nil {
		logrus.Error(err)
		t.FailNow()
		return
	}

	err = trainersClientTest.GetTrainerStatsToken(username, authClientTest.AuthToken)

	if err != nil {
		logrus.Error(err)
		t.FailNow()
		return
	}

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
		if transaction.Id.Hex() == id.Hex() {
			contains = true
			break
		}
	}

	assert.True(t, contains)

	if updatedTkn == trainersClient.TrainerStatsToken {
		t.Log("Token sent: " + trainersClient.TrainerStatsToken)
		t.Log("Token rcvd: " + updatedTkn)
		t.Error("Stats token did not update")
		t.Fail()
	}
}

func TestGetPerformedTransactions(t *testing.T) {

	username := RandomString(10)
	password := RandomString(10)
	err := authClientTest.Register(username, password)

	if err != nil {
		logrus.Error(err)
		t.FailNow()
		return
	}

	err = authClientTest.LoginWithUsernameAndPassword(username, password)

	if err != nil {
		logrus.Error(err)
		t.FailNow()
		return
	}

	performedTransactions, err := transactionsClientTest.GetTransactionRecords(authClientTest.AuthToken)

	if err != nil {
		logrus.Error(err)
		t.FailNow()
		return
	}
	assert.Empty(t, performedTransactions)
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
