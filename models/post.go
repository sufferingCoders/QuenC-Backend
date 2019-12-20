package models

import (
	"context"
	"quenc/database"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Post -Post Schema
type Post struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Author       string             `json:"author" bson:"author"`
	AuthorDomain string             `json:"authorDomain" bson:"authorDomain"`
	AuthorGender int                `json:"authorGender" bson:"authorGender"`
	Title        string             `json:"title" bson:"title"`
	Content      string             `json:"content" bson:"content"`
	CreatedAt    primitive.DateTime `json:"createdAt" bson:"createdAt"`
	UpdatedAt    primitive.DateTime `json:"updatedAt" bson:"updatedAt"`
	Anonymous    bool               `json:"anonymous" bson:"anonymous"`
	PreviewText  string             `json:"previewText" bson:"previewText"`
	PreviewPhoto string             `json:"previewPhoto" bson:"previewPhoto"`
	LikeCount    int                `json:"likeCount" bson:"likeCount"`
	Category     string             `json:"category" bson:"category"`
}

var (
	postCollection = database.DB.Collection("Post")
)

// AddPost - Adding Post to MongoDB
func AddPost(inputPost *Post) (interface{}, error) {

	result, err := postCollection.InsertOne(context.TODO(), inputPost)
	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil
}
