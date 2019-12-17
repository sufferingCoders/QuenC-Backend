package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// InitRouter -initialise all the routers
func InitRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, []string{"123", "321"})
	})

	router.GET("/ws", func(c *gin.Context) {

		println("Get in the websocket service")

		// Upgrade writer and Reader
		ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)

		defer ws.Close()

		if err != nil {
			log.Printf("The websocket is not working due to the error : %+v \n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}

		for {
			mt, message, err := ws.ReadMessage()

			log.Printf("Get message %+v", string(message))

			if err != nil {
				log.Printf("Error occur: %+v\b", err)
				break
			}

			if string(message) == "ping" {
				message = []byte("pong")
			}

			err = ws.WriteMessage(mt, message)

			if err != nil {
				break
			}
		}
	})

	return router
}
