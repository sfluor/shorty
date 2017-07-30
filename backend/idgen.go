package main

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randIdGen(size int) string {
	b := make([]byte, size)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}

	return string(b)
}
