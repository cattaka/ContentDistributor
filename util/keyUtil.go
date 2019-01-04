package util

import (
	"math/rand"
	"time"
)

var randSrc = rand.NewSource(time.Now().UnixNano())

const (
	letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func RandString(length int) string {
	b := make([]byte, length)
	var r = randSrc.Int63()
	for i := 0; i < length; i++ {
		if r == 0 {
			r = randSrc.Int63()
			if r < 0 {
				r = -r
			}
		}
		b[i] = letters[int(r%int64(len(letters)))]
		r /= int64(len(letters))
	}
	return string(b)
}
