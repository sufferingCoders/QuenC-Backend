package router

import "github.com/gin-gonic/gin"

import "quenc/middlewares"

import "quenc/apis"

func InitReportRouter(router *gin.Engine) {
	reportRouter := router.Group("/report")
	{
		reportRouter.POST("/", middlewares.UserAuth(), apis.AddComment)
		reportRouter.PATCH("/:rid", middlewares.AdminAuth(), apis.UpdateReport)
		reportRouter.DELETE("/:rid", middlewares.AdminAuth(), apis.DeleteReport)
		reportRouter.GET("/", middlewares.AdminAuth(), apis.FindReportsForPreview)
		reportRouter.GET("/detail/:rid", middlewares.AdminAuth(), apis.FindSingleReport)
	}

}
