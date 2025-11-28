package webserver

import (
	"fmt"
	"github.com/jibuene/true-rand/twitch"
	"math/big"
	"net/http"
	"strconv"
)

func hello(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "hello\n")
}

func headers(w http.ResponseWriter, req *http.Request) {

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func getTrueRandom(w http.ResponseWriter, req *http.Request) {

	twitchmsgs := twitchmsg.DoTwitchRequest()
	num := numbeFromString(twitchmsgs[:]...)

	fmt.Fprintf(w, "Twitch Messages Used:\n")
	for _, msg := range twitchmsgs {
		fmt.Fprintf(w, "Twitch Message: %s\n", msg)
	}

	fmt.Fprintf(w, "True Random Number: %s\n", num.String())
}

func isInterger(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func numbeFromString(strings ...string) *big.Int {
	var result string
	for _, s := range strings {
		for _, r := range s {
			intStr := fmt.Sprintf("%X", r)
			if isInterger(intStr) {
				result += intStr
			}
		}
	}

	bigInt, success := new(big.Int).SetString(result, 16)
	if !success {
		panic("Failed to convert string to big.Int")
	}

	return bigInt
}

func StartServer() {

	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)
	http.HandleFunc("/random", getTrueRandom)

	fmt.Println("Web server is running on http://localhost:8090")
	http.ListenAndServe(":8090", nil)
}
