package router

import "github.com/gin-gonic/gin"

import "quenc/middlewares"

import "quenc/apis"

func InitPostRouter(router *gin.Engine) {
	postRouter := router.Group("/post")
	{
		postRouter.POST("/", middlewares.UserAuth(), apis.AddPost)
		postRouter.PATCH("/like/:pid/:condition", middlewares.UserAuth(), apis.LikePost)
		postRouter.PATCH("/detail/:pid", middlewares.UserAuth(), apis.UpdatePost)
		postRouter.DELETE("/:pid", middlewares.UserAuth(), apis.DeletePost)
		postRouter.GET("/all", apis.FindAllPost)
		postRouter.GET("/author/:aid", apis.FindPostByAuthor)
		postRouter.GET("/detail/:pid", apis.FindPostById)
		postRouter.GET("/saved", middlewares.UserAuth(), apis.FindSavedPost)
		postRouter.GET("/array", apis.FindArrayOfPosts)
	}
}
