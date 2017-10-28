package ws

import (
	"github.com/gorilla/websocket"
)

//Hub Websocket client store
type Hub struct {
	clients   map[string]*Client
	broadcast chan action
}

type action struct {
	action string
	data   map[string]interface{}
}

var hub = &Hub{
	clients:   map[string]*Client{},
	broadcast: make(chan action),
}

//GetHub Obtains the system hub, which is responsible for processing all
//websocket requests and processing them
func GetHub() *Hub {
	return hub
}

//Send Send messages to users
func (hub Hub) Send(conn string, message []byte) error {
	hub.broadcast <- action{
		action: "send",
		data: map[string]interface{}{
			"conn":    conn,
			"message": message,
		},
	}
	return nil
}

//Start d
func (hub Hub) Start() {
	go func() {
		for {
			action := <-hub.broadcast
			switch action.action {
			case "send":
				conn := action.data["conn"].(string)
				message := action.data["message"].([]byte)
				if conection, ok := hub.clients[conn]; ok {
					conection.conn.WriteMessage(websocket.TextMessage, message)
				}
			case "close":
				token := action.data["token"].(string)
				if _, ok := hub.clients[token]; ok {
					delete(hub.clients, token)
				}
			default:
			}
		}
	}()
}
