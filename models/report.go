package models

import (
	"context"
	"quenc/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Report struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Content      primitive.ObjectID `json:"content" bson:"content"`
	Author       primitive.ObjectID `json:"author" bson:"author"`
	AuthorDomain primitive.ObjectID `json:"authorDomain" bson:"authorDomain"`
	AuthorGender primitive.ObjectID `json:"authorGender" bson:"authorGender"`
	PreviewText  primitive.ObjectID `json:"previewText" bson:"previewText"`
	PreviewPhoto primitive.ObjectID `json:"previewPhoto" bson:"previewPhoto"`
	ReportTarget primitive.ObjectID `json:"reportTarget" bson:"reportTarget"`
	ReportType   primitive.ObjectID `json:"reportType" bson:"reportType"`
	CreatedAt    primitive.ObjectID `json:"createdAt" bson:"createdAt"`
	ReportId     primitive.ObjectID `json:"reportId" bson:"reportId"`
	Solve        primitive.ObjectID `json:"solve" bson:"solve"`
}

var (
	reportTypeCodeList = []string{
		"其他",             // 0
		"謾罵他人",           // 1
		"惡意洗版",           // 2
		"惡意洩漏他人資料",       // 3
		"包含色情, 血腥, 暴力內容", // 4
		"廣告和宣傳內容",        // 5
	}
	reportCollection = database.DB.Collection("Report")
)

func AddReport(inputReport *Report) (interface{}, error) {

	result, err := reportCollection.InsertOne(context.TODO(), inputReport)

	return result.InsertedID, err
}

// UpdateReports - Update Report in MongoDB
func UpdateReports(filterDetail bson.M, updateDetail bson.M) (*mongo.UpdateResult, error) {

	result, err := reportCollection.UpdateMany(context.TODO(), filterDetail, bson.M{"$set": updateDetail})

	return result, err
}

// UpdateReportByOID - Update Report in MongoDB by its OID
func UpdateReportByOID(oid primitive.ObjectID, updateDetail bson.M) (*mongo.UpdateResult, error) {

	result, err := reportCollection.UpdateOne(context.TODO(), bson.M{"_id": oid}, bson.M{"$set": updateDetail})

	return result, err
}

// DeleteReportByOID - Delete Report by its OID
func DeleteReportByOID(oid primitive.ObjectID) error {
	_, err := reportCollection.DeleteOne(context.TODO(), bson.M{"_id": oid})
	return err
}

// FindReportByOID - Find Report by its OID
func FindReportByOID(oid primitive.ObjectID) (*Report, error) {
	var report Report

	err := reportCollection.FindOne(context.TODO(), bson.M{"_id": oid}).Decode(&report)

	return &report, err
}

// FindReports - Find Multiple Reports by filterDetail
func FindReports(filterDetail bson.M, findOptions *options.FindOptions) ([]*Report, error) {
	var reports []*Report
	result, err := reportCollection.Find(context.TODO(), filterDetail, findOptions)
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
