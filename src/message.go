package main

type Message struct {
	/// type of the message (mandatory)
	messageType string

	/// the room ID this message was sent from (mandatory)
	roomId string

	/// content of this message (optional)
	content string

	/// the username of the sender of this message (optional)
	senderUsername string

	/// the id of the sender of this message (optional)
	senderId int

	/// the username of the inteded recipient of this message (optional)
	recipientUsername string

	/// the id of the inteded recipient of this message (optional)
	recipientId int
}
