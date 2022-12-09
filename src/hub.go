package main

import "log"

type Hub struct {
	/// map that maps room IDs to rooms
	rooms map[string]*Room

	/// channel for taking in any new clients that come in
	register chan *Client

	/// channel for removing dead clients
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		rooms:    map[string]*Room{},
		register: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			if client.isHost {
				client.room = newRoom(h, client.room.id)
				client.room.clients[client.id] = client
			} else {
				room, err := getRoom(h, client.room.id)

				if err != nil {
					log.Panic("Failed to get a room from the room ID given.")
					return
				}

				client.room = room
				client.room.clients[client.id] = client
			}
		}
	}
}
