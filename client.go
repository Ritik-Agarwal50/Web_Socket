package main

import "github.com/gorilla/websocket"

type Client struct {
	connection *websocket.Conn
	manager *Manager
}
type ClientList map[*Client]bool

func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		connection: conn,
		manager: manager,
	}
}