package main

import (
	"log"
	"net/http"
	
)

func setupAPI() {
	manager := NewManager()
	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	http.HandleFunc("/ws", manager.serveWS)
}

func main() {
	setupAPI()
	log.Fatal(http.ListenAndServe(":8000", nil))
}