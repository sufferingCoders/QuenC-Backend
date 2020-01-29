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

	var post models.PostAdding
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
	post.Author = user.ID
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()
	post.Likers = []primitive.ObjectID{}

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

	// if we have category fields, then convert it to ObjectID

	if err = c.ShouldBind(&updateFields); err != nil {
		errStr := fmt.Sprintf("Cannot bind the input json: %+v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err":          errStr,
			"updateFields": updateFields,
			"pid":          pid,
		})
		return
	}

	if cStr, ok := updateFields["category"]; ok {

		categoryOID, err := primitive.ObjectIDFromHex(cStr.(string))

		if err != nil {
			errStr := fmt.Sprintf("Cannot get the category OID: %+v", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"err":          errStr,
				"updateFields": updateFields,
				"pid":          pid,
			})
			return
		}

		updateFields["category"] = categoryOID

	}

	delete(updateFields, "_id")
	delete(updateFields, "createdAt")
	delete(updateFields, "author")
	delete(updateFields, "updatedAt")

	user := utils.GetUserFromContext(c)
	if user == nil {
		return
	}

	// Only Admin and Author can update the post
	pOID := utils.GetOID(pid, c)
	if pOID == nil {
		return
	}

	updateFields["updatedAt"] = time.Now()

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

func FindAllPostWithCategory(c *gin.Context) {

	// findOption := options.Find()
	// err := utils.SetupFindOptions(findOption, c)

	// if err != nil {
	// 	return
	// }

	var cOID *primitive.ObjectID

	cid := c.Param("cid")

	if cid == "all" || cid == "" {
		cOID = nil
	} else {
		cOID = utils.GetOID(cid, c)
		if cOID == nil {
			return
		}
	}

	skip, limit, sort, err := utils.GetSkipLimitSortFromContext(c)

	if err != nil {
		return
	}

	sortByLikeCount := false

	if sort != nil {
		if strings.ToLower(*sort) == "likecount" {
			sortByLikeCount = true
		} else {
			sortByLikeCount = false
		}
	}

	posts, err := models.FindAllCategoryPostsWithPreview(cOID, *skip, *limit, sortByLikeCount)

	// posts, err := models.FindPosts(bson.M{}, findOption)
	if err != nil {
		errStr := fmt.Sprintf("Cannot retreive the posts: %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": errStr,
		})
		return
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
	post, err := models.FindSinglePostWithDetail(*pOID)
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

	posts, err := models.FindPostsWithPreview(&[]bson.M{bson.M{"$match": bson.M{"author": aOID}}}, -1, -1, false)
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

	// savedOIDs := []*primitive.ObjectID{}
	// for _, savedId := range user.SavedPosts {
	// 	oid := utils.GetOID(savedId, c)
	// 	if oid == nil {
	// 		return
	// 	}

	// 	savedOIDs = append(savedOIDs, oid)
	// }

	posts, err := models.FindPostsWithPreview(&[]bson.M{bson.M{"$match": bson.M{"_id": bson.M{"$in": user.SavedPosts}}}}, -1, -1, false)

	// posts, err := models.FindPosts(bson.M{"_id": bson.M{"$in": savedOIDs}}, findOption)

	if err != nil {
		errStr := fmt.Sprintf("Cannot find the SavedPosts: %+v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"err": errStr,
		})
		return
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
	posts, err := models.FindPostsWithPreview(&[]bson.M{bson.M{"_id": bson.M{"$in": postsOID}}}, -1, -1, false)
	// posts, err := models.FindPosts(}, findOption)

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

func LikePost(c *gin.Context) {
	pid := c.Param("pid")
	pOID := utils.GetOID(pid, c)

	con := c.Query("condition")

	var like bool

	if con == "1" {
		like = true
	} else {
		like = false
	}

	if pOID == nil {
		return
	}

	user := utils.GetUserFromContext(c)

	if user == nil {
		return
	}

	result, err := models.ToggleLikerForPost(*pOID, user.ID, like)

	if err != nil {
		errStr := fmt.Sprint("Cannnot toggle the post like: %+v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"err": errStr,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": result,
		"uid":    user.ID.Hex(),
		"pid":    pid,
	})

}
