package randomsuff

import (
	"math/rand"
	"time"
)

func RandSeq(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n+1)
	b[0] = rune('/')
	for i := range b[1:] {
		b[i+1] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
