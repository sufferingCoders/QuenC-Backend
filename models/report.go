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

// 回傳Report時, 需要Populate什麼
// Author
// ReportID

type ReportAdding struct {
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

type ReportPreview struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	PreviewText  string             `json:"previewText" bson:"previewText"`
	PreviewPhoto string             `json:"previewPhoto" bson:"previewPhoto"`
	ReportTarget int                `json:"reportTarget" bson:"reportTarget"`
	Solve        bool               `json:"solve" bson:"solve"`
	Author       User               `json:"author" bson:"author"`
	CreatedAt    time.Time          `json:"createdAt" bson:"createdAt"`
}
type ReportDetail struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Content      string             `json:"content" bson:"content"`
	ReportTarget int                `json:"reportTarget" bson:"reportTarget"`
	Solve        bool               `json:"solve" bson:"solve"`
	Author       User               `json:"author" bson:"author"`
	ReportID     primitive.ObjectID `json:"reportId" bson:"reportId"`
	CreatedAt    time.Time          `json:"createdAt" bson:"createdAt"`
}

func AddReport(inputReport *ReportAdding) (interface{}, error) {

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
func FindReportByOID(oid primitive.ObjectID) (*ReportAdding, error) {
	var report ReportAdding

	err := database.ReportCollection.FindOne(context.TODO(), bson.M{"_id": oid}).Decode(&report)

	return &report, err
}

// FindReports - Find Multiple Reports by filterDetail
func FindReports(filterDetail bson.M, findOptions *options.FindOptions) ([]*ReportAdding, error) {
	var reports []*ReportAdding
	result, err := database.ReportCollection.Find(context.TODO(), filterDetail, findOptions)
	if result != nil {
		defer result.Close(context.TODO())
	}

	if err != nil {
		return nil, err
	}

	for result.Next(context.TODO()) {
		var elem ReportAdding
		err := result.Decode(&elem)
		if err != nil {
			return nil, err
		}
		reports = append(reports, &elem)
	}

	return reports, nil
}

func FindAllReporstWithPreview(skip int, limit int) ([]*ReportPreview, error) {
	reports, err := FindReportsWithPreview(nil, skip, limit)
	return reports, err
}

func FindReportsWithPreview(matchingCond *[]bson.M, skip int, limit int) ([]*ReportPreview, error) {
	// This will return the report sort by createdAt
	var reports []*ReportPreview

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
				"let":  bson.M{"author": "$author"},
				"pipeline": bson.A{
					bson.M{"$match": bson.M{"$expr": bson.M{"$eq": bson.A{"$_id", "$$author"}}}},
					bson.M{"$project": bson.M{"_id": 1, "gender": 1, "domain": bson.M{"$cond": bson.M{"if": bson.M{"$eq": bson.A{"$$anonymous", true}}, "then": "", "else": "$domain"}}}},
				},
				"as": "author",
			},
		},

		bson.M{
			"$project": bson.M{
				"reportTarget": 1,
				"PreviewText":  1,
				"PreviewPhoto": 1,
				"author":       1,
				"solve":        1,
				"createdAt":    1,
				"_id":          1,
			},
		},
		// Sort by createdAt
		// Old one go higher
		bson.M{
			"$sort": bson.M{"createdAt": 1},
		},
	}...)

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

	result, err := database.ReportCollection.Aggregate(context.TODO(), pipeline)
	if result != nil {
		defer result.Close(context.TODO())
	}

	if err != nil {
		return nil, err
	}

	err = result.All(context.TODO(), &reports)

	if err != nil {
		return nil, err
	}

	return reports, nil
}

func FindSingleReportWithDetail(rOID primitive.ObjectID) (*ReportDetail, error) {
	reports, err := FindReportWithDetail(&[]bson.M{bson.M{"$match": bson.M{"_id": rOID}}})
	return reports[0], err
}

func FindReportWithDetail(matchingCond *[]bson.M) ([]*ReportDetail, error) {
	// This will return the report sort by createdAt
	var reports []*ReportDetail

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
				"let":  bson.M{"author": "$author"},
				"pipeline": bson.A{
					bson.M{"$match": bson.M{"$expr": bson.M{"$eq": bson.A{"$_id", "$$author"}}}},
					bson.M{"$project": bson.M{"_id": 1, "gender": 1, "domain": bson.M{"$cond": bson.M{"if": bson.M{"$eq": bson.A{"$$anonymous", true}}, "then": "", "else": "$domain"}}}},
				},
				"as": "author",
			},
		},
		// Sort by createdAt
		// Old one go higher
		bson.M{
			"$project": bson.M{
				"reportTarget": 1,
				"content":      1,
				"author":       1,
				"solve":        1,
				"createdAt":    1,
				"reportId":     1,
				"_id":          1,
			},
		},
	}...)

	result, err := database.ReportCollection.Aggregate(context.TODO(), pipeline)
	if result != nil {
		defer result.Close(context.TODO())
	}

	if err != nil {
		return nil, err
	}

	err = result.All(context.TODO(), &reports)

	if err != nil {
		return nil, err
	}

	return reports, nil
}
