package ws

import (
	logger "golang-api-starter/internal/helper/logger/zap_log"
	"golang-api-starter/internal/modules/user"
	"log"

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
		logger.Debugf("currLoginUserId: %+v", currLoginUserId)
		client := &Client{}

		defer func() {
			client.hub.unregister <- client
			client.conn.Close() // Close the connection
			logger.Debugf("closed>>>>>>>>>>>>>>>")
		}()

		done := make(chan struct{}) // Channel to signal when to stop broadcasting
		// currUser, err := user.Srvc.GetById(map[string]interface{}{"id": currLoginUserId})
		currUser, _ := user.Repo.Get(map[string]interface{}{"id": currLoginUserId})

		// if err != nil {
		if len(currUser) == 0 {
			logger.Errorf("failed to fetch curr login user by id: %+v", currLoginUserId)
			c.Close()
			done <- struct{}{}
		} else {
			currUser[0].Password = nil
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
