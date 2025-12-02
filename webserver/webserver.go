package webserver

import (
	"fmt"
	"github.com/jibuene/true-rand/random-number"
	"net/http"
)

var realRandom = randomnumber.New(randomnumber.TwitchMessages)
var notRandom = randomnumber.New(randomnumber.NotRandom)

func getTrueRandom(w http.ResponseWriter, req *http.Request) {
	realRandom.Generate()

	response := fmt.Sprintf("Random Number: %s\nBasis: %s\n", realRandom.Number.String(), realRandom.Basis)
	w.Write([]byte(response))
}

func getNotRandom(w http.ResponseWriter, req *http.Request) {
	notRandom.Generate()

	response := fmt.Sprintf("Random Number: %s\nBasis: %s\n", notRandom.Number.String(), notRandom.Basis)
	w.Write([]byte(response))
}

func StartServer() {
	http.HandleFunc("/random", getTrueRandom)
	http.HandleFunc("/notrandom", getNotRandom)
	http.HandleFunc("/", mainFrontEnd)
	http.HandleFunc("/randomNumber", randomNumberFrontend)

	port := ":8090"
	fmt.Println("Web server is running on http://localhost" + port)
	http.ListenAndServe(port, nil)
}
