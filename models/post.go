package models

import (
	"context"
	"quenc/database"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

)

// PostAdding -PostAdding Schema
type PostAdding struct {
	ID           primitive.ObjectID   `json:"_id" bson:"_id,omitempty"`
	Anonymous    bool                 `json:"anonymous" bson:"anonymous"`
	Title        string               `json:"title" bson:"title"`
	Author       primitive.ObjectID   `json:"author" bson:"author"`
	Content      string               `json:"content" bson:"content"`
	PreviewText  string               `json:"previewText" bson:"previewText"`
	PreviewPhoto string               `json:"previewPhoto" bson:"previewPhoto"`
	Category     primitive.ObjectID   `json:"category" bson:"category"`
	Likers       []primitive.ObjectID `json:"likers" bson:"likers"`
	UpdatedAt    time.Time            `json:"updatedAt" bson:"updatedAt"`
	CreatedAt    time.Time            `json:"createdAt" bson:"createdAt"`
}

type PostPreview struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Anonymous bool               `json:"anonymous" bson:"anonymous"`

	Title        string       `json:"title" bson:"title"`
	Author       User         `json:"author" bson:"author"`
	Content      string       `json:"content" bson:"content"`
	PreviewText  string       `json:"previewText" bson:"previewText"`
	PreviewPhoto string       `json:"previewPhoto" bson:"previewPhoto"`
	Category     PostCategory `json:"category" bson:"category"`
	UpdatedAt    time.Time    `json:"updatedAt" bson:"updatedAt"`
	CreatedAt    time.Time    `json:"createdAt" bson:"createdAt"`
	LikeCount    int          `json:"likeCount" bson:"likeCount"`
}

type PostDetail struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Anonymous    bool               `json:"anonymous" bson:"anonymous"`
	Title        string             `json:"title" bson:"title"`
	Author       User               `json:"author" bson:"author"`
	Content      string             `json:"content" bson:"content"`
	Category     PostCategory       `json:"category" bson:"category"`
	UpdatedAt    time.Time          `json:"updatedAt" bson:"updatedAt"`
	CreatedAt    time.Time          `json:"createdAt" bson:"createdAt"`
	LikeCount    int                `json:"likeCount" bson:"likeCount"`
	PreviewText  string             `json:"previewText" bson:"previewText"`
	PreviewPhoto string             `json:"previewPhoto" bson:"previewPhoto"`
}

// AddPost - Adding Post to MongoDB
func AddPost(inputPost *PostAdding) (interface{}, error) {

	result, err := database.PostCollection.InsertOne(context.TODO(), inputPost)

	return result.InsertedID, err
}

// UpdatePosts - Update Post in MongoDB
func UpdatePosts(filterDetail bson.M, updateDetail bson.M) (*mongo.UpdateResult, error) {

	result, err := database.PostCollection.UpdateMany(context.TODO(), filterDetail, bson.M{"$set": updateDetail})

	return result, err
}

// UpdatePostByOID - Update Post in MongoDB by its OID
func UpdatePostByOID(oid primitive.ObjectID, updateDetail bson.M) (*mongo.UpdateResult, error) {

	result, err := database.PostCollection.UpdateOne(context.TODO(), bson.M{"_id": oid}, bson.M{"$set": updateDetail})

	return result, err
}

// DeletePostByOID - Delete Post by its OID
func DeletePostByOID(oid primitive.ObjectID) error {
	_, err := database.PostCollection.DeleteOne(context.TODO(), bson.M{"_id": oid})
	return err
}

// FindPostByOID - Find Post by its OID
func FindPostByOID(oid primitive.ObjectID) (*PostAdding, error) {
	var post PostAdding

	err := database.PostCollection.FindOne(context.TODO(), bson.M{"_id": oid}).Decode(&post)

	return &post, err
}

// FindPosts - Find Multiple Posts by filterDetail
func FindPosts(filterDetail bson.M, findOptions *options.FindOptions) ([]*PostAdding, error) {
	var posts []*PostAdding
	result, err := database.PostCollection.Find(context.TODO(), filterDetail, findOptions)
	if result != nil {
		defer result.Close(context.TODO())
	}

	if err != nil {
		return nil, err
	}

	for result.Next(context.TODO()) {
		var elem PostAdding
		err := result.Decode(&elem)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &elem)
	}

	return posts, nil
}

func ToggleLikerForPost(pOID primitive.ObjectID, uOID primitive.ObjectID, like bool) (*mongo.UpdateResult, error) {
	var pullOrPush string
	var reverse string

	if like {
		pullOrPush = "$push"
		reverse = "$pull"

	} else {
		pullOrPush = "$pull"
		reverse = "$push"

	}

	_, err := database.UserCollection.UpdateOne(
		context.TODO(),
		bson.M{"_id": uOID},
		bson.M{pullOrPush: bson.M{"likePosts": pOID}},
	)

	if err != nil {
		return nil, err
	}

	result, err := database.PostCollection.UpdateOne(
		context.TODO(),
		bson.M{"_id": pOID},
		bson.M{pullOrPush: bson.M{"likers": uOID}},
	)

	if err != nil {
		// Reverse the Post one
		_, err := database.UserCollection.UpdateOne(
			context.TODO(),
			bson.M{"_id": uOID},
			bson.M{reverse: bson.M{"likePosts": pOID}},
		)

		return nil, err
	}

	return result, err
}

// // Giving pipeline here
// func FindPostsPreview() ([]*PostPreview, error) {
// 	var posts []*PostPreview

// 	pipeline := []bson.M{
// 		bson.M{"$sort": bson.M{"createdAt": -1}},
// 		bson.M{"$limit": 30},
// 		bson.M{"$lookup": bson.M{
// 			"from":         "User",
// 			"localField":   "author",
// 			"foreignField": "_id",
// 			"as":           "authorObj",
// 		},
// 		},
// 		bson.M{
// 			"$project": bson.M{
// 				"authorObj":    bson.M{"$arrayElemAt": bson.A{"$authorObj", 0}},
// 				"_id":          1,
// 				"previewText":  1,
// 				"previewPhoto": 1,
// 				"title":        1,
// 				"category":     1,
// 			},
// 		},
// 	}

// 	result, err := database.PostCategoryCollection.Aggregate(context.TODO(), pipeline)
// 	if result != nil {
// 		defer result.Close(context.TODO())
// 	}

// 	result.All(context.TODO(), &posts)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return posts, nil

// }

// FindPostByAuthor - find posts for certain author
func FindPostByAuthor(uOID primitive.ObjectID, findOptions *options.FindOptions) ([]*PostAdding, error) {
	posts, err := FindPosts(bson.M{"author": uOID}, findOptions)
	return posts, err
}

func FindAllCategoryPostsWithPreview(cOID *primitive.ObjectID, skip int, limit int, sortByLikeCount bool) ([]*PostPreview, error) {
	cond := []bson.M{}
	if cOID != nil {
		cond = append(cond, bson.M{"$match": bson.M{"category": cOID}})
	} else {
		cond = nil
	}

	posts, err := FindPostsWithPreview(&cond, skip, limit, sortByLikeCount)
	return posts, err
}

func FindPostsWithPreview(matchingCond *[]bson.M, skip int, limit int, sortByLikeCount bool) ([]*PostPreview, error) {
	var posts []*PostPreview

	var pipeline = []bson.M{}

	if matchingCond != nil {
		// Find match
		pipeline = append(pipeline, *matchingCond...)
	}

	pipeline = append(pipeline, []bson.M{
		// Populate Author
		bson.M{
			"$lookup": bson.M{
				"from": "user",
				"let":  bson.M{"author": "$author", "anonymous": "$anonymous"},
				"pipeline": bson.A{
					bson.M{"$match": bson.M{"$expr": bson.M{"$eq": bson.A{"$_id", "$$author"}}}},
					bson.M{"$project": bson.M{"_id": 1, "gender": 1, "domain": bson.M{"$cond": bson.M{"if": bson.M{"$eq": bson.A{"$$anonymous", true}}, "then": "", "else": "$domain"}}}},
				},
				"as": "author",
			},
		},
		// Populate Category
		bson.M{
			"$lookup": bson.M{
				"from": "postCategory",
				"let":  bson.M{"category": "$category"},
				"pipeline": bson.A{
					bson.M{"$match": bson.M{"$expr": bson.M{"$eq": bson.A{"$_id", "$$category"}}}},
					bson.M{"$project": bson.M{"categoryName": 1, "_id": 1}},
				},
				"as": "category",
			},
		},
		// Project
		bson.M{
			"$project": bson.M{
				"_id":          1,
				"likeCount":    bson.M{"$size": "$likers"},
				"author":       bson.M{"$arrayElemAt": bson.A{"$author", 0}},
				"category":     bson.M{"$arrayElemAt": bson.A{"$category", 0}},
				"title":        1,
				"previewText":  1,
				"previewPhoto": 1,
				"createdAt":    1,
				"anonymous":    1,
			},
		},
		// Sorting
		bson.M{
			"$sort": bson.M{"createdAt": -1},
		},
	}...)

	if sortByLikeCount {
		pipeline = append(pipeline, bson.M{
			"$sort": bson.M{
				"likeCount": -1,
			},
		})
	}

	if skip > 0 {
		pipeline = append(pipeline, bson.M{
			"$skip": skip,
		})
	}

	if limit > 0 {
		pipeline = append(pipeline, bson.M{
			"$limit": limit,
		})
	}

	result, err := database.PostCollection.Aggregate(context.TODO(), pipeline)
	if result != nil {
		defer result.Close(context.TODO())
	}
	if err != nil {
		return nil, err
	}

	err = result.All(context.TODO(), &posts)

	if err != nil {
		return nil, err
	}

	return posts, nil

}

func FindSinglePostWithDetail(pOID primitive.ObjectID) (*PostDetail, error) {

	posts, err := FindPostWithDetail(&[]bson.M{
		bson.M{"$match": bson.M{
			"_id": pOID,
		}},
	})

	return posts[0], err
}

func FindPostWithDetail(matchingCond *[]bson.M) ([]*PostDetail, error) {
	var posts []*PostDetail

	var pipeline = []bson.M{}

	if matchingCond != nil {
		// Find match
		pipeline = append(pipeline, *matchingCond...)
	}

	pipeline = append(pipeline, []bson.M{

		// Populate Author
		bson.M{
			"$lookup": bson.M{
				"from": "user",
				"let":  bson.M{"author": "$author", "anonymous": "$anonymous"},
				"pipeline": bson.A{
					bson.M{"$match": bson.M{"$expr": bson.M{"$eq": bson.A{"$_id", "$$author"}}}},
					bson.M{"$project": bson.M{"_id": 1, "gender": 1, "domain": bson.M{"$cond": bson.M{"if": bson.M{"$eq": bson.A{"$$anonymous", true}}, "then": "", "else": "$domain"}}}},
				},
				"as": "author",
			},
		},
		// Populate Category
		bson.M{
			"$lookup": bson.M{
				"from": "postCategory",
				"let":  bson.M{"category": "$category"},
				"pipeline": bson.A{
					bson.M{"$match": bson.M{"$expr": bson.M{"$eq": bson.A{"$_id", "$$category"}}}},
					bson.M{"$project": bson.M{"categoryName": 1, "_id": 1}},
				},
				"as": "category",
			},
		},
		// Project
		bson.M{
			"$project": bson.M{
				"_id":          1,
				"likeCount":    bson.M{"$size": "$likers"},
				"author":       bson.M{"$arrayElemAt": bson.A{"$author", 0}},
				"category":     bson.M{"$arrayElemAt": bson.A{"$category", 0}},
				"title":        1,
				"content":      1,
				"createdAt":    1,
				"updatedAt":    1,
				"anonymous":    1,
				"previewText":  1,
				"previewPhoto": 1,
			},
		},
		// Sorting

	}...)

	result, err := database.PostCollection.Aggregate(context.TODO(), pipeline)
	if result != nil {
		defer result.Close(context.TODO())
	}

	if err != nil {
		return nil, err
	}

	err = result.All(context.TODO(), &posts)

	if err != nil {
		return nil, err
	}

	return posts, nil

}
