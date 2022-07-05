package util

import (
	"math/rand"
	"time"

	"github.com/croixxant/go-sandbox/entity"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomCurrency テスト用。ランダムな通貨を返す。
func RandomCurrency() string {
	currencies := []string{entity.USD, entity.EUR, entity.CAD}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
