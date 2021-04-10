package socket

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

// clientWrite sends JSON messages to client connection.
func (c *Client) clientWrite() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			jsonStr, err := json.Marshal(message)

			err = c.conn.WriteJSON(string(jsonStr))
			if err != nil {
				log.WithFields(log.Fields{
					"context": "clientWrite WriteJSON",
					"error":   err.Error(),
				}).Error("Error sending message")
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
