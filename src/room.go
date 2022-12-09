package main

type Room struct {
	/// id of this room
	id string

	/// a list of clients this room contains
	clients []*Client

	/// the hub that this room is a part of
	hub *Hub

	/// messages that have been sent to this room
	messageReceiveQueue chan *Message
}
