package main

import (
	"math/rand"
	"time"
)

func randSeq(n int) string {
	rand.Seed(time.Now().UnixNano())

	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)

	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

func getRoomMaxClientId(room *Room) int {
	max := 0

	for _, client := range room.clients {
		if client.id > max {
			max = client.id
		}
	}

	return max
}
