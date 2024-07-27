package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	websocketUpgrader = websocket.Upgrader{
		CheckOrigin:     checkOrigin,
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

var ErrEventNotSupported = errors.New("this event type is not supported")

type Manager struct {
	clients ClientList
	sync.RWMutex
	handlers map[string]EventHandler
}

func NewManager() *Manager {
	m := &Manager{
		clients:  make(ClientList),
		handlers: make(map[string]EventHandler),
	}
	m.setupEventHandlers()
	return m
}

func (m *Manager) setupEventHandlers() {
	m.handlers[EventSendMessage] = func(e Event, c *Client) error {
		fmt.Println(e)
		return nil
	}
}

func (m *Manager) routeEvent(event Event, c *Client) error {
	if handler, ok := m.handlers[event.Type]; ok {
		if err := handler(event, c); err != nil {
			return err
		}
		return nil
	} else {
		return ErrEventNotSupported
	}
}

func (m *Manager) serveWS(w http.ResponseWriter, r *http.Request) {
	log.Println("New Connection")
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := NewClient(conn, m)
	m.addClient(client)
	go client.writeMessages()
	go client.readMessages()
}

func (m *Manager) addClient(client *Client) {
	m.Lock()
	defer m.Unlock()
	m.clients[client] = true

}

func (m *Manager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()
	if _, ok := m.clients[client]; ok {
		client.connection.Close()
		delete(m.clients, client)
	}
}

// func (c *Client) readMessages() {
// 	defer func() {
// 		c.manager.removeClient(c)
// 	}()

// 	for {
// 		messageType, payload, err := c.connection.ReadMessage()
// 		if err != nil {
// 			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
// 				log.Printf("error: %v", err)
// 			}
// 			break
// 		}
// 		log.Println("Message Type: ", messageType)
// 		log.Println("Payload", string(payload))

// 		for wsclient := range c.manager.clients {
// 			wsclient.egress <- payload
// 		}
// 	}
// }

// func (c *Client) WriteMessages() {
// 	defer func() {
// 		c.manager.removeClient(c)
// 	}()

// 	for {
// 		select {
// 		case message, ok := <-c.egress:
// 			if !ok {
// 				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
// 					log.Println("connection closed: ", err)
// 				}
// 				return
// 			}
// 			if err := c.connection.WriteMessage(websocket.TextMessage,message); err != nil {
// 				log.Println(err)
// 			}
// 			log.Println("sent message")
// 		}
// 	}
// }

func checkOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	switch origin {
	case "http://localhost:8080":
		return true
	default:
		return false
	}
}
