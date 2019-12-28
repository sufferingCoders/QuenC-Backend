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

// Post -Post Schema
type Post struct {
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

// AddPost - Adding Post to MongoDB
func AddPost(inputPost *Post) (interface{}, error) {

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
func FindPostByOID(oid primitive.ObjectID) (*Post, error) {
	var post Post

	err := database.PostCollection.FindOne(context.TODO(), bson.M{"_id": oid}).Decode(&post)

	return &post, err
}

// FindPosts - Find Multiple Posts by filterDetail
func FindPosts(filterDetail bson.M, findOptions *options.FindOptions) ([]*Post, error) {
	var posts []*Post
	result, err := database.PostCollection.Find(context.TODO(), filterDetail, findOptions)
	if result != nil {
		defer result.Close(context.TODO())
	}

	if err != nil {
		return nil, err
	}

	for result.Next(context.TODO()) {
		var elem Post
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
	if like {
		pullOrPush = "$push"
	} else {
		pullOrPush = "$pull"
	}

	result, err := database.PostCategoryCollection.UpdateOne(
		context.TODO(),
		bson.M{"_id": pOID},
		bson.M{pullOrPush: bson.M{"likers": uOID}},
	)
	return result, err
}

// Giving pipeline here
func FindPostsPreview() ([]*Post, error) {
	var posts []*Post

	pipeline := []bson.M{
		bson.M{"$sort": bson.M{"createdAt": -1}},
		bson.M{"$limit": 30},
		bson.M{"$lookup": bson.M{
			"from":         "User",
			"localField":   "author",
			"foreignField": "_id",
			"as":           "authorObj",
		},
		},
		bson.M{
			"$project": bson.M{
				"authorObj":    bson.M{"$arrayElemAt": bson.A{"$authorObj", 0}},
				"_id":          1,
				"previewText":  1,
				"previewPhoto": 1,
				"title":        1,
				"category":     1,
			},
		},
	}

	result, err := database.PostCategoryCollection.Aggregate(context.TODO(), pipeline)
	if result != nil {
		defer result.Close(context.TODO())
	}

	result.All(context.TODO(), &posts)

	if err != nil {
		return nil, err
	}

	// for result.Next(context.TODO()) {
	// 	var elem Post
	// 	err := result.Decode(&elem)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	posts = append(posts, &elem)
	// }

	return posts, nil

}

// FindPostByAuthor - find posts for certain author
func FindPostByAuthor(uOID primitive.ObjectID, findOptions *options.FindOptions) ([]*Post, error) {
	posts, err := FindPosts(bson.M{"author": uOID}, findOptions)
	return posts, err
}

func FindAllCategoryPostsWithPreview(cOID *primitive.ObjectID, skip int, limit int, sortByLikeCount bool) ([]*Post, error) {
	cond := []bson.M{}
	if cOID != nil {
		cond = append(cond, bson.M{"$match": bson.M{"category": cOID}})
	} else {
		cond = nil
	}

	posts, err := FindPostsWithPreview(&cond, skip, limit, sortByLikeCount)
	return posts, err
}

func FindPostsWithPreview(matchingCond *[]bson.M, skip int, limit int, sortByLikeCount bool) ([]*Post, error) {
	var posts []*Post

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
					bson.M{"$project": bson.M{"categoryName1": 1, "_id": 1}},
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
				"previewText":  1,
				"title":        1,
				"previewPhoto": 1,
				"createdAt":    1,
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
				"$sort": -1,
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

func FindSinglePostWithDetail(pOID primitive.ObjectID) (*Post, error) {

	posts, err := FindPostWithDetail(&[]bson.M{
		bson.M{"$match": bson.M{
			"_id": pOID,
		}},
	})

	return posts[0], err
}

func FindPostWithDetail(matchingCond *[]bson.M) ([]*Post, error) {
	var posts []*Post

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
					bson.M{"$project": bson.M{"categoryName1": 1, "_id": 1}},
				},
				"as": "category",
			},
		},
		// Project
		bson.M{
			"$project": bson.M{
				"_id":       1,
				"likeCount": bson.M{"$size": "$likers"},
				"author":    bson.M{"$arrayElemAt": bson.A{"$author", 0}},
				"category":  bson.M{"$arrayElemAt": bson.A{"$category", 0}},
				"title":     1,
				"content":   1,
				"createdAt": 1,
				"updatedAt": 1,
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
