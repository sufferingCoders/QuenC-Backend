package models

import (
	"context"
	"quenc/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

)

type Comment struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	BelongPost   string             `json:"belongPost" bson:"belongPost"`
	Author       string             `json:"author" bson:"author"`
	AuthorDomain string             `json:"authorDomain" bson:"authorDomain"`
	Content      string             `json:"content" bson:"content"`
	AuthorGender int                `json:"authorGender" bson:"authorGender"`
	LikeCount    int                `json:"likeCount" bson:"likeCount"`
	UpdatedAt    primitive.DateTime `json:"updatedAt" bson:"updatedAt"`
	CreatedAt    primitive.DateTime `json:"createdAt" bson:"createdAt"`
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
	defer result.Close(context.TODO())

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


