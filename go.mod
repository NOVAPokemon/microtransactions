module github.com/NOVAPokemon/microtransactions

go 1.13

require (
	github.com/NOVAPokemon/utils v0.0.62
	github.com/bruno-anjos/archimedesHTTPClient v0.0.2
	github.com/golang/geo v0.0.0-20200730024412-e86565bf3f35
	github.com/gorilla/mux v1.8.0
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.6.0
	github.com/stretchr/testify v1.4.0
)

replace (
	github.com/NOVAPokemon/utils v0.0.62 => ../utils
	github.com/bruno-anjos/archimedesHTTPClient v0.0.2 => ./../../bruno-anjos/archimedesHTTPClient
	github.com/bruno-anjos/cloud-edge-deployment v0.0.1 => ../../bruno-anjos/cloud-edge-deployment
)
