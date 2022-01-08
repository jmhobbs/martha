package util

import (
	"math/rand"
	"time"
)

var rng *rand.Rand
func init() {
rng = rand.New(rand.NewSource(time.Now().UnixNano())) //#nosec
}

func MakeId() string {
	return MakeIdFull(8, 4)
}

func MakeIdFull(length, split int) string {
	runes := []rune{}
	for i := 0; i < length; i++ {
		n := rune(rng.Intn(62))
		if n < 10 {
			n = n + 48
		} else if n < 36 {
			n = n + 65 - 10
		} else {
			n = n + 97 - 36
		}
		if i % split == 0 && i != 0 {
			runes = append(runes, '-', n)
		} else {
			runes = append(runes, n)
		}
	}
	return string(runes)
}