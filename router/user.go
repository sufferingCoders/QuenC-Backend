package router

import "github.com/gin-gonic/gin"

import "quenc/apis"

import "quenc/middlewares"

func InitUserRouter(router *gin.Engine) {
	userRouter := router.Group("/user")
	{
		userRouter.POST("/signup", apis.SingupUser)
		userRouter.POST("/login", apis.LoginUser)
		userRouter.POST("/auto-login", middlewares.UserAuth(), apis.TokenAutoLogin)
		userRouter.GET("/send-verification-email", middlewares.UserAuth(), apis.SendVerificationEmailForUser)
		userRouter.GET("/email/activate/:uid", apis.ActivateUserEmail)
		userRouter.PATCH("/detail/:uid", middlewares.UserAuth(), apis.UpdateUser)
		userRouter.PATCH("/friends/:id/:condition", middlewares.UserAuth(), apis.ToggleFunc("friends"))
		userRouter.PATCH("/chat-rooms/:id/:condition", middlewares.UserAuth(), apis.ToggleFunc("chatRooms"))
		userRouter.PATCH("/like-posts/:id/:condition", middlewares.UserAuth(), apis.ToggleFunc("likePosts"))
		userRouter.PATCH("/like-comments/:id/:condition", middlewares.UserAuth(), apis.ToggleFunc("likeComments"))
		userRouter.PATCH("/saved-posts/:id/:condition", middlewares.UserAuth(), apis.ToggleFunc("savedPosts"))
		userRouter.PATCH("/block-post/:id/:condition", middlewares.UserAuth(), apis.ToggleFunc("blockedPosts"))
		userRouter.PATCH("/block-user/:id/:condition", middlewares.UserAuth(), apis.ToggleFunc("blockedUsers"))
		userRouter.GET("/subsrible", middlewares.UserAuth(), apis.SubscribeUser)
	}
}
