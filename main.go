package main

import (
	"log"
	"net/http"
)
func setupAPI() {
	http.Handle("/", http.FileServer(http.Dir("./frontend")))
}

func main() {
	setupAPI()
	log.Fatal(http.ListenAndServe(":8000", nil))
}