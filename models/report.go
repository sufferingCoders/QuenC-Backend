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
	Content      string             `json:"content" bson:"content"`
	AuthorDomain string             `json:"authorDomain" bson:"authorDomain"`
	PreviewText  string             `json:"previewText" bson:"previewText"`
	PreviewPhoto string             `json:"previewPhoto" bson:"previewPhoto"`
	AuthorGender int                `json:"authorGender" bson:"authorGender"`
	ReportTarget int                `json:"reportTarget" bson:"reportTarget"`
	ReportType   int                `json:"reportType" bson:"reportType"`
	Solve        bool               `json:"solve" bson:"solve"`
	Author       primitive.ObjectID `json:"author" bson:"author"`
	ReportID     primitive.ObjectID `json:"reportId" bson:"reportId"`
	CreatedAt    primitive.DateTime `json:"createdAt" bson:"createdAt"`
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
	ReportCollection = database.DB.Collection("Report")
)

func AddReport(inputReport *Report) (interface{}, error) {

	result, err := ReportCollection.InsertOne(context.TODO(), inputReport)

	return result.InsertedID, err
}

// UpdateReports - Update Report in MongoDB
func UpdateReports(filterDetail bson.M, updateDetail bson.M) (*mongo.UpdateResult, error) {

	result, err := ReportCollection.UpdateMany(context.TODO(), filterDetail, bson.M{"$set": updateDetail})

	return result, err
}

// UpdateReportByOID - Update Report in MongoDB by its OID
func UpdateReportByOID(oid primitive.ObjectID, updateDetail bson.M) (*mongo.UpdateResult, error) {

	result, err := ReportCollection.UpdateOne(context.TODO(), bson.M{"_id": oid}, bson.M{"$set": updateDetail})

	return result, err
}

// DeleteReportByOID - Delete Report by its OID
func DeleteReportByOID(oid primitive.ObjectID) error {
	_, err := ReportCollection.DeleteOne(context.TODO(), bson.M{"_id": oid})
	return err
}

// FindReportByOID - Find Report by its OID
func FindReportByOID(oid primitive.ObjectID) (*Report, error) {
	var report Report

	err := ReportCollection.FindOne(context.TODO(), bson.M{"_id": oid}).Decode(&report)

	return &report, err
}

// FindReports - Find Multiple Reports by filterDetail
func FindReports(filterDetail bson.M, findOptions *options.FindOptions) ([]*Report, error) {
	var reports []*Report
	result, err := ReportCollection.Find(context.TODO(), filterDetail, findOptions)
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
