package main

import "log"

type Hub struct {
	/// map that maps room IDs to rooms
	rooms map[string]*Room

	/// channel for taking in any new clients that come in
	register chan *Client
}

func newHub() *Hub {
	return &Hub{
		rooms:    map[string]*Room{},
		register: make(chan *Client),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			room, ok := getRoom(h, client.room.id)

			if !ok {
				log.Panic("Failed to get a room from the room ID given.")
				return
			}

			client.room = room
			client.room.clients[client.id] = client

			// TODO: send a join signal to all clients
		}
	}
}
