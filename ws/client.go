package ws

import (
	"net/http"
	"time"

	"fmt"

	"encoding/json"

	"github.com/gorilla/websocket"
	"gopkg.in/mgo.v2/bson"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	maxMessageSize = 8
)

//Client Websocket connections
type Client struct {
	token string
	conn  *websocket.Conn
}

type user struct {
	clients map[*Client]bool
}

//Clean Clean dead connections
func (c user) Clean() {
	for key := range c.clients {
		err := key.conn.WriteMessage(websocket.PingMessage, make([]byte, 0))
		if err != nil {
			delete(c.clients, key)
		}
	}
}

// ServeWs Here is where we establish the connection websocket with the user
func (hub *Hub) ServeWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	connection := bson.NewObjectId().Hex()
	client := &Client{token: connection, conn: conn}
	hub.clients[connection] = client
	client.Listen()
}

//Listen It's just to shoot the event close
func (c *Client) Listen() {
	defer func() {
		recover()
	}()
	open := map[string]interface{}{
		"type":  "websocket/OPEN",
		"token": c.token,
	}
	mensaje, _ := json.Marshal(open)
	c.conn.WriteMessage(websocket.TextMessage, mensaje)
	for {
		if _, _, err := c.conn.ReadMessage(); err != nil {
			break
		}
	}
	hub.broadcast <- action{
		action: "close",
		data: map[string]interface{}{
			"token": c.token,
		},
	}
	println("Conexion cerrada")
}
