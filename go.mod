module github.com/NOVAPokemon/microtransactions

go 1.13

require (
	github.com/NOVAPokemon/utils v0.0.62
	github.com/gorilla/mux v1.7.4
	github.com/pkg/errors v0.8.1
	github.com/sirupsen/logrus v1.5.0
	github.com/stretchr/testify v1.4.0
)

replace github.com/NOVAPokemon/utils v0.0.62 => ../utils
