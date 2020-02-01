package models

import (
	"context"
	"net/smtp"
	"os"
	"quenc/database"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// User - User Schema
type User struct {
	ID             primitive.ObjectID   `json:"_id" bson:"_id,omitempty"`
	RandomChatRoom *primitive.ObjectID  `json:"randomChatRoom" bson:"randomChatRoom"`
	Name           string               `json:"name" bson:"name"`
	Domain         string               `json:"domain" bson:"domain"`
	Email          string               `json:"email" bson:"email"`
	Password       string               `json:"password" bson:"password"`
	PhotoURL       string               `json:"photoURL" bson:"photoURL"`
	Major          string               `json:"major" bson:"major"`
	Dob            string               `json:"dob" bson:"dob"`
	Role           int                  `json:"role" bson:"role"`
	Gender         int                  `json:"gender" bson:"gender"`
	EmailVerified  bool                 `json:"emailVerified" bson:"emailVerified"`
	LastSeen       time.Time            `json:"lastSeen" bson:"lastSeen"`
	CreatedAt      time.Time            `json:"createdAt" bson:"createdAt"`
	BlockedPosts   []primitive.ObjectID `json:"blockedPosts" bson:"blockedPosts"`
	BlockedUsers   []primitive.ObjectID `json:"blockedUsers" bson:"blockedUsers"`
	ChatRooms      []primitive.ObjectID `json:"chatRooms" bson:"chatRooms"`
	LikePosts      []primitive.ObjectID `json:"likePosts" bson:"likePosts"`
	LikeComments   []primitive.ObjectID `json:"likeComments" bson:"likeComments"`
	Friends        []primitive.ObjectID `json:"friends" bson:"friends"`
	SavedPosts     []primitive.ObjectID `json:"savedPosts" bson:"savedPosts"`
}

var ( // Changing to env variables
	// username                      = os.Getenv("EMAIL_SENDING_USERNAME")
	// password                      = os.Getenv("EMAIL_SENDING_PASSWORD")
	projectionForRemovingPassword = bson.D{
		{"password", 0},
	}
	verificationBaseURL = "http://quenc-hlc.appspot.com/user/email/activate/"
)

func (u *User) IsAmin() bool {
	return u.Role == 0
}

// AddUser - Adding User to MongoDB
func AddUser(inputUser *User) (interface{}, error) {

	result, err := database.UserCollection.InsertOne(context.TODO(), inputUser)
	return result.InsertedID, err
}

// UpdateUsers - Update User in MongoDB
func UpdateUsers(filterDetail bson.M, updateDetail bson.M) (*mongo.UpdateResult, error) {

	result, err := database.UserCollection.UpdateMany(context.TODO(), filterDetail, bson.M{"$set": updateDetail})

	return result, err
}

// UpdateUserByOID - Update User in MongoDB by its OID
func UpdateUserByOID(oid primitive.ObjectID, updateDetail bson.M) (*mongo.UpdateResult, error) {

	result, err := database.UserCollection.UpdateOne(context.TODO(), bson.M{"_id": oid}, bson.M{"$set": updateDetail})

	return result, err
}

// DeleteUserByOID - Delete user by its OID
func DeleteUserByOID(oid primitive.ObjectID) error {
	_, err := database.UserCollection.DeleteOne(context.TODO(), bson.M{"_id": oid})
	return err
}

// FindUserByOID - Find User by its OID
func FindUserByOID(oid primitive.ObjectID) (*User, error) {
	var user User

	err := database.UserCollection.FindOne(context.TODO(), bson.M{"_id": oid},
		options.FindOne().SetProjection(projectionForRemovingPassword)).Decode(&user)

	return &user, err
}

// FindUserByEmail - Find User by its email
func FindUserByEmail(email string) (*User, error) {
	var user User

	options.FindOne().SetProjection(projectionForRemovingPassword)
	err := database.UserCollection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, err
}

// FindUsers - Find Multiple Users by filterDetail
func FindUsers(filterDetail bson.M) ([]*User, error) {
	var users []*User
	result, err := database.UserCollection.Find(context.TODO(), filterDetail,
		options.Find().SetProjection(projectionForRemovingPassword))
	if result != nil {
		defer result.Close(context.TODO())
	}

	if err != nil {
		return nil, err
	}

	for result.Next(context.TODO()) {
		var elem User
		err := result.Decode(&elem)
		if err != nil {
			return nil, err
		}
		users = append(users, &elem)
	}

	return users, nil
}

// SendingVerificationEmail - Sending Email Verification to a User
func SendingVerificationEmail(user *User) error {
	var host = "smtp.gmail.com:587"
	var username = os.Getenv("EMAIL_SENDING_USERNAME")
	var password = os.Getenv("EMAIL_SENDING_PASSWORD")
	plainAuth := smtp.PlainAuth(host, username, password, "smtp.gmail.com")
	to := []string{user.Email}
	msg := []byte(
		"Subject: 激活您的Quenc帳號!\r\n" +
			"From: no-reply@gmail.com\r\n" +
			`Content-Type: text/plain;` +
			"\r\n" +
			"\r\n" +
			"您好，歡迎您加入昆嗑社群\r\n" +
			"請點擊以下的連結激活帳號：" +
			"\r\n" +
			verificationBaseURL + user.ID.Hex(),
	)
	err := smtp.SendMail(
		host,
		plainAuth,
		username,
		to,
		msg,
	)
	return err
}

// CheckingTheAuth - Checking if email and password are valid
func CheckingTheAuth(email string, password string) (*User, error) {
	var user User
	err := database.UserCollection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func WatchUser(pipeline []bson.M, changeStreamOption *options.ChangeStreamOptions) (*mongo.ChangeStream, error) {
	collectionStream, err := database.UserCollection.Watch(context.TODO(), pipeline, changeStreamOption)
	return collectionStream, err
}

func WatchUserByOID(oid primitive.ObjectID) (*mongo.ChangeStream, error) {
	// pipeline := mongo.Pipeline{bson.D{{"$match", bson.D{{"fullDocument._id", oid}}}}}
	pipeline := []bson.M{
		bson.M{
			"$match": bson.M{"fullDocument._id": oid},
		},
	}
	changeStreamOption := options.ChangeStream().SetFullDocument(options.UpdateLookup)
	stream, err := WatchUser(pipeline, changeStreamOption)
	return stream, err
}

// Liked Comment, Saved Comment, Saved Posts

func ToggleElementToUserArray(field string, adding bool, element primitive.ObjectID, uOID primitive.ObjectID) (*mongo.UpdateResult, error) {

	var result *mongo.UpdateResult
	var err error
	if adding {
		result, err = database.UserCollection.UpdateOne(context.TODO(), bson.M{"_id": uOID}, bson.M{"$push": bson.M{field: element}})
	} else {
		result, err = database.UserCollection.UpdateOne(context.TODO(), bson.M{"_id": uOID}, bson.M{"$pull": bson.M{field: element}})
	}
	return result, err
}
