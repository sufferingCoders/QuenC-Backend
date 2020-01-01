package main

import (
	"quenc/database"
	"quenc/router"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"

)

func main() {
	gin.ForceConsoleColor()
	database.InitDB()
	r := router.InitRouter()
	r.Run()
}
