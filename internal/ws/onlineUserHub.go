package ws

import (
	"github.com/gofiber/contrib/websocket"

	logger "golang-api-starter/internal/helper/logger/zap_log"
	"golang-api-starter/internal/modules/groupUser"
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *OnlineUsersHub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	// send chan []byte
	send chan struct{}

	user *groupUser.User
}

type OnlineUsersHub struct {
	// Registered clients.
	// clients map[*Client]bool
	clients map[*Client]groupUser.User

	// Inbound messages from the clients.
	// broadcast chan []byte
	broadcast chan struct{}

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newOnlineUserHub() *OnlineUsersHub {
	return &OnlineUsersHub{
		// broadcast:  make(chan []byte),
		broadcast:  make(chan struct{}, 10),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		// clients:    make(map[*Client]bool),
		clients: make(map[*Client]groupUser.User),
	}
}

func (h *OnlineUsersHub) run() {
	// var onlineUserList sync.Map
	var keyPrefix = "online_user:"
	var onlineUserList = NewOnlineUserList(h)

	for {
		select {
		case client := <-h.register:
			// h.clients[client] = true
			logger.Debugf(">>>>>>>>> new user online, %+v", *client.user)
			h.clients[client] = *client.user
			onlineUserList.Set(keyPrefix+client.user.GetId(), client.user)

			// h.broadcast <- struct{}{}
		case client := <-h.unregister:
			logger.Debugf(">>>>>>>>> user left")
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				// onlineUserList.Delete(client.user.GetId())
				onlineUserList.Del(keyPrefix + client.user.GetId())

				// h.broadcast <- struct{}{}
			}
		// case message := <-h.broadcast:
		case <-h.broadcast:
			userList := onlineUserList.GetList()
			logger.Debugf(">>>>>>>> no. of clients: %+v", len(h.clients))
			for client := range h.clients {
				client.conn.WriteJSON(userList)
			}
		}
	}
}
