package routers

import "github.com/gin-gonic/gin"

// InitRouter -initialise all the routers
func InitRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, []string{"123", "321"})
	})

	return router
}
