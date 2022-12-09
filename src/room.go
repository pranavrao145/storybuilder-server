package main

type Room struct {
	/// id of this room
	id string

	/// a map that maps client IDs to their clients
	clients map[int]*Client

	/// the hub that this room is a part of
	hub *Hub

	/// messages that have been sent to this room
	messageReceiveQueue chan *Message
}

func getOrCreateRoom(hub *Hub, id string) *Room {
	if room, ok := hub.rooms[id]; !ok {
		return room
	}

	return &Room{
		id:                  id,
		clients:             map[int]*Client{},
		hub:                 hub,
		messageReceiveQueue: make(chan *Message),
	}
}
