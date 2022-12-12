package main

import (
	"math/rand"
	"sort"
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

func findPositionInSlice(slice []int, requiredElenment int) (int, bool) {
	for idx, elem := range slice {
		if elem == requiredElenment {
			return idx, true
		}
	}

	return -1, false
}

func getNextClientId(room *Room, clientId int) (int, bool) {
	// generate slice of client IDs in this room
	keys := make([]int, len(room.clients))

	i := 0

	for k := range room.clients {
		keys[i] = k
		i++
	}

	// sort the slice in ascending order
	sort.Ints(keys)

	// find where the client ID is
	clientIdPosition, found := findPositionInSlice(keys, clientId)

	if !found {
		return -1, false
	}

	// return the next largest client ID
	return keys[(clientIdPosition+1)%len(keys)], true
}
