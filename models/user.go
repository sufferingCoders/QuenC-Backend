package models

import (
	"context"
	"net/smtp"
	"os"
	"quenc/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"

)

// User - User Schema
type User struct {
	ID            primitive.ObjectID   `json:"_id" bson:"_id,omitempty"`
	Domain        string               `json:"domain" bson:"domain"`
	Email         string               `json:"email" bson:"email"`
	Password      string               `json:"password" bson:"password"`
	PhotoURL      string               `json:"photoURL" bson:"photoURL"`
	Role          int                  `json:"role" bson:"role"`
	Gender        int                  `json:"gender" bson:"gender"`
	EmailVerified bool                 `json:"emailVerified" bson:"emailVerified"`
	LastSeen      primitive.DateTime   `json:"lastSeen" bson:"lastSeen"`
	Dob           primitive.DateTime   `json:"dob" bson:"dob"`
	CreatedAt     primitive.DateTime   `json:"createdAt" bson:"createdAt"`
	ChatRooms     []primitive.ObjectID `json:"chatRooms" bson:"chatRooms"`
	Friends       []primitive.ObjectID `json:"friends" bson:"friends"`
	Major         []primitive.ObjectID `json:"major" bson:"major"`
	LikePosts     []primitive.ObjectID `json:"likePosts" bson:"likePosts"`
	LikeComments  []primitive.ObjectID `json:"likeComments" bson:"likeComments"`
	SavedPosts    []primitive.ObjectID `json:"savedPosts" bson:"savedPosts"`
}

var ( // Changing to env variables
	host                          = "smtp.gmail.com:587"
	username                      = os.Getenv("EMAIL_SENDING_USERNAME")
	password                      = os.Getenv("EMAIL_SENDING_PASSWORD")
	projectionForRemovingPassword = bson.D{
		{"password", 0},
	}
	verificationBaseURL = "https://shopping-au.appspot.com/user/activate/"
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

	err := database.UserCollection.FindOne(context.TODO(), bson.M{"email": email},
		options.FindOne().SetProjection(projectionForRemovingPassword)).Decode(&user)

	return &user, err
}

// FindUsers - Find Multiple Users by filterDetail
func FindUsers(filterDetail bson.M) ([]*User, error) {
	var users []*User
	result, err := database.UserCollection.Find(context.TODO(), filterDetail,
		options.Find().SetProjection(projectionForRemovingPassword))
	defer result.Close(context.TODO())

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
	err := database.DB.Collection("user").FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return nil, err
	}
	return &user, nil
}
