package apis

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"quenc/models"
	"quenc/utils"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

)

// Need a login Auth
func AddChatRoom(c *gin.Context) {
	var chatRoom models.ChatRoomAdding
	var err error

	if err = c.ShouldBindJSON(&chatRoom); err != nil {
		errStr := fmt.Sprintf("Cannot bind the input json: %+v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": errStr,
		})
		return
	}

	InsertedID, err := models.AddChatRoom(&chatRoom)

	if err != nil {
		errStr := fmt.Sprintf("Cannot add this ChatRoom: %+v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": errStr,
		})
		return
	}

	chatRoom.ID = InsertedID.(primitive.ObjectID)

	c.JSON(http.StatusOK, gin.H{
		"chatRoom": chatRoom,
	})
}

func AddMessageToChatRoom(c *gin.Context) {
	var message models.Message
	var err error

	// getting the chat roomOID

	rid := c.Param("rid")
	rOID := utils.GetOID(rid, c)
	if rOID == nil {
		return
	}

	if err = c.ShouldBind(&message); err != nil {
		errStr := fmt.Sprintf("Cannot bind the input json: %+v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": errStr,
		})
		return
	}

	user := utils.GetUserFromContext(c)
	if user == nil {
		return
	}

	message.Author = user.ID
	message.CreatedAt = time.Now()
	message.LikedBy = []primitive.ObjectID{}
	message.ReadBy = []primitive.ObjectID{} // author doens't need to be here
	message.ID = primitive.NewObjectID()

	result, err := models.AddMessageToChatRoom(*rOID, message)

	if err != nil {
		errStr := fmt.Sprintf("Cannot add this message : %+v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"err": errStr,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result":  result,
		"id":      message.ID,
		"message": message,
	})
}

func UpdateChatRoom(c *gin.Context) {

	var err error
	var result *mongo.UpdateResult
	var updateFields map[string]interface{}
	rid := c.Param("rid") // Get the room id

	if err = c.ShouldBind(&updateFields); err != nil {
		errStr := fmt.Sprintf("Cannot bind the input json: %+v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err":          errStr,
			"updateFields": updateFields,
			"rid":          rid,
		})
		return
	}

	delete(updateFields, "_id")
	delete(updateFields, "createdAt")
	delete(updateFields, "isGroup")

	// Only Admin and Author can update the Comment
	rOID := utils.GetOID(rid, c)
	if rOID == nil {
		return
	}

	result, err = models.UpdateChatRoomByOID(*rOID, updateFields)

	if err != nil {
		errStr := fmt.Sprintf("Cannot update the ChatRoom with Given User: %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"err":          errStr,
			"updateFields": updateFields,
			"rid":          rid,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"result":       result,
		"updateFields": updateFields,
		"rid":          rid,
	})
}

func DeleteChatRoom(c *gin.Context) {
	var err error
	rid := c.Param("rid")

	rOID := utils.GetOID(rid, c)

	if rOID == nil {
		return
	}

	// Only Admin and Author Can delete the Comment

	err = models.DeleteChatRoomByOID(*rOID)

	if err != nil {
		errStr := fmt.Sprintf("Cannot delete this chatRoom: %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": errStr,
			"rid": rid,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"rid": rid,
	})
}

// Find User ChatRoom Detail

// login Require
// This need to be watching? get then watching

// then how to get message?
func FindUserChatRoomDetailWithLastMessages(c *gin.Context) {
	user := utils.GetUserFromContext(c)
	if user == nil {
		return
	}

	chatRooms, err := models.FindChatRoomDetailWithLastMessage(&[]bson.M{
		bson.M{
			"$match": bson.M{
				"$in": bson.A{"$_id", user.ChatRooms},
			},
		},
	},
		-1,
		-1,
	)

	if err != nil {
		errStr := fmt.Sprintf("Cannot retreive the chatRooms: %+v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"err": errStr,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"chatRooms": chatRooms,
	})
}

func FindMessagesForRoom(c *gin.Context) {
	// you have to be one of the member

	var err error

	user := utils.GetUserFromContext(c)
	if user == nil {
		return
	}

	rid := c.Param("rid")
	rOID := utils.GetOID(rid, c)
	if rOID == nil {
		return
	}

	var sOID *primitive.ObjectID
	sid := c.Query("sid")
	if strings.TrimSpace(sid) != "" {
		sOID = utils.GetOID(sid, c)
		if sOID == nil {
			return
		}
	}

	numInt := 50
	num := c.Query("num")
	if strings.TrimSpace(num) != "" {
		numInt, err = strconv.Atoi(num)
		if err != nil {
			errStr := fmt.Sprintf("Cannot transform the num to int: %+v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"err": errStr,
			})
			return
		}
	}

	messages, err := models.FindMessagesForChatRoomByStartOID(user.ID, *rOID, sOID, numInt)

	if err != nil {
		errStr := fmt.Sprintf("Cannot find the messages: %+v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"err": errStr,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messages": messages,
	})
}

// Do the subscribing here

// 若User的 ChatRoom Field有改變的話, 那就重新Subscribe一次
func SubscribeUserChatRoomDetail(c *gin.Context) {

	upGrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	user := utils.GetUserFromContext(c)
	if user == nil {
		return
	}

	stream, err := models.WatchChatRooms(user.ChatRooms)

	if err != nil {
		errStr := fmt.Sprintf("Cannot get the stream: %+v \n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errStr,
		})
		return
	}

	if stream != nil {
		defer stream.Close(context.TODO())
	}

	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		errStr := fmt.Sprintf("The websocket is not working due to the error: %+v \n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errStr,
		})
		return
	}

	if stream != nil {
		defer ws.Close()
	}

	for {
		ok := stream.Next(context.TODO())
		if ok {
			next := stream.Current
			// err = ws.WriteJSON(next) // 是否要直接傳送 不用重複encode & decode

			var m map[string]interface{}

			err := bson.Unmarshal(next, &m)
			if err != nil {
				log.Print(err)
			}
			err = ws.WriteJSON(m)
			if err != nil {
				log.Print(err)
			}
		}
	}
}

// 先取得基礎資訊再追蹤流?
// 基礎資訊若有斷層?
// 先做全觀察
// 觀察時也回傳FullDocuemnt? // First try this
