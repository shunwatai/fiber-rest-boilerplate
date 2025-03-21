package ws

import (
	"encoding/json"
	"fmt"
	logger "golang-api-starter/internal/helper/logger/zap_log"
	"golang-api-starter/internal/modules/groupUser"
	"golang-api-starter/internal/modules/user"
	"log"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func ServeCurrLoginUserWs(router fiber.Router) error {
	onlineUserHub := newOnlineUserHub()
	go onlineUserHub.run()

	router.Get("/curr-users", websocket.New(func(c *websocket.Conn) {
		// Read user details from the JWT claims
		claims := c.Locals("claims").(jwt.MapClaims)
		currLoginUserId := claims["userId"]
		logger.Debugf("???????????????currLoginUserId: %+v", currLoginUserId)
		client := &Client{}
		defer func() {
			client.hub.unregister <- client
			logger.Debugf("here>>>>>>>>>>>>>>>")
			client.conn.Close() // Close the connection
			logger.Debugf("here>>>>>>>>>>>>>>>")
		}()

		done := make(chan struct{}) // Channel to signal when to stop broadcasting
		currUser, err := user.Srvc.GetById(map[string]interface{}{"id": currLoginUserId})
		if err != nil {
			logger.Errorf("failed to fetch curr login user by id: %+v", currLoginUserId)
			c.Close()
			done <- struct{}{}
		} else {
			currUser[0].Password = nil
			// currLoginUsers.Store(currUser[0].GetId(), currUser[0])
			client = &Client{hub: onlineUserHub, conn: c, send: make(chan struct{}), user: currUser[0]}
			client.hub.register <- client

			// Start a goroutine to listen for incoming messages
			go func(c *Client) {
				for {
					if _, _, err := c.conn.ReadMessage(); err != nil {
						log.Println("read:", err)
						c.hub.unregister <- c
						done <- struct{}{}

						break
					}
				}
			}(client)
		}

		for {
			select {
			case <-done:
				return // Exit the broadcasting loop
			}
		}
	}))

	return nil
}

func ServeCurrLoginUserWsOld(router fiber.Router) error {
	// currLoginUsers := map[string]*groupUser.User{}
	var currLoginUsers sync.Map
	// clients := make(map[*websocket.Conn]bool)              // Track connected clients
	var clients sync.Map
	broadcast := make(chan map[string]*groupUser.User, 10) // Broadcast channel

	router.Get("/ws/curr-users", websocket.New(func(c *websocket.Conn) {
		// Read user details from the JWT claims
		claims := c.Locals("claims").(jwt.MapClaims)
		currLoginUserId := claims["userId"]

		currUser, err := user.Srvc.GetById(map[string]interface{}{"id": currLoginUserId})
		if err != nil {
			logger.Errorf("failed to fetch curr login user by id: %+v", currLoginUserId)
		} else {
			currUser[0].Password = nil
			fmt.Println("new user connected")
			currLoginUsers.Store(currUser[0].GetId(), currUser[0])
		}

		// Handle disconnection
		done := make(chan struct{}) // Channel to signal when to stop broadcasting
		defer func() {
			logger.Debugf("here>>>>>>>>>>>>>>>")
			close(done) // Signal the broadcasting loop to stop
			logger.Debugf("here>>>>>>>>>>>>>>>")
			clients.Delete(c)
			logger.Debugf("here>>>>>>>>>>>>>>>")
			if len(currUser) > 0 {
				currLoginUsers.Delete(currUser[0].GetId())
				broadcast <- getCurrentUsers(&currLoginUsers) // Notify remaining clients
			}
			logger.Debugf("here>>>>>>>>>>>>>>>")
			c.Close() // Close the connection
			logger.Debugf("here>>>>>>>>>>>>>>>")
		}()

		// Add the new client to the clients map
		clients.Store(c, true)

		// Notify all clients about the new user
		broadcast <- getCurrentUsers(&currLoginUsers)

		// Start a goroutine to listen for incoming messages
		go func(c *websocket.Conn) {
			for {
				if _, _, err := c.ReadMessage(); err != nil {
					log.Println("read:", err)
					done <- struct{}{}

					// clients.Delete(c)
					// if len(currUser) > 0 {
					// 	currLoginUsers.Delete(currUser[0].GetId())
					// 	broadcast <- getCurrentUsers(&currLoginUsers) // Notify remaining clients
					// }
					// c.Close() // Close the connection
					break
				}
			}
		}(c)

		// Handle broadcasting
		for {
			select {
			case users := <-broadcast:
				// Convert currLoginUsers to JSON
				resp, err := json.Marshal(users)
				if err != nil {
					logger.Errorf("failed to Marshal, err: %s", err.Error())
					continue
				}

				// Send the message to all connected clients
				clients.Range(func(key, value interface{}) bool {
					if err := key.(*websocket.Conn).WriteMessage(websocket.TextMessage, resp); err != nil {
						log.Println("write:", err)
						key.(*websocket.Conn).Close()
						clients.Delete(key.(*websocket.Conn))
						done <- struct{}{}
					}
					return true // continue iteration
				})
			case <-done:
				return // Exit the broadcasting loop
			}
		}
	}))

	return nil
}

// Helper function to get a copy of currLoginUsers safely
func getCurrentUsers(currLoginUsers *sync.Map) map[string]*groupUser.User {
	usersCopy := make(map[string]*groupUser.User)
	currLoginUsers.Range(func(key, value interface{}) bool {
		usersCopy[key.(string)] = value.(*groupUser.User)
		return true // continue iteration
	})
	return usersCopy
}
