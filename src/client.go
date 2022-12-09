package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

type Client struct {
	/// id of this client
	id int

	/// username of this client
	username string

	/// room that this client is connected to
	room *Room

	/// messages that this client must receive
	messageReceiveQueue chan *Message

	/// connection to the client
	conn *websocket.Conn
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	clientId, ok := r.URL.Query()["clientId"]

	if !ok {
		log.Println("Failed to get client ID from websocket connection request. Abandoning client...")
		return
	}

	var finalClientId int

	if num, err := strconv.Atoi(clientId[0]); err == nil {
		finalClientId = num
	} else {
		log.Println("Failed to parse client ID. Abandoning client...")
		return
	}

	username, ok := r.URL.Query()["username"]

	if !ok {
		log.Println("Failed to get client username from websocket connection request. Abandoning client...")
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{
		id:                  finalClientId,
		username:            username[0],
		room:                &Room{},
		messageReceiveQueue: make(chan *Message),
		conn:                conn,
	}
	hub.register <- client

	// start the client's writing
	// start the client's reading
}
