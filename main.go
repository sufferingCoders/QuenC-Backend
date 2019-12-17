package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"quenc/database"
	"quenc/routers"
)

func main() {
	gin.ForceConsoleColor()
	database.InitDB()
	r := routers.InitRouter()
	r.Run()
}
