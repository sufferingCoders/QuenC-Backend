package models

import (
	"context"
	"errors"
	"log"
	"quenc/database"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

)

type ChatRoomAdding struct {
	ID        primitive.ObjectID   `json:"_id" bson:"_id,omitempty"`
	Members   []primitive.ObjectID `json:"members" bson:"members"`
	Messages  []Message            `json:"messages" bson:"messages"`
	CreatedAt time.Time            `json:"createdAt" bson:"createdAt"`
	IsGroup   bool                 `json:"isGroup" bson:"isGroup"`
	GroupName string               `json:"groupName" bson:"groupName"`
}

// After Populating
type ChatRoomDetail struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Members   []User             `json:"members" bson:"members"`
	Messages  []Message          `json:"messages" bson:"messages"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	IsGroup   bool               `json:"isGroup" bson:"isGroup"`
	GroupName string             `json:"groupName" bson:"groupName"`
}

// GroupChatRoom will generate a ID and the normal chatRoom will use the members' id

func AddChatRoom(inputChatRoom *ChatRoomAdding) (interface{}, error) {
	// Generate the id here

	if !inputChatRoom.IsGroup {
		if len(inputChatRoom.Members) != 2 {
			return nil, errors.New("the number of memebers in this chatroom is not acceptable")
		}

		if inputChatRoom.Members[0].Hex() == inputChatRoom.Members[1].Hex() {
			return nil, errors.New("Cannot set two same members in the same room")

		} else if inputChatRoom.Members[0].Hex() > inputChatRoom.Members[1].Hex() {
			combineId, err := primitive.ObjectIDFromHex(inputChatRoom.Members[0].Hex() + "-" + inputChatRoom.Members[1].Hex())
			if err != nil {
				return nil, err
			}
			inputChatRoom.ID = combineId
		}
	}

	inputChatRoom.CreatedAt = time.Now()
	inputChatRoom.Messages = []Message{}

	result, err := database.ChatRoomCollection.InsertOne(context.TODO(), inputChatRoom)

	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil
}

func UpdateChatRooms(filterDetail bson.M, updateDetail bson.M) (*mongo.UpdateResult, error) {

	result, err := database.ChatRoomCollection.UpdateMany(context.TODO(), filterDetail, bson.M{"$set": updateDetail})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func UpdateChatRoomByOID(oid primitive.ObjectID, updateDetail bson.M) (*mongo.UpdateResult, error) {

	result, err := database.ChatRoomCollection.UpdateOne(context.TODO(), bson.M{"_id": oid}, bson.M{"$set": updateDetail})

	return result, err
}

func DeleteChatRoomByOID(oid primitive.ObjectID) error {
	_, err := database.ChatRoomCollection.DeleteOne(context.TODO(), bson.M{"_id": oid})
	return err
}

// Find the chatroom without messages
// showing what's in the chatroom to customer
func FindChatRoomDetailWithoutMessages(matchingCond *[]bson.M, skip int, limit int) ([]*ChatRoomDetail, error) { // will populate members
	var chatRooms []*ChatRoomDetail

	var pipeline = []bson.M{}

	if matchingCond != nil {
		// Find match
		pipeline = append(pipeline, *matchingCond...)
	}

	pipeline = append(pipeline, []bson.M{

		// unwind the members
		bson.M{"$unwind": bson.M{"path": "$members"}},

		// get the member detail
		bson.M{
			"$lookup": bson.M{
				"from": "user",
				"let":  bson.M{"member": "$members"},
				"pipeline": bson.A{
					bson.M{"$match": bson.M{"$expr": bson.M{"$eq": bson.A{"$_id", "$$member"}}}},
					bson.M{"$project": bson.M{"_id": 1, "gender": 1, "domain": 1, "major": 1, "photoURL": 1, "role": 1, "email": 1}},
				},
				"as": "member",
			},
		},

		// Project
		bson.M{
			"$project": bson.M{
				"_id":       1,
				"member":    bson.M{"$arrayElemAt": bson.A{"$member", 0}},
				"isGroup":   1,
				"createdAt": 1,
				"groupName": 1,
			},
		},

		// Group

		bson.M{
			"$group": bson.M{
				"_id": bson.M{
					"_id":       "$_id",
					"createdAt": "$createdAt",
					"isGroup":   "$isGroup",
					"groupName": "$groupName",
				},
				"members": bson.M{"$push": bson.M{
					"_id":      "$member._id",
					"domain":   "$member.domain",
					"email":    "$member.email",
					"photoURL": "$member.photoURL",
					"major":    "$member.major",
					"role":     "$member.role",
					"gender":   "$member.gender",
				}},
			},
		},

		// Project

		bson.M{
			"$project": bson.M{
				"_id":       "$_id._id",
				"createdAt": "$_id.createdAt",
				"isGroup":   "$_id.isGroup",
				"groupName": "$_id.groupName",
				"members":   1,
			},
		},
	}...,
	)

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

	result, err := database.ChatRoomCollection.Aggregate(context.TODO(), pipeline)

	if result != nil {
		defer result.Close(context.TODO())
	}

	if err != nil {
		return nil, err
	}

	err = result.All(context.TODO(), &chatRooms)

	if err != nil {
		return nil, err
	}

	return chatRooms, nil
}

// WE don't need this right now
func FindMessagesForChatRoomByTime(chatRoomOID primitive.ObjectID, startTime time.Time) (*[]Message, error) {

	// 怎麼測定那些Message我應該汲取? -> 用ID 或是 時間?, id 有可能會消失, 但時間不會 時間還可以Sort. 用時間的缺點 -> 不確定是否為此ChatRoom
	// A: 時間, 因為用ID後也要取時間來Sort, 直接用$gt取值後再Sort

	var chatRooms []*ChatRoomDetail

	pipeline := []bson.M{
		bson.M{
			"$match": bson.M{
				"_id": chatRoomOID,
			},
		},

		bson.M{
			"$sort": bson.M{
				"messages.createdAt": -1,
			},
		},

		bson.M{"$project": bson.M{
			"messages": bson.M{
				"$filter": bson.M{
					"input": "$messages",
					"as":    "m",
					"cond":  bson.M{"$lte": bson.A{"$$m.createdAt", startTime}},
				},
			},
		}},
	}

	result, err := database.ChatRoomCollection.Aggregate(context.TODO(), pipeline)
	if result != nil {
		defer result.Close(context.TODO())
	}

	if err != nil {
		return nil, err
	}

	err = result.All(context.TODO(), &chatRooms)

	if err != nil {
		return nil, err
	}

	log.Println(chatRooms)

	return &chatRooms[0].Messages, nil
}

// Getting message here

// bson.M{
// 	"$match": bson.M{
// 		"_id": oid,
// 	},
// },
// Find the chatroom with messages

// find the chatroom with the last 50 messages

func FindMessagesForChatRoomByStartOID(userOID primitive.ObjectID, chatRoomOID primitive.ObjectID, StartOID primitive.ObjectID, retreiveNum int) (*[]Message, error) {

	var chatRooms []*ChatRoomDetail

	pipeline := []bson.M{
		bson.M{
			"$match": bson.M{
				"$and": bson.A{
					bson.M{
						"_id": chatRoomOID,
					},
					bson.M{
						"member": bson.M{"$in": bson.A{
							userOID,
						}},
					},
				},
			},
		},

		// do populate the author here

		bson.M{
			"$sort": bson.M{
				"messages.createdAt": -1,
			},
		},

		bson.M{
			"$project": bson.M{
				"messages": 1,
				"midx":     bson.M{"$indexOfArray": bson.A{"$messages._id", StartOID}},
				"_id":      1,
			},
		},

		bson.M{
			"$project": bson.M{
				"messages": bson.M{"$slice": bson.A{"$messages", bson.M{"$substact": bson.A{"$midx", retreiveNum}}, retreiveNum}},
				"_id":      1,
			},
		},

		// We should populate the author here
		// Get the author detail

		// First

		bson.M{
			"$unwind": bson.M{
				"path": "$messages",
			},
		},

		// bson.M{
		// 	"from":         "user",
		// 	"localField":   "messages.author",
		// 	"foreignField": "_id",
		// 	"as":           "messages.author",
		// },

		// Second

		// Look up for message authors
		bson.M{
			"$lookup": bson.M{
				"from": "user",
				"let":  bson.M{"messages": "$messages"},
				"pipeline": bson.A{
					bson.M{"match": bson.M{"$expr": bson.M{"$eq": bson.A{"$_id", "$$messages.author"}}}},
					bson.M{"_id": 1., "gender": 1, "domain": 1, " major": 1, "photoURL": 1, "role": 1, "email": 1},
				},
				"as": "messages.author",
			},
		},

		// Group the message back

		bson.M{"$group": bson.M{
			"_id":      bson.M{"_id": "$_id"},
			"messages": bson.M{"$push": bson.M{"_id": "$messages._id", "author": bson.M{"$arrayElemAt": bson.A{"$messages.author", 0}}, "likeBy": "$messages.likeBy", "readBy": "$messages.readBy", "messageType": "$messages.messageType", "content": "$messages.content", "createdAt": "$messages.createdAt"}},
		}},

		bson.M{
			"$project": bson.M{
				"_id":      "$_id._id",
				"messages": 1,
			},
		},
	}

	result, err := database.ChatRoomCollection.Aggregate(context.TODO(), pipeline)
	if result != nil {
		defer result.Close(context.TODO())
	}

	if err != nil {
		return nil, err
	}

	err = result.All(context.TODO(), &chatRooms)

	if err != nil {
		return nil, err
	}

	return &chatRooms[0].Messages, nil
}

func WatchChatRooms(chatRooms []primitive.ObjectID) (*mongo.ChangeStream, error) {
	pipeline := []bson.M{
		bson.M{
			"$in": bson.M{"documentKey._id": chatRooms},
		},
	}

	// changeStreamOption := options.ChangeStream().SetFullDocument(options.UpdateLookup) // 不用FullDocument, 不然每次會更新整個聊天室, 只有第一次打開時要Loading全部, 用Time來確定
	collectionStream, err := database.ChatRoomCollection.Watch(context.TODO(), pipeline)

	if err != nil {
		return nil, err
	}

	return collectionStream, nil
}
