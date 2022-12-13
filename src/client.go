package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	// id of this client
	id int

	// username of this client
	username string

	// room that this client is connected to
	room *Room

	// messages that are to be sent to the peer
	messageSendQueue chan *Message

	// if the client is a host or not
	isHost bool

	// connection to the client
	conn *websocket.Conn
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

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

	roomId, ok := r.URL.Query()["roomId"]

	if !ok {
		log.Println("Failed to get client room ID from websocket connection request. Abandoning client...")
		return
	}

	isHost, ok := r.URL.Query()["isHost"]

	if !ok {
		log.Println("Failed to get client host status from websocket connection request. Abandoning client...")
		return
	}

	var finalIsHost bool

	if isHost, err := strconv.ParseBool(isHost[0]); err == nil {
		finalIsHost = isHost
	} else {
		log.Println("Failed to parse client ID. Abandoning client...")
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{
		id:       finalClientId,
		username: username[0],
		room: &Room{
			id: roomId[0],
		},
		messageSendQueue: make(chan *Message),
		isHost:           finalIsHost,
		conn:             conn,
	}

	hub.register <- client

	go client.read()
	go client.write()
}

// / read: this method handles reading from the websocket connection and broadcasting messages that this client receives to all other clients
func (c *Client) read() {
	defer func() {
		c.conn.Close()
		c.room.unregister <- c
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, rawMessage, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		rawMessage = bytes.TrimSpace(bytes.Replace(rawMessage, newline, space, -1))

		message := &Message{}

		if err := json.Unmarshal(rawMessage, message); err != nil {
			log.Printf("error: %v", err)
			break
		}

		nextClientId, ok := getNextClientId(c.room, c.id)

		if !ok {
			log.Println("Failed to get next client ID in room.")
			return
		}

		if message.MessageType == "story" {
			message.RecipientId = nextClientId
			message.RecipientUsername = c.room.clients[nextClientId].username

			// add the user's name to the story line before storing it
			finalStoryLine := message.SenderUsername + ": " + message.Content

			err := addStoryLineToDb(c.room.id, finalStoryLine)

			if err != nil {
				log.Println(err)
				return
			}
		}

		c.room.broadcast <- message
	}
}

func (c *Client) write() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case rawMessage, ok := <-c.messageSendQueue:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))

			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)

			if err != nil {
				return
			}

			message, err := json.Marshal(rawMessage)

			if err != nil {
				log.Printf("error: %v", err)
				break
			}

			w.Write(message)

			// Add queued messages to the current websocket message.
			n := len(c.messageSendQueue)
			for i := 0; i < n; i++ {
				w.Write(newline)

				message, err := json.Marshal(<-c.messageSendQueue)

				if err != nil {
					log.Printf("error: %v", err)
					break
				}

				w.Write(message)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
