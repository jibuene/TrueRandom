package randomnumber

import (
	"fmt"
	"math/big"
	"strconv"

	twitchmsg "github.com/jibuene/true-rand/twitch"
)

const NotRandomNumber = 42

type numberWithBasis struct {
	Number *big.Int
	Basis  string
}

type RandomNumberGenerator struct {
	Number          *big.Int                // The generated random number
	Basis           string                  // The basis or source of randomness
	numberGenerator func() *numberWithBasis // Function to generate the random number
}

type RandomnessType int

const (
	TwitchMessages RandomnessType = iota
	NotRandom
)

// New creates a new RandomNumberGenerator based on the specified RandomnessType.
func New(randomnessType RandomnessType) *RandomNumberGenerator {
	var rng RandomNumberGenerator

	switch randomnessType {
	case TwitchMessages:
		rng.numberGenerator = getTwitchRandom
	case NotRandom:
		rng.numberGenerator = func() *numberWithBasis {
			return &numberWithBasis{
				Number: big.NewInt(NotRandomNumber),
				Basis:  "Not Random",
			}
		}
	default:
		panic("Unsupported randomness type")
	}

	return &rng
}

func (rng *RandomNumberGenerator) Generate() {
	result := rng.numberGenerator()
	rng.Number = result.Number
	rng.Basis = result.Basis
}

// getTwitchRandom fetches Twitch messages and generates a random number from them.
func getTwitchRandom() *numberWithBasis {
	twitchmsgs := twitchmsg.FetchTwitchMessages()
	num := numberFromString(twitchmsgs[:]...)

	var basis string
	for _, msg := range twitchmsgs {
		basis += msg
	}

	return &numberWithBasis{
		Number: num,
		Basis:  basis,
	}
}

// numberFromString converts a list of strings into a big.Int by extracting integer characters.
func numberFromString(strings ...string) *big.Int {
	var result string
	for _, s := range strings {
		for _, r := range s {
			intStr := fmt.Sprintf("%X", r)
			if isInteger(intStr) {
				result += intStr
			}
		}
	}

	bigInt, success := new(big.Int).SetString(result, 16)
	if !success {
		panic("Failed to convert string to big.Int")
	}

	bigInt = bigInt.Mod(bigInt, big.NewInt(1_337_000))

	return bigInt
}

// isInteger checks if a string can be converted to an integer.
func isInteger(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}
