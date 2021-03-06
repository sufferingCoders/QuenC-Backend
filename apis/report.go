package apis

import (
	"fmt"
	"net/http"
	"quenc/models"
	"quenc/utils"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

)

// 若回傳nil 直接Return
func AddReport(c *gin.Context) {
	var report models.ReportAdding
	var err error

	if err = c.ShouldBindJSON(&report); err != nil {
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

	switch report.ReportTarget {
	case 0:
		// find post
		post, err := models.FindSinglePostWithDetail(report.ReportID)
		if err != nil {
			errStr := fmt.Sprintf("Cannot find this Report post: %+v", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"err": errStr,
			})
			return
		}

		retreivedMap, err := utils.StructToMap(post)

		if err != nil {
			errStr := fmt.Sprintf("Cannot convert this post to map: %+v", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"err": errStr,
			})
			return
		}

		report.ReportObject = retreivedMap

	case 1:
		// find comment
		comment, err := models.FindCommentByOID(report.ReportID)
		if err != nil {
			errStr := fmt.Sprintf("Cannot find this Report comment: %+v", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"err": errStr,
			})
			return
		}

		retreivedMap, err := utils.StructToMap(comment)

		if err != nil {
			errStr := fmt.Sprintf("Cannot convert this comment to map: %+v", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"err": errStr,
			})
			return
		}

		report.ReportObject = retreivedMap

		break
	case 2:
		// find chat
		room, err := models.FindChatRoomByOID(report.ReportID)
		if err != nil {
			errStr := fmt.Sprintf("Cannot find this Report chat room: %+v", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"err": errStr,
			})
			return
		}

		retreivedMap, err := utils.StructToMap(room)

		if err != nil {
			errStr := fmt.Sprintf("Cannot convert this room to map: %+v", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"err": errStr,
			})
			return
		}

		report.ReportObject = retreivedMap

		break
	}

	report.Author = user.ID
	report.CreatedAt = time.Now()
	report.Solve = false

	InsertedID, err := models.AddReport(&report)
	if err != nil {
		errStr := fmt.Sprintf("Cannot add this Report: %+v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": errStr,
		})
		return
	}

	report.ID = InsertedID.(primitive.ObjectID)
	c.JSON(http.StatusOK, gin.H{
		"report": report,
	})
}

func UpdateReport(c *gin.Context) {

	var err error
	var result *mongo.UpdateResult
	var updateFields map[string]interface{}
	rid := c.Param("rid")

	if err = c.ShouldBind(&updateFields); err != nil {
		errStr := fmt.Sprintf("Cannot bind the input json: %+v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err":          errStr,
			"updateFields": updateFields,
			"pid":          rid,
		})
		return
	}

	user := utils.GetUserFromContext(c)
	if user == nil {
		return
	}

	// Only Admin and Author can update the Report
	rOID := utils.GetOID(rid, c)
	if rOID == nil {
		return
	}

	delete(updateFields, "_id")
	delete(updateFields, "author")
	delete(updateFields, "createdAt")

	result, err = models.UpdateReportByOID(*rOID, updateFields)

	if err != nil {
		errStr := fmt.Sprintf("Cannot update the Report with Given User: %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"err":          errStr,
			"user":         user,
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

func DeleteReport(c *gin.Context) {
	var err error
	rid := c.Param("rid")
	rOID := utils.GetOID(rid, c)
	if rOID == nil {
		return
	}

	// Only Admin and Author Can delete the Report

	err = models.DeleteReportByOID(*rOID)

	if err != nil {
		errStr := fmt.Sprintf("Cannot transfrom the given id to ObjectId: %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": errStr,
			"pid": rid,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"pid": rid,
	})
}

func FindAllReports(c *gin.Context) {
	findOption := options.Find()
	err := utils.SetupFindOptions(findOption, c)
	if err != nil {
		return
	}
	reports, err := models.FindReports(bson.M{}, findOption)

	if err != nil {
		errStr := fmt.Sprintf("Cannot find the reports: %+v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"err": errStr,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"reports": reports,
	})
}

// Admin Middleware
func FindSingleReport(c *gin.Context) {
	rid := c.Param("rid")

	rOID := utils.GetOID(rid, c)
	if rOID == nil {
		return
	}

	report, err := models.FindSingleReportWithDetail(*rOID)

	if err != nil {
		errStr := fmt.Sprintf("Cannot find the rsport: %+v", err)
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"err": errStr,
				"rid": rid,
			},
		)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"report": report,
	})
}

// Admin Middleware

func FindReportsForPreview(c *gin.Context) {

	skip, limit, _, err := utils.GetSkipLimitSortFromContext(c)

	if err != nil {
		return
	}

	reports, err := models.FindAllReporstWithPreview(*skip, *limit)

	if err != nil {
		errStr := fmt.Sprintf("Cannot fetch the reports: %+v", reports)
		c.JSON(http.StatusInternalServerError, gin.H{
			"err":   errStr,
			"limit": limit,
			"skip":  skip,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"reports": reports,
	})

}

func FindReportsWithDetail(c *gin.Context) {

	skip, limit, _, err := utils.GetSkipLimitSortFromContext(c)

	if err != nil {
		return
	}

	reports, err := models.FindAllReporstWithDetail(*skip, *limit)

	if err != nil {
		errStr := fmt.Sprintf("Cannot fetch the reports: %+v", reports)
		c.JSON(http.StatusInternalServerError, gin.H{
			"err":   errStr,
			"limit": limit,
			"skip":  skip,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"reports": reports,
	})

}
