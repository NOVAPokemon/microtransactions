package main

import (
	"fmt"
	"github.com/NOVAPokemon/utils"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"time"
)

const host = utils.ServeHost
const port = utils.MicrotransactionsPort

var addr = fmt.Sprintf("%s:%d", host, port)

func main() {
	rand.Seed(time.Now().Unix())
	r := utils.NewRouter(routes)
	log.Infof("Starting MICROTRANSACTIONS server in port %d...\n", port)
	log.Fatal(http.ListenAndServe(addr, r))
}
