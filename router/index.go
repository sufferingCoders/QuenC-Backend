package router

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"quenc/database"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Testing - for testing
type Testing struct {
	ID    primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Email string             `json:"email" bson:"email"`
}

// TestInfo - for adddin the test
type TestInfo struct {
	Email string `json:"email" bson:"email"`
}

// InitRouter -initialise all the routers
func InitRouter() *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, []string{"123", "321"})
	})

	// Inserting the test value
	router.POST("/test", func(c *gin.Context) {
		var testAdding TestInfo
		err := c.ShouldBindJSON(&testAdding)

		if err != nil {
			errStr := fmt.Sprintf("Cannot bind the given info : %+v \n", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"err": errStr,
				"msg": "Cannot bind the given info",
			})
			return
		}

		testingClient := Testing{
			Email: testAdding.Email,
		}

		result, err := database.DB.Collection("test").InsertOne(context.TODO(), testingClient)

		if err != nil {
			errStr := fmt.Sprintf("Can't insert a test due to the error : %+v \n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": errStr,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id": result.InsertedID,
		})

	})

	// Changing the value of a test
	router.PUT("/test/:id", func(c *gin.Context) {

		id := c.Param("id")

		oid, err := primitive.ObjectIDFromHex(id)

		if err != nil {
			errStr := fmt.Sprintf("The given id cannot be transform to oid : %+v \n", err)

			c.JSON(http.StatusInternalServerError, gin.H{
				"err": errStr,
			})
			return
		}

		var testInfo TestInfo
		err = c.ShouldBindJSON(&testInfo)

		if err != nil {
			errStr := fmt.Sprintf("Cannot bind the given info: %+v \n", err)

			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"err": errStr,
			})
			return
		}

		result, err := database.DB.Collection("test").UpdateOne(
			context.TODO(),
			bson.M{"_id": oid},
			bson.M{"$set": bson.M{"email": testInfo.Email}},
		)

		if err != nil {
			errStr := fmt.Sprintf("Cannot update a test due to the error: %+v \n", err)

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": errStr,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"result": result,
		})

	})

	// This for creating the WebSocket and listen to the change stream in the MongoDB
	router.GET("/test/subscribe/:id", func(c *gin.Context) {

		ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)

		defer ws.Close()

		if err != nil {
			errStr := fmt.Sprintf("The websocket is not working due to the error: %+v \n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": errStr,
			})
			return
		}

		id := c.Param("id")

		oid, err := primitive.ObjectIDFromHex(id)

		if err != nil {
			errStr := fmt.Sprintf("The given id cannot be transform to oid: %+v \n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"err": errStr,
			})
			return
		}

		pipeline := mongo.Pipeline{bson.D{{"$match", bson.D{{"fullDocument._id", oid}}}}}

		collectionStream, err := database.DB.Collection("test").Watch(context.TODO(), pipeline, options.ChangeStream().SetFullDocument(options.UpdateLookup))

		if err != nil {
			errStr := fmt.Sprintf("Cannot get the stream: %+v \n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": errStr,
			})
			return
		}

		defer collectionStream.Close(context.TODO())

		// ws.WriteMessage(websocket.TextMessage, []byte("hello"))

		for {
			ok := collectionStream.Next(context.TODO())
			if ok {
				next := collectionStream.Current

				var m map[string]interface{}

				err := bson.Unmarshal(next, &m)
				if err != nil {
					log.Print(err)
					break
				}

				// fmt.Printf("map: %+v", m["fullDocument"])

				// if m["fullDocument"].(map[string]interface{})["email"] == "123" {
				// 	ws.WriteMessage(websocket.TextMessage, []byte("You gave me 123"))
				// 	fmt.Printf("Send out 123")

				// }

				err = ws.WriteJSON(m["fullDocument"].(map[string]interface{}))
				if err != nil {
					log.Print(err)
					break
				}

			}

		}

		fmt.Println("Send out the message_3")

	})

	// This for changing the filed in the MongoDB

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

	InitUserRouter(router)
	InitReportRouter(router)
	InitPostCategoryRouter(router)
	InitPostRouter(router)
	InitCommentRouter(router)

	return router
}
