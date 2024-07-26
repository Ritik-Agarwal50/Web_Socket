package main

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	connection *websocket.Conn
	manager    *Manager
	egress     chan Event
}
type ClientList map[*Client]bool

func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		egress:     make(chan Event),
	}
}

func (c *Client) readMessage() {
	defer func() {
		c.manager.removeClient(c)
	}()
	for {
		_, payload, err := c.connection.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure, websocket.CloseGoingAway) {
				log.Println("error: %v", err)
			}
			break
		}
		var request Event
		if err := json.Unmarshal(payload, &request); err != nil {
			log.Println("error in marhsalling: %v", err)
			break
		}
		if err := c.manager.routeEvent(request, c); err != nil {
			log.Println("error in routing event: %v", err)
		}
	}
}

func (c *Client) writeMessage() {
	defer func() {
		c.manager.removeClient(c)
	}()

	for {
		select {
		case message, ok := <-c.egress:
			if !ok {
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("connection is closed or blocked: %v", err)
				}
				return
			}
			data, err := json.Marshal(message)
			if err != nil {
				log.Println(err)
				return
			}
			if err := c.connection.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Println(err)
			}
			log.Println("send Message")
		}
	}
}
