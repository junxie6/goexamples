package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"time"
)

// client represents a single chatting user
type client struct {
	socket   *websocket.Conn
	send     chan *message
	room     *room
	userData map[string]interface{}
}

func (c *client) read() {
	defer c.socket.Close()

	for {
		var msg *message
		err := c.socket.ReadJSON(&msg)
		if err != nil {
			fmt.Printf("err ReadJSON: %#v\n", err)
			return
		}
		msg.When = time.Now()
		msg.Name = c.userData["name"].(string)
		if avatarURL, ok := c.userData["avatar_url"]; ok {
			msg.AvatarURL = avatarURL.(string)
		}
		c.room.forward <- msg
	}
}

func (c *client) write() {
	defer c.socket.Close()

	for msg := range c.send {
		err := c.socket.WriteJSON(msg)
		if err != nil {
			break
		}
	}
}
