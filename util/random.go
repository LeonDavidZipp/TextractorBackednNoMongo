package util

import (
    "math/rand"
    "time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func RandomName() string {
	return RandomString(8)
}

func RandomEmail() string {
	return RandomString(10) + "@example.com"
}

func RandomURL() string {
	return RandomString(10) + ".com"
}

func RandomText() string {
	return RandomString(100)
}
