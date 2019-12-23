package apis

import (
	"context"
	"fmt"
	"net/http"
	"quenc/database"
	"quenc/models"
	"quenc/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

)

func AddPost(c *gin.Context) {

	var post models.Post
	var err error

	if err = c.ShouldBindJSON(&post); err != nil {
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
	post.AuthorDomain = user.Domain
	post.AuthorGender = user.Gender
	post.Author = user.ID.Hex()
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()

	InsertedID, err := models.AddPost(&post)
	if err != nil {
		errStr := fmt.Sprintf("Cannot add this post: %+v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": errStr,
		})
		return
	}

	post.ID = InsertedID.(primitive.ObjectID)

	c.JSON(http.StatusOK, gin.H{
		"post": post,
	})
}

func UpdatePost(c *gin.Context) {

	var err error
	var result *mongo.UpdateResult
	var updateFields map[string]interface{}
	pid := c.Param("pid")

	// Remove some fields so it can be modified

	delete(updateFields, "_id")
	delete(updateFields, "createdAt")
	delete(updateFields, "authorGender")
	delete(updateFields, "authorDomain")
	delete(updateFields, "author")

	if err = c.ShouldBind(&updateFields); err != nil {
		errStr := fmt.Sprintf("Cannot bind the input json: %+v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err":          errStr,
			"updateFields": updateFields,
			"pid":          pid,
		})
		return
	}

	user := utils.GetUserFromContext(c)
	if user == nil {
		return
	}

	// Only Admin and Author can update the post
	pOID := utils.GetOID(pid, c)
	if pOID == nil {
		return
	}

	if user.Role == 0 {
		result, err = models.UpdatePostByOID(*pOID, updateFields)
	} else {
		result, err = database.PostCollection.UpdateOne(context.TODO(),
			bson.M{"_id": pOID, "author": user.ID},
			bson.M{"$set": updateFields},
		)
	}

	if err != nil {
		errStr := fmt.Sprintf("Cannot update the Post with Given User: %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"err":          errStr,
			"user":         user,
			"updateFields": updateFields,
			"pid":          pid,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"result":       result,
		"updateFields": updateFields,
		"pid":          pid,
	})
}

func DeletePost(c *gin.Context) {
	var err error
	pid := c.Param("pid")
	pOID, err := primitive.ObjectIDFromHex(pid)
	if err != nil {
		errStr := fmt.Sprintf("Cannot transfrom the given id to ObjectId: %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": errStr,
			"pid": pid,
		})
	}
	user := utils.GetUserFromContext(c)
	if user == nil {
		return
	}

	// Only Admin and Author Can delete the post

	if user.Role == 0 {
		err = models.DeletePostByOID(pOID)
	} else {
		_, err = database.PostCollection.DeleteOne(context.TODO(),
			bson.M{"_id": pOID, "author": user.ID},
		)
	}

	if err != nil {
		errStr := fmt.Sprintf("Cannot transfrom the given id to ObjectId: %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": errStr,
			"pid": pid,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"pid": pid,
	})
}

func FindAllPost(c *gin.Context) {

	findOption := options.Find()
	err := utils.SetupFindOptions(findOption, c)

	if err != nil {
		return
	}
	posts, err := models.FindPosts(bson.M{}, findOption)
	if err != nil {
		errStr := fmt.Sprintf("Cannot retreive the posts: %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": errStr,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
	})
}

func FindPostById(c *gin.Context) {

	pid := c.Param("pid")
	pOID := utils.GetOID(pid, c)
	if pOID == nil {
		return
	}
	post, err := models.FindPostByOID(*pOID)
	if err != nil {
		errStr := fmt.Sprintf("Cannot retreive the post: %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": errStr,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"post": post,
	})
}

func FindPostByAuthor(c *gin.Context) {
	aid := c.Param("aid")

	aOID := utils.GetOID(aid, c)
	if aOID == nil {
		return
	}

	findOption := options.Find()
	err := utils.SetupFindOptions(findOption, c)
	if err != nil {
		return
	}

	posts, err := models.FindPostByAuthor(*aOID, findOption)
	if err != nil {
		errStr := fmt.Sprintf("Cannot retreive the post: %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": errStr,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
	})
}

func FindSavedPost(c *gin.Context) {
	user := utils.GetUserFromContext(c)
	if user == nil {
		return
	}

	findOption := options.Find()
	utils.SetupFindOptions(findOption, c)

	// makin the save post to ObjectID

	savedOIDs := []*primitive.ObjectID{}
	for _, savedId := range user.SavedPosts {
		oid := utils.GetOID(savedId, c)
		if oid == nil {
			return
		}

		savedOIDs = append(savedOIDs, oid)
	}

	posts, err := models.FindPosts(bson.M{"_id": bson.M{"$in": savedOIDs}}, findOption)

	if err != nil {
		errStr := fmt.Sprintf("Cannot find the SavedPosts: %+v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"err": errStr,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
	})
}

func FindArrayOfPosts(c *gin.Context) {

	postsStr := c.Param("posts")

	postsStrIDs := strings.Split(postsStr, ",")

	postsOID := []*primitive.ObjectID{}

	for _, id := range postsStrIDs {
		oid := utils.GetOID(id, c)
		if oid == nil {
			return
		}

		postsOID = append(postsOID, oid)
	}

	findOption := options.Find()
	utils.SetupFindOptions(findOption, c)

	// makin the save post to ObjectID

	posts, err := models.FindPosts(bson.M{"_id": bson.M{"$in": postsOID}}, findOption)

	if err != nil {
		errStr := fmt.Sprintf("Cannot find the posts : %+v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"err": errStr,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
	})
}
