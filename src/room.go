package main

type Room struct {
	// id of this room
	id string

	// a map that maps client IDs to their clients
	clients map[int]*Client

	// the hub that this room is a part of
	hub *Hub

	// messages to broadcast to all clients
	broadcast chan *Message

	// channel for removing dead clients
	unregister chan *Client
}

func getRoom(hub *Hub, id string) (*Room, bool) {
	if room, ok := hub.rooms[id]; ok {
		return room, true
	}

	return nil, false
}

func newRoom(hub *Hub, id string) *Room {
	return &Room{
		id:        id,
		clients:   map[int]*Client{},
		hub:       hub,
		broadcast: make(chan *Message),
	}
}

func (r *Room) run() {
	for {
		select {
		case client := <-r.unregister:
			delete(r.clients, client.id)

			leaveMessage := &Message{
				MessageType:       "leave",
				RoomId:            client.room.id,
				Content:           "",
				SenderUsername:    "",
				SenderId:          -1,
				RecipientUsername: "",
				RecipientId:       -1,
			}

			for _, client := range client.room.clients {
				client.messageSendQueue <- leaveMessage
			}
		case message := <-r.broadcast:
			for _, client := range r.clients {
				client.messageSendQueue <- message
			}
		}
	}
}
