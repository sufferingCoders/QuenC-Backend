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
	err := c.ShouldBindJSON(&postCategory)

	if err != nil {
		errStr := fmt.Sprintf("Cannot bind the input json: %+v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": errStr,
		})
		return
	}

	user := utils.GetUserFromContext(c)

	// Only admin can add post category

	if !user.IsAmin() {
		errStr := fmt.Sprintf("Only Admin can add post category")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": errStr,
		})
		return
	}

	InsertedID, err := models.AddPostCategory(&postCategory)

	postCategory.ID = InsertedID.(primitive.ObjectID)

	c.JSON(http.StatusOK, gin.H{
		"postCategory": postCategory,
	})

}

func UpdatePostCategory(c *gin.Context) {
	var updateFields map[string]interface{}

	cid := c.Param("cid")
	err := c.ShouldBindJSON(&updateFields)
	if err != nil {
		errStr := fmt.Sprintf("Cannot bind the input json: %+v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": errStr,
		})
		return
	}

	cOID := utils.GetOID(cid, c)

	user := utils.GetUserFromContext(c)

	if !user.IsAmin() {
		errStr := fmt.Sprintf("Only Admin can update post category")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": errStr,
		})
		return
	}

	result, err := models.UpdatePostCategoryByOID(*cOID, updateFields)

	if err != nil {
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
	cOID := utils.GetOID(cid, c)

	user := utils.GetUserFromContext(c)

	if !user.IsAmin() {
		errStr := fmt.Sprintf("Only Admin can delete post category")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": errStr,
		})
		return
	}

	err := models.DeletePostCategoryByOID(*cOID)

	if err != nil {
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
	utils.SetupFindOptions(findOption, c)

	postCategories, err := models.FindAllPostCategorys(findOption)

	if err != nil {
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
	cOID := utils.GetOID(cid, c)

	postCategory, err := models.FindPostCategoryByOID(*cOID)

	if err != nil {
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
