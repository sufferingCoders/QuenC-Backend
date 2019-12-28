package models

import (
	"context"
	"quenc/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

)

type PostCategory struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	CategoryName primitive.ObjectID `json:"categoryName" bson:"categoryName"`
}

// AddPostCategory - Adding PostCategory to MongoDB
func AddPostCategory(inputPostCategory *PostCategory) (interface{}, error) {

	result, err := database.PostCategoryCollection.InsertOne(context.TODO(), inputPostCategory)

	return result.InsertedID, err
}

/// Not YET

// UpdatePostCategorys - Update PostCategory in MongoDB
func UpdatePostCategorys(filterDetail bson.M, updateDetail bson.M) (*mongo.UpdateResult, error) {

	result, err := database.PostCategoryCollection.UpdateMany(context.TODO(), filterDetail, bson.M{"$set": updateDetail})

	return result, err
}

// UpdatePostCategoryByOID - Update PostCategory in MongoDB by its OID
func UpdatePostCategoryByOID(oid primitive.ObjectID, updateDetail bson.M) (*mongo.UpdateResult, error) {

	result, err := database.PostCategoryCollection.UpdateOne(context.TODO(), bson.M{"_id": oid}, bson.M{"$set": updateDetail})

	return result, err
}

// DeletePostCategoryByOID - Delete PostCategory by its OID
func DeletePostCategoryByOID(oid primitive.ObjectID) error {
	_, err := database.PostCategoryCollection.DeleteOne(context.TODO(), bson.M{"_id": oid})
	return err
}

// FindPostCategoryByOID - Find PostCategory by its OID
func FindPostCategoryByOID(oid primitive.ObjectID) (*PostCategory, error) {
	var PostCategory PostCategory

	err := database.PostCategoryCollection.FindOne(context.TODO(), bson.M{"_id": oid}).Decode(&PostCategory)

	return &PostCategory, err
}

// FindPostCategorys - Find Multiple PostCategorys by filterDetail
func FindPostCategorys(filterDetail bson.M, findOptions *options.FindOptions) ([]*PostCategory, error) {
	var postCategorys []*PostCategory
	result, err := database.PostCategoryCollection.Find(context.TODO(), filterDetail, findOptions)
	if result != nil {
		defer result.Close(context.TODO())
	}

	if err != nil {
		return nil, err
	}

	for result.Next(context.TODO()) {
		var elem PostCategory
		err := result.Decode(&elem)
		if err != nil {
			return nil, err
		}
		postCategorys = append(postCategorys, &elem)
	}

	return postCategorys, nil
}

func FindAllPostCategorys(findOptions *options.FindOptions) ([]*PostCategory, error) {
	postCategorys, err := FindPostCategorys(bson.M{}, findOptions)
	return postCategorys, err
}
