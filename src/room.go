package main

import "errors"

type Room struct {
	/// id of this room
	id string

	/// a map that maps client IDs to their clients
	clients map[int]*Client

	/// the hub that this room is a part of
	// hub *Hub

	/// messages that have been sent to this room
	messageReceiveQueue chan *Message
}

func getRoom(hub *Hub, id string) (*Room, error) {
	if room, ok := hub.rooms[id]; !ok {
		return room, nil
	}

	return nil, errors.New("Failed to get room with room ID supplied.")
}

func newRoom(hub *Hub, id string) *Room {
	return &Room{
		id:                  id,
		clients:             map[int]*Client{},
		messageReceiveQueue: make(chan *Message),
	}
}
