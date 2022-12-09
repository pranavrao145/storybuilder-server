package main

type Hub struct {
	/// map that maps room IDs to rooms
	rooms map[string]*Room

	/// channel for taking in any new clients that come in
	register chan *Client
}

func NewHub() *Hub {
	return &Hub{
		rooms:    map[string]*Room{},
		register: make(chan *Client),
	}
}
