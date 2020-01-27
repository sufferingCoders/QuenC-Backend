package router

import "github.com/gin-gonic/gin"

import "quenc/middlewares"

import "quenc/apis"

func InitChatRoomRouter(router *gin.Engine) {
	chatRoomRouter := router.Group("/chat-room")
	{
		chatRoomRouter.POST("/", middlewares.UserAuth(), apis.AddChatRoom)
		chatRoomRouter.POST("/message/:rid", middlewares.UserAuth(), apis.AddMessageToChatRoom)
		chatRoomRouter.POST("/test/message", middlewares.UserAuth(), apis.TestAddMessageToChatRoom)
		chatRoomRouter.PATCH("/:rid", middlewares.UserAuth(), apis.UpdateChatRoom)
		chatRoomRouter.DELETE("/detail/:rid", middlewares.UserAuth(), apis.DeleteChatRoom)
		chatRoomRouter.GET("/rooms", middlewares.UserAuth(), apis.FindUserChatRoomDetailWithLastMessages)
		chatRoomRouter.GET("/message/:rid", middlewares.UserAuth(), apis.FindMessagesForRoom)
		chatRoomRouter.GET("/user/subscribe", middlewares.UserAuth(), apis.SubscribeUserChatRoomDetail)
		chatRoomRouter.POST("/random/connect", middlewares.UserAuth(), apis.AssignRandomChatRoomForUser)
		chatRoomRouter.GET("/random/room", middlewares.UserAuth(), apis.FindDetailOfRandomRoom)
		chatRoomRouter.GET("/random/message", middlewares.UserAuth(), apis.FindMessageForRandomChatRoom)
		chatRoomRouter.GET("/random/subscribe", middlewares.UserAuth(), apis.SubscribeUserRandomChatRoomDetail)
		chatRoomRouter.DELETE("/random", middlewares.UserAuth(), apis.LeavingRandomChatRoom)
	}
}
