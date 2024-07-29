package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

func setupAPI(ctx context.Context) {
	manager := NewManager(ctx)
	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	http.HandleFunc("/login", manager.loginHandler)
	http.HandleFunc("/ws", manager.serveWS)

	http.HandleFunc("/debug", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, len(manager.clients))
	})
}

func main() {
	rootCtx := context.Background()
	ctx, cancel := context.WithCancel(rootCtx)
	defer cancel()
	setupAPI(ctx)
	err := http.ListenAndServeTLS(":8000", "server.crt", "server.key", nil)
	if err != nil {
		log.Fatal("listen and server err:", err)
	}
	log.Fatal(http.ListenAndServe(":8080", nil))
}
