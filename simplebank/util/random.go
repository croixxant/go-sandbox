package util

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomCurrency generates a random currency code
func RandomCurrency() string {
	currencies := []string{USD, EUR, CAD}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
