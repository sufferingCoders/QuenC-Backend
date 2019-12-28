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

type Comment struct {
	ID         primitive.ObjectID   `json:"_id" bson:"_id,omitempty"`
	BelongPost primitive.ObjectID   `json:"belongPost" bson:"belongPost"`
	Author     primitive.ObjectID   `json:"author" bson:"author"`
	Content    string               `json:"content" bson:"content"`
	Likers     []primitive.ObjectID `json:"likers" bson:"likers"`
	UpdatedAt  time.Time            `json:"updatedAt" bson:"updatedAt"`
	CreatedAt  time.Time            `json:"createdAt" bson:"createdAt"`
}

// AddComment - Adding Comment to MongoDB
func AddComment(inputComment *Comment) (interface{}, error) {

	result, err := database.CommentCollection.InsertOne(context.TODO(), inputComment)

	return result.InsertedID, err
}

// UpdateComments - Update Comment in MongoDB
func UpdateComments(filterDetail bson.M, updateDetail bson.M) (*mongo.UpdateResult, error) {

	result, err := database.CommentCollection.UpdateMany(context.TODO(), filterDetail, bson.M{"$set": updateDetail})

	return result, err
}

// UpdateCommentByOID - Update Comment in MongoDB by its OID
func UpdateCommentByOID(oid primitive.ObjectID, updateDetail bson.M) (*mongo.UpdateResult, error) {

	result, err := database.CommentCollection.UpdateOne(context.TODO(), bson.M{"_id": oid}, bson.M{"$set": updateDetail})

	return result, err
}

// DeleteCommentByOID - Delete Comment by its OID
func DeleteCommentByOID(oid primitive.ObjectID) error {
	_, err := database.CommentCollection.DeleteOne(context.TODO(), bson.M{"_id": oid})
	return err
}

// FindCommentByOID - Find Comment by its OID
func FindCommentByOID(oid primitive.ObjectID) (*Comment, error) {
	var comment Comment

	err := database.CommentCollection.FindOne(context.TODO(), bson.M{"_id": oid}).Decode(&comment)

	return &comment, err
}

// FindComments - Find Multiple Comments by filterDetail
func FindComments(filterDetail bson.M, findOptions *options.FindOptions) ([]*Comment, error) {
	var comments []*Comment
	result, err := database.CommentCollection.Find(context.TODO(), filterDetail, findOptions)
	if result != nil {
		defer result.Close(context.TODO())
	}

	if err != nil {
		return nil, err
	}

	for result.Next(context.TODO()) {
		var elem Comment
		err := result.Decode(&elem)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &elem)
	}

	return comments, nil
}

// FindCommentByPost - find Comments for certain post
func FindCommentByPost(uOID primitive.ObjectID, findOptions *options.FindOptions) ([]*Comment, error) {
	comments, err := FindComments(bson.M{"author": uOID}, findOptions)
	return comments, err
}

func ToggleLikerForComment(cOID primitive.ObjectID, uOID primitive.ObjectID, like bool) (*mongo.UpdateResult, error) {
	var pullOrPush string
	if like {
		pullOrPush = "$push"
	} else {
		pullOrPush = "$pull"
	}

	result, err := database.CommentCollection.UpdateOne(
		context.TODO(),
		bson.M{"_id": cOID},
		bson.M{pullOrPush: bson.M{"likers": uOID}},
	)
	return result, err
}

func FindCommentsWithDetailForPost(pOID primitive.ObjectID, skip int, limit int, sortByLikeCount bool) ([]*Comment, error) {
	var comments []*Comment
	pipeline := []bson.M{
		bson.M{"$match": bson.M{
			"belongPost": pOID,
		}},

		// Populate Author
		bson.M{
			"$lookup": bson.M{
				"from": "user",
				"let":  bson.M{"author": "$author"},
				"pipeline": bson.A{
					bson.M{"$match": bson.M{"$expr": bson.M{"$eq": bson.A{"$_id", "$$author"}}}},
					bson.M{"$project": bson.M{"_id": 1, "gender": 1, "domain": 1}},
				},
				"as": "author",
			},
		},

		// Project needed fields
		// Get likeCount,
		bson.M{
			"$project": bson.M{
				"_id":        1,
				"belongPost": 1,
				"likeCount":  bson.M{"$size": "$likers"},
				"author":     bson.M{"$arrayElemAt": bson.A{"$author", 0}},
				"content":    1,
				"createdAt":  1,
			},
		},

		// Sort by created at
		bson.M{
			"$sort": bson.M{"createdAt": -1},
		},
	}

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

	result, err := database.CommentCollection.Aggregate(context.TODO(), pipeline)

	if result != nil {
		defer result.Close(context.TODO())
	}

	if err != nil {
		return nil, err
	}

	err = result.All(context.TODO(), &comments)

	if err != nil {
		return nil, err
	}

	return comments, nil

}

func GetSingleCommentWithDetail(cOID primitive.ObjectID) (*Comment, error) {
	var comments []*Comment
	pipeline := []bson.M{
		bson.M{"$match": bson.M{
			"_id": cOID,
		}},

		// Populate Author
		bson.M{
			"$lookup": bson.M{
				"from": "user",
				"let":  bson.M{"author": "$author"},
				"pipeline": bson.A{
					bson.M{"$match": bson.M{"$expr": bson.M{"$eq": bson.A{"$_id", "$$author"}}}},
					bson.M{"$project": bson.M{"_id": 1, "gender": 1, "domain": 1}},
				},
				"as": "author",
			},
		},

		// Project needed fields
		// Get likeCount,
		bson.M{
			"$project": bson.M{
				"_id":        1,
				"belongPost": 1,
				"likeCount":  bson.M{"$size": "$likers"},
				"author":     bson.M{"$arrayElemAt": bson.A{"$author", 0}},
				"content":    1,
				"createdAt":  1,
				"updatedAt":  1,
			},
		},
	}

	result, err := database.CommentCollection.Aggregate(context.TODO(), pipeline)
	if result != nil {
		defer result.Close(context.TODO())
	}

	if err != nil {
		return nil, err
	}

	err = result.All(context.TODO(), &comments)

	if err != nil {
		return nil, err
	}

	return comments[0], nil

}
