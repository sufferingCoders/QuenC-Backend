package apis

import (
	"fmt"
	"net/http"
	"quenc/models"
	"quenc/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AddPostCategory(c *gin.Context) {
	var postCategory models.PostCategory

	if err := c.ShouldBindJSON(&postCategory); err != nil {
		errStr := fmt.Sprintf("Cannot bind the input json: %+v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": errStr,
		})
		return
	}

	if user := utils.GetUserFromContext(c); user == nil {
		return
	}

	// Only admin can add post category

	// Do this in the Middleware
	// if !user.IsAmin() {
	// 	errStr := fmt.Sprintf("Only Admin can add post category")
	// 	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
	// 		"err": errStr,
	// 	})
	// 	return
	// }

	InsertedID, err := models.AddPostCategory(&postCategory)

	postCategory.ID = InsertedID.(primitive.ObjectID)

	c.JSON(http.StatusOK, gin.H{
		"postCategory": postCategory,
	})

}

func UpdatePostCategory(c *gin.Context) {
	var updateFields map[string]interface{}

	cid := c.Param("cid")

	if err := c.ShouldBindJSON(&updateFields); err != nil {
		errStr := fmt.Sprintf("Cannot bind the input json: %+v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": errStr,
		})
		return
	}

	if cOID := utils.GetOID(cid, c); cOID = nil {
		return
	}


	if user := utils.GetUserFromContext(c); user == nil {
		return 
	}

	if !user.IsAmin() {
		errStr := fmt.Sprintf("Only Admin can update post category")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": errStr,
		})
		return
	}

	if result, err := models.UpdatePostCategoryByOID(*cOID, updateFields); err != nil {
		errStr := fmt.Sprintf("Cannot update the category: %+v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": errStr,
		})
		return
	}

	c.JSON(
		http.StatusOK, gin.H{
			"result": result,
		},
	)
}

func DeletePostCategoryById(c *gin.Context) {
	cid := c.Param("cid")
	if cOID := utils.GetOID(cid, c); cOID == nil {
		return 
	}

	if user := utils.GetUserFromContext(c), user == nil {
		return 
	}

	if !user.IsAmin() {
		errStr := fmt.Sprintf("Only Admin can delete post category")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": errStr,
		})
		return
	}

	if err := models.DeletePostCategoryByOID(*cOID); err != nil {
		errStr := fmt.Sprintf("Cannot delete the category: %+v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": errStr,
		})
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"cid": cid,
		},
	)
}

func FindAllPostCategorys(c *gin.Context) {
	findOption := options.Find()
	if err := utils.SetupFindOptions(findOption, c); err != nil {
		return
	}

	if postCategories, err := models.FindAllPostCategorys(findOption); err != nil {
		errStr := fmt.Sprintf("Cannot find the categories: %+v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": errStr,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"postCategories": postCategories,
	})
}

func FindPostCategoryByOID(c *gin.Context) {
	cid := c.Param("cid")
	if cOID := utils.GetOID(cid, c); cOID == nil {
		return 
	}

	

	if postCategory, err := models.FindPostCategoryByOID(*cOID); err != nil {
		errStr := fmt.Sprintf("Cannot find the categories: %+v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": errStr,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"postCategory": postCategory,
	})

}
