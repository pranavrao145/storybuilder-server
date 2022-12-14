package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type GenerateRoomPayload struct {
	RoomId string `json:"roomId"`
}

type ValidateRoomPayload struct {
	Exists bool `json:"exists"`
}

type GetMembersPayload struct {
	RoomMembers []string `json:"roomMembers"`
}

type GenerateClientIdPayload struct {
	ClientId int `json:"clientId"`
}

type GetStoryLinesPayload struct {
	StoryLines []string `json:"storyLines"`
}

func handleGenerateRoom(hub *Hub, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	generatedId := randSeq(8)
	room := newRoom(hub, generatedId)
	hub.rooms[room.id] = room

	go room.run()

	p := GenerateRoomPayload{RoomId: generatedId}

	json.NewEncoder(w).Encode(p)
}

func handleGetMembers(hub *Hub, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	roomId, ok := r.URL.Query()["roomId"]

	if !ok {
		log.Println("No room ID provided. Returning 400.")
		w.WriteHeader(400)
		return
	}

	room, ok := hub.rooms[roomId[0]]

	if !ok {
		log.Println("Unable to find room in room list. Returning 404.")
		w.WriteHeader(404)
		return
	}

	p := GetMembersPayload{
		RoomMembers: []string{},
	}

	for _, client := range room.clients {
		if client.isHost {
			p.RoomMembers = append(p.RoomMembers, client.username+" (Host)")
		} else {
			p.RoomMembers = append(p.RoomMembers, client.username)
		}
	}

	json.NewEncoder(w).Encode(p)
}

func handleValidateRoom(hub *Hub, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	roomId, ok := r.URL.Query()["roomId"]

	if !ok {
		log.Println("No room ID provided. Returning 400.")
		w.WriteHeader(400)
		return
	}

	_, ok = hub.rooms[roomId[0]]

	var p ValidateRoomPayload

	if ok {
		p = ValidateRoomPayload{
			Exists: true,
		}
	} else {
		p = ValidateRoomPayload{
			Exists: false,
		}

		log.Println("Unable to find room in room list.")
	}

	json.NewEncoder(w).Encode(p)
}

func handleGenerateClientId(hub *Hub, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	roomId, ok := r.URL.Query()["roomId"]

	if !ok {
		log.Println("No room ID provided. Returning 400.")
		w.WriteHeader(400)
		return
	}

	var p GenerateClientIdPayload

	room, ok := hub.rooms[roomId[0]]

	if !ok {
		log.Println("Room does not exist. Returning 404.")
		w.WriteHeader(404)
		return
	} else {
		p.ClientId = getRoomMaxClientId(room) + 1
	}

	json.NewEncoder(w).Encode(p)
}

func handleGetAllStoryLines(hub *Hub, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	roomId, ok := r.URL.Query()["roomId"]

	if !ok {
		log.Println("No room ID provided. Returning 400.")
		w.WriteHeader(400)
		return
	}

	var p GetStoryLinesPayload

	room, ok := hub.rooms[roomId[0]]

	if !ok {
		log.Println("Room does not exist. Returning 404.")
		w.WriteHeader(404)
		return
	} else {
		storyLines, err := getStoryLinesFromDb(room.id)

		if err != nil {
			log.Println("Unable to get story lines from database. Returning 500.")
			w.WriteHeader(500)
			return
		}

		p.StoryLines = storyLines
	}

	json.NewEncoder(w).Encode(p)
}

func attachApiHandlers(hub *Hub) {
	// returns a new room id
	http.HandleFunc("/generate_room", func(w http.ResponseWriter, r *http.Request) {
		handleGenerateRoom(hub, w, r)
	})

	// // get members
	http.HandleFunc("/get_members", func(w http.ResponseWriter, r *http.Request) {
		handleGetMembers(hub, w, r)
	})

	// validate room ID
	http.HandleFunc("/validate_room", func(w http.ResponseWriter, r *http.Request) {
		handleValidateRoom(hub, w, r)
	})

	http.HandleFunc("/generate_client_id", func(w http.ResponseWriter, r *http.Request) {
		handleGenerateClientId(hub, w, r)
	})

	http.HandleFunc("/get_all_story_lines", func(w http.ResponseWriter, r *http.Request) {
		handleGetAllStoryLines(hub, w, r)
	})
}
