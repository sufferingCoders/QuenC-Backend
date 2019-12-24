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

type Report struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Content      string             `json:"content" bson:"content"`
	PreviewText  string             `json:"previewText" bson:"previewText"`
	PreviewPhoto string             `json:"previewPhoto" bson:"previewPhoto"`
	ReportTarget int                `json:"reportTarget" bson:"reportTarget"`
	Solve        bool               `json:"solve" bson:"solve"`
	Author       primitive.ObjectID `json:"author" bson:"author"`
	ReportID     primitive.ObjectID `json:"reportId" bson:"reportId"`
	CreatedAt    time.Time          `json:"createdAt" bson:"createdAt"`
}

func AddReport(inputReport *Report) (interface{}, error) {

	result, err := database.ReportCollection.InsertOne(context.TODO(), inputReport)

	return result.InsertedID, err
}

// UpdateReports - Update Report in MongoDB
func UpdateReports(filterDetail bson.M, updateDetail bson.M) (*mongo.UpdateResult, error) {

	result, err := database.ReportCollection.UpdateMany(context.TODO(), filterDetail, bson.M{"$set": updateDetail})

	return result, err
}

// UpdateReportByOID - Update Report in MongoDB by its OID
func UpdateReportByOID(oid primitive.ObjectID, updateDetail bson.M) (*mongo.UpdateResult, error) {

	result, err := database.ReportCollection.UpdateOne(context.TODO(), bson.M{"_id": oid}, bson.M{"$set": updateDetail})

	return result, err
}

// DeleteReportByOID - Delete Report by its OID
func DeleteReportByOID(oid primitive.ObjectID) error {
	_, err := database.ReportCollection.DeleteOne(context.TODO(), bson.M{"_id": oid})
	return err
}

// FindReportByOID - Find Report by its OID
func FindReportByOID(oid primitive.ObjectID) (*Report, error) {
	var report Report

	err := database.ReportCollection.FindOne(context.TODO(), bson.M{"_id": oid}).Decode(&report)

	return &report, err
}

// FindReports - Find Multiple Reports by filterDetail
func FindReports(filterDetail bson.M, findOptions *options.FindOptions) ([]*Report, error) {
	var reports []*Report
	result, err := database.ReportCollection.Find(context.TODO(), filterDetail, findOptions)
	defer result.Close(context.TODO())

	if err != nil {
		return nil, err
	}

	for result.Next(context.TODO()) {
		var elem Report
		err := result.Decode(&elem)
		if err != nil {
			return nil, err
		}
		reports = append(reports, &elem)
	}

	return reports, nil
}
