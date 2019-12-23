package apis

import (
	"fmt"
	"net/http"
	"quenc/models"
	"quenc/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

)

func AddComment(c *gin.Context) {

	var comment models.Comment
	var err error

	if err = c.ShouldBindJSON(&comment); err != nil {
		errStr := fmt.Sprintf("Cannot bind the input json: %+v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": errStr,
		})
		return
	}
	InsertedID, err := models.AddComment(&comment)
	if err != nil {
		errStr := fmt.Sprintf("Cannot add this Comment: %+v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": errStr,
		})
		return
	}

	comment.ID = InsertedID.(primitive.ObjectID)

	c.JSON(http.StatusOK, gin.H{
		"comment": comment,
	})
}

func UpdateComment(c *gin.Context) {

	var err error
	var result *mongo.UpdateResult
	var updateFields map[string]interface{}
	cid := c.Param("cid")

	if err = c.ShouldBind(&updateFields); err != nil {
		errStr := fmt.Sprintf("Cannot bind the input json: %+v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err":          errStr,
			"updateFields": updateFields,
			"cid":          cid,
		})
		return
	}

	// Only Admin and Author can update the Comment
	cOID := utils.GetOID(cid, c)
	if cOID == nil {
		return
	}

	result, err = models.UpdateCommentByOID(*cOID, updateFields)

	if err != nil {
		errStr := fmt.Sprintf("Cannot update the Comment with Given User: %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"err":          errStr,
			"updateFields": updateFields,
			"cid":          cid,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"result":       result,
		"updateFields": updateFields,
		"cid":          cid,
	})
}

func DeleteComment(c *gin.Context) {
	var err error
	cid := c.Param("cid")
	pOID, err := primitive.ObjectIDFromHex(cid)
	if err != nil {
		errStr := fmt.Sprintf("Cannot transfrom the given id to ObjectId: %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": errStr,
			"cid": cid,
		})
	}

	// Only Admin and Author Can delete the Comment

	err = models.DeleteCommentByOID(pOID)

	if err != nil {
		errStr := fmt.Sprintf("Cannot transfrom the given id to ObjectId: %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": errStr,
			"cid": cid,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"cid": cid,
	})
}

func FindAllComment(c *gin.Context) {

	findOption := options.Find()
	err := utils.SetupFindOptions(findOption, c)

	if err != nil {
		return
	}
	comments, err := models.FindComments(bson.M{}, findOption)
	if err != nil {
		errStr := fmt.Sprintf("Cannot retreive the Comments: %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": errStr,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"Comments": comments,
	})
}

func FindCommentsByPost(c *gin.Context) {
	pid := c.Param("pid")

	pOID := utils.GetOID(pid, c)
	if pOID == nil {
		return
	}

	findOption := options.Find()
	err := utils.SetupFindOptions(findOption, c)
	if err != nil {
		return
	}

	comments, err := models.FindCommentByPost(*pOID, findOption)
	if err != nil {
		errStr := fmt.Sprintf("Cannot retreive the Comment: %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": errStr,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"comments": comments,
	})
}

func FindCommentById(c *gin.Context) {
	cid := c.Param("cid")
	cOID := utils.GetOID(cid, c)
	if cOID == nil {
		return
	}

	comment, err := models.FindCommentByOID(*cOID)

	if err != nil {
		errStr := fmt.Sprintf("Cannot find the commennt: %+v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"err": errStr,
			"cid": cid,
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"comment": comment,
	})
}
