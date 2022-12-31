package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	flag.Parse()

	hub := newHub()
	go hub.run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	attachApiHandlers(hub)

	initializeRedis()

	log.Println("Starting server...")
	log.Println("Listening on port 8080.") 

    err := http.ListenAndServe("0.0.0.0:8080", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
