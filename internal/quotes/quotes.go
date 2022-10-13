package quotes

import (
	"math/rand"
	"strings"
	"time"
)

var (
	db []string
	r  *rand.Rand
)

func init() {
	db = strings.Split(body, "\n")
	r = rand.New(rand.NewSource(time.Now().Unix()))
}

func GetRandom() string {
	return db[r.Int()%len(db)]
}
