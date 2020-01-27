package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

)

type Message struct {
	ID          primitive.ObjectID   `json:"_id" bson:"_id,omitempty"`
	Author      primitive.ObjectID   `json:"author" bson:"author"`
	MessageType int                  `json:"messageType" bson:"messageType"`
	Content     string               `json:"content" bson:"content"`
	CreatedAt   time.Time            `json:"createdAt" bson:"createdAt"`
	LikedBy     []primitive.ObjectID `json:"likeBy" bson:"likeBy"`
	ReadBy      []primitive.ObjectID `json:"readBy" bson:"readBy"`
	// We Don't need the belongChatRoom since we save the messages as a field of chatRoom
}

type MessageWithAuthorDetail struct {
	ID          primitive.ObjectID   `json:"_id" bson:"_id,omitempty"`
	Author      User                 `json:"author" bson:"author"`
	MessageType int                  `json:"messageType" bson:"messageType"`
	Content     string               `json:"content" bson:"content"`
	CreatedAt   time.Time            `json:"createdAt" bson:"createdAt"`
	LikedBy     []primitive.ObjectID `json:"likeBy" bson:"likeBy"`
	ReadBy      []primitive.ObjectID `json:"readBy" bson:"readBy"`
}
