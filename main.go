package main

import (
	"quenc/database"
	"quenc/router"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

func main() {

	database.InitDB()
	gin.ForceConsoleColor()
	r := router.InitRouter()
	r.Run()

	/*
		Testing Fetching
	*/

	// oid, err := primitive.ObjectIDFromHex("5e1033d61c9d4400002fa261")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// messages, err := models.FindMessagesForChatRoom(oid, time.Date(2020, time.January, 3, 0, 0, 0, 0, time.Local))

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Printf("messages: %+v", messages)
}
