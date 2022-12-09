package main

type Client struct {
	/// id of this client
	id int

	/// username of this client
	username string

	/// room that this client is connected to
	room *Room

	/// messages that this client must receive
	messageReceiveQueue chan *Message
}
