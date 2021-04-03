package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const PORT = ":8080"

var addr = flag.String("addr", PORT, "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/room" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

// func main() {
// 	fmt.Println("Starting server")
// 	flag.Parse()
// 	hub := newHub()
// 	go hub.run()
// 	http.HandleFunc("/", serveHome)
// 	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
// 		serveWs(hub, w, r)
// 	})
// 	err := http.ListenAndServe(*addr, nil)
// 	if err != nil {
// 		log.Fatal("ListenAndServe: ", err)
// 	}
// 	fmt.Println("Listening")
// }

var rooms = make(map[uuid.UUID]Room)

func main() {
	fmt.Printf("Starting server on port %s", PORT)
	// move me later!!
	l, _ := zap.NewProduction()
	defer l.Sync()

	flag.Parse()

	registerRoutes(l)

	fmt.Println("Listening")
}

type RegisterRoomRequest struct {
	items []string
	owner string
}

type RegisterRoomResponse struct {
	ID uuid.UUID "JSON:id"
}

func registerRoom(w http.ResponseWriter, r *http.Request, l *zap.Logger) {
	var b = &RegisterRoomRequest{}
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := uuid.NewRandom()
	if err != nil {
		l.Error("Error while creating new room id", zap.Error(err))
	}

	rooms[id] = *newRoom(b.items, b.owner, l)
	// TODO: This should be moved to when the first user visits the room.
	go rooms[id].hub.run()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&RegisterRoomResponse{
		ID: id,
	})
}

func serveRoom(w http.ResponseWriter, r *http.Request, l *zap.Logger) {
	s := strings.TrimPrefix(r.URL.Path, "/ws/")
	id, err := uuid.Parse(s)
	if err != nil {
		l.Error("Unable to parse room uuid", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	room, ok := rooms[id]
	if !ok {
		l.Info("Room does not exist", zap.String("room id", id.String()))
		http.Redirect(w, r, "/room", 418)
	}
	serveWs(room.hub, w, r)
}

func registerRoutes(l *zap.Logger) {
	http.HandleFunc("/room", serveHome)
	http.HandleFunc("/new", func(w http.ResponseWriter, r *http.Request) {
		registerRoom(w, r, l)
	})
	http.HandleFunc("/ws/", func(w http.ResponseWriter, r *http.Request) {
		serveRoom(w, r, l)
	})

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		l.Fatal("ListenAndServe: ", zap.Error(err))
	}
}
