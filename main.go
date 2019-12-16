package main

import (
	"Quenc/database"
	"Quenc/routers"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.ForceConsoleColor()
	database.InitDB()

	r := routers.InitRouter()

	r.Run()
}
