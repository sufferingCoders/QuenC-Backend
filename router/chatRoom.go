package router

import "github.com/gin-gonic/gin"

import "quenc/middlewares"

import "quenc/apis"

func InitChatRoomRouter(router *gin.Engine) {
	chatRoomRouter := router.Group("/chat-room")
	{
		chatRoomRouter.POST("/", middlewares.UserAuth(), apis.AddChatRoom)
		chatRoomRouter.PATCH("/:rid", middlewares.UserAuth(), apis.UpdateChatRoom)
		chatRoomRouter.DELETE("/:rid", middlewares.UserAuth(), apis.DeleteChatRoom)
		chatRoomRouter.GET("/rooms", middlewares.UserAuth(), apis.FindUserChatRoomDetailWithoutMessages)
		

		chatRoomRouter.GET("/user/subscribe", middlewares.UserAuth(), apis.SubscribeUserChatRoomDetail)
	}
}
