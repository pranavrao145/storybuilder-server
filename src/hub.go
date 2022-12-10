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

			joinMessage := &Message{
				messageType:       "join",
				roomId:            client.room.id,
				content:           "",
				senderUsername:    "",
				senderId:          -1,
				recipientUsername: "",
				recipientId:       -1,
			}

			for _, client := range room.clients {
				client.messageSendQueue <- joinMessage
			}
		}
	}
}
