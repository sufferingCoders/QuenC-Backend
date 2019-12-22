package models

import (
	"context"
	"quenc/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

)

// Post -Post Schema
type Post struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	AuthorGender int                `json:"authorGender" bson:"authorGender"`
	LikeCount    int                `json:"likeCount" bson:"likeCount"`
	Anonymous    bool               `json:"anonymous" bson:"anonymous"`
	Title        string             `json:"title" bson:"title"`
	Author       string             `json:"author" bson:"author"`
	AuthorDomain string             `json:"authorDomain" bson:"authorDomain"`
	Content      string             `json:"content" bson:"content"`
	PreviewText  string             `json:"previewText" bson:"previewText"`
	PreviewPhoto string             `json:"previewPhoto" bson:"previewPhoto"`
	Category     string             `json:"category" bson:"category"`
	UpdatedAt    primitive.DateTime `json:"updatedAt" bson:"updatedAt"`
	CreatedAt    primitive.DateTime `json:"createdAt" bson:"createdAt"`
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
	defer result.Close(context.TODO())

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

// FindPostByAuthor - find posts for certain author
func FindPostByAuthor(uOID primitive.ObjectID, findOptions *options.FindOptions) ([]*Post, error) {
	posts, err := FindPosts(bson.M{"author": uOID}, findOptions)
	return posts, err
}
