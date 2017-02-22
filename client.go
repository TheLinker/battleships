// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 10 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	hubs   []*Hub
	conn   *websocket.Conn
	send   chan []byte
	Player *Player
}


// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Println("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		log.Println("recibido: ", string(message))

		var msg json.RawMessage
		env := Envelope{
			Msg: &msg,
		}

		if err := json.Unmarshal(message, &env); err != nil {
			log.Println("Error: ", err)
			return
		}

		switch strings.ToUpper(env.Type) {
		case "REGISTER":
			var r Registration
			if err := json.Unmarshal(msg, &r); err != nil {
				log.Println("Error: ", err)
				return
			}

			var uname = r.Playername

			log.Println("uname", uname)

            c.Player = CreatePlayer(c, uname)

            // c.hubs = append(c.hubs, GlobalLobby)
            // GlobalLobby.register <- c

			resp := Envelope{
				Type: "RegistrationOK",
				Msg: RegistrationOK{
					Playername: c.Player.Playername,
					Playerhash: c.Player.Playerhash,
				},
			}

			buf, _ := json.Marshal(resp)
			c.send <- buf

            break

        case "CHAT":
            var r Chat

            if err := json.Unmarshal(msg, &r); err != nil {
                log.Println("Error: ", err)
                return
            }

            // if c.hubs != nil {
            //     c.hubs.broadcast <- message
            // }
		}
    }
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
        if c.Player != nil {
            log.Printf("Client killed (%s)\n",  c.Player.Playername)
        } else {
            log.Printf("Client killed (no player)\n")
        }

        for _, h := range(c.hubs) {
           h.unregister <- c
        }

        DeletePlayer(c.Player)

		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{conn: conn, send: make(chan []byte, 256)}

	go client.writePump()
	client.readPump()
}
