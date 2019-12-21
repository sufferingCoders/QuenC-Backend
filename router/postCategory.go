package router

import "github.com/gin-gonic/gin"

import "quenc/middlewares"

import "quenc/apis"

func InitPostCategoryRouter(router *gin.Engine) {
	postCategoryRouter := router.Group("/post-category")
	{
		postCategoryRouter.POST("/", middlewares.AdminAuth(), apis.AddPostCategory)
		postCategoryRouter.PATCH("/:cid", middlewares.AdminAuth(), apis.UpdatePostCategory)
		postCategoryRouter.DELETE("/:cid", middlewares.AdminAuth(), apis.DeletePostCategoryById)
		postCategoryRouter.GET("/", apis.FindAllPostCategorys)
		postCategoryRouter.GET("/detail/:cid", apis.FindPostCategoryByID)
	}

}
