package randomsuff

import (
	"fmt"
	"math/rand"
	"time"
)

func RandSeq(n int) (string, error) {
	if n < 1 {
		err := fmt.Errorf("wromg argument: number %v less than 1\n ", n)
		return "", err
	}
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b), nil
}
/* 
b := make([]rune, n+1)
	b[0] = rune('/')
	for i := range b[1:] {
		b[i+1] = letters[rand.Intn(len(letters))]
	}
	return string(b), nil
}
 */