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
		userRouter.PATCH("/:uid",middlewares.UserAuth(), apis.UpdateUser)
	}
}
