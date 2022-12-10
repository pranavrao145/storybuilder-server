package main

type Message struct {
	// type of the message (mandatory)
	MessageType string `json:"messageType"`

	// the room ID this message was sent from (mandatory)
	RoomId string `json:"roomId"`

	// Content of this message (optional)
	Content string `json:"content"`

	// the username of the sender of this message (optional)
	SenderUsername string `json:"senderUsername"`

	// the id of the sender of this message (optional)
	SenderId int `json:"senderId"`

	// the username of the inteded recipient of this message (optional)
	RecipientUsername string `json:"recipientUsername"`

	// the id of the inteded recipient of this message (optional)
	RecipientId int `json:"recipientId"`
}
