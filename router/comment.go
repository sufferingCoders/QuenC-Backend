package router

import "github.com/gin-gonic/gin"

import "quenc/middlewares"

import "quenc/apis"

func InitCommentRouter(router *gin.Engine) {
	commentRouter := router.Group("/comment")
	{
		commentRouter.POST("/", middlewares.UserAuth(), apis.AddComment)
		commentRouter.PATCH("/detail/:cid", middlewares.AdminAuth(), apis.UpdateComment)
		commentRouter.PATCH("/like/:cid/:condition", middlewares.UserAuth(), apis.LikeComment)
		commentRouter.DELETE("/:cid", middlewares.AdminAuth(), apis.DeleteComment)
		commentRouter.GET("/post/:pid", apis.FindCommentsByPost)
		commentRouter.GET("/detail/:cid", apis.FindCommentById)
	}
}
