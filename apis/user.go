package apis

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"quenc/models"
	"quenc/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"

)

type SingupInfo struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ChangingPasswordInfo struct {
	OldPassword string `json:"oldPassword" bson:"oldPassword"`
	NewPasswrod string `json:"newPassword" bson:"newPassword"`
}

type LoginInfo struct {
	Eamil    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

type UpdateUserInfo struct {
	UpdateDetail map[string]interface{} `json:"updateDetail" bson:"updateDetail"`
}

func SingupUser(c *gin.Context) {

	// Hava to use the uni email
	var singupInfo SingupInfo
	err := c.ShouldBindJSON(&singupInfo)

	// Create a User in the backend

	if err != nil {
		errStr := fmt.Sprint("Cannot bind the given signup info: %+v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": errStr,
			"msg": "Cannot bind the given LoginInfo",
		})
		return
	}

	if foundUser, _ := models.FindUserByEmail(singupInfo.Email); foundUser != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err":        "The email has been used",
			"msg":        "The email has been used",
			"singupInfo": singupInfo,
			"foundUser":  foundUser,
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(singupInfo.Password), bcrypt.DefaultCost)

	if err != nil {
		errStr := fmt.Sprint("Cannot has the password: %+v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"err":        errStr,
			"msg":        "Cannot hash the Password",
			"singupInfo": singupInfo,
		})
		return
	}

	// Creating user here

	userDomain := utils.GetDomainFromEmail(singupInfo.Email)

	if valid := utils.CheckDomainValid(userDomain); !valid {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"err":        "Please use supported email domain for registering",
			"msg":        "Please use supported email domain for registering",
			"singupInfo": singupInfo,
		})
		return
	}

	user := models.User{
		Domain:        userDomain,
		Email:         singupInfo.Email,
		Name:          "",
		Password:      string(hashedPassword),
		PhotoURL:      "",
		Major:         "",
		Role:          1,
		Gender:        -1,
		EmailVerified: false,
		Dob:           "",
		LastSeen:      time.Now(),
		CreatedAt:     time.Now(),
		LikePosts:     []primitive.ObjectID{},
		LikeComments:  []primitive.ObjectID{},
		ChatRooms:     []primitive.ObjectID{},
		Friends:       []primitive.ObjectID{},
		SavedPosts:    []primitive.ObjectID{},
	}

	InsertedID, err := models.AddUser(&user)

	if err != nil {
		errStr := fmt.Sprint("Unable to add the user to Database: %+v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"err": errStr,
			"msg": "Unable to add the user to Database",
		})
		return
	}

	user.ID = InsertedID.(primitive.ObjectID)
	user.Password = ""

	authToken, err := utils.GenerateAuthToken(user.ID.Hex())

	if err != nil {
		errStr := fmt.Sprint("Cannot Generate the Auth Token: %+v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"err":  errStr,
			"msg":  "Cannot Generate the Auth Token",
			"user": user,
		})
		return
	}

	// Send verification email here

	err = models.SendingVerificationEmail(&user)

	if err != nil {
		errStr := fmt.Sprint("Cannot send the email to this account: %+v", err)
		c.AbortWithStatusJSON(
			http.StatusBadGateway, gin.H{
				"err": errStr,
				"msg": "Cannot send the email to this account",
			},
		)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":  user,
		"token": authToken,
		"msg":   "Email verification has been sent to" + user.Email,
	})
}

func TokenAutoLogin(c *gin.Context) {
	user := utils.GetUserFromContext(c)

	c.JSON(
		http.StatusOK,
		gin.H{
			"user": user,
		},
	)
}

func SendVerificationEmailForUser(c *gin.Context) {
	// Have to login to do this
	user := utils.GetUserFromContext(c)

	err := models.SendingVerificationEmail(user)

	if err != nil {
		errStr := fmt.Sprintf("Cannot send the email to this account: %+v", err)
		c.AbortWithStatusJSON(
			http.StatusBadGateway, gin.H{
				"err": errStr,
				"msg": "Cannot send the email to this account",
			},
		)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
		"msg":  "Email verification has been sent to" + user.Email,
	})
}

func ActivateUserEmail(c *gin.Context) {
	uid := c.Param("uid")

	uOID, err := primitive.ObjectIDFromHex(uid)

	if err != nil {
		errStr := fmt.Sprintf("Cannot get ObjectID from Hex: %+v", err)
		c.HTML(http.StatusBadRequest, "EmailVerificationFail.tmpl", gin.H{
			"email": uid,
			"error": errStr,
			"msg":   "無法找到此使用者",
		})
		return
	}

	user, err := models.FindUserByOID(uOID)

	if err != nil {
		c.HTML(http.StatusBadRequest, "EmailVerificationFail.tmpl", gin.H{
			"email": user.Email,
			"error": err.Error,
			"msg":   "無法找到此使用者",
		})
		return
	}

	_, err = models.UpdateUserByOID(uOID, bson.M{"emailVerified": true})

	if err != nil {

		c.HTML(http.StatusBadRequest, "EmailVerificationFail.tmpl", gin.H{
			"email": user.Email,
			"error": err.Error,
			"msg":   "無法激活此使用者",
		})
		return
	}

	user.EmailVerified = true

	c.HTML(http.StatusOK, "EmailVerificationSuccessful.tmpl", gin.H{
		"email": user.Email,
	})
}

func LoginUser(c *gin.Context) {
	var loginInfo LoginInfo
	err := c.ShouldBindJSON(&loginInfo)

	if err != nil {
		errStr := fmt.Sprintf("Cannot bind the given LoginInfo: %+v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": errStr,
			"msg": "Cannot bind the given LoginInfo",
		})
		return
	}

	// Checking that the email and password are provided
	if loginInfo.Eamil == "" || loginInfo.Password == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err":       "Must Provide Password and Email",
			"msg":       "Must Provide Password and Email",
			"loginInfo": loginInfo,
		})
		return
	}

	user, err := models.CheckingTheAuth(loginInfo.Eamil, loginInfo.Password)

	if err != nil {
		errStr := fmt.Sprintf("Email or Passwrod is not correct: %+v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err":       errStr,
			"msg":       "Email or Passwrod is not correct",
			"loginInfo": loginInfo,
		})
		return
	}

	authToken, err := utils.GenerateAuthToken(user.ID.Hex())

	if err != nil {
		errStr := fmt.Sprintf("Cannot generate auth token for this user: %+v", err)

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"err":  errStr,
			"msg":  "Cannot generate auth token for this user",
			"user": user,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": authToken,
		"user":  user,
	})
}

func UpdateUser(c *gin.Context) {
	var updateFields map[string]interface{}
	err := c.ShouldBindJSON(&updateFields)

	if _, ok := updateFields["dob"]; ok {
		_, err := time.Parse("2006-01-02 15:04:05.000Z", updateFields["dob"].(string))
		if err != nil {
			errStr := fmt.Sprint("The given DOB is not valid time %+v", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"err": errStr,
				"msg": "The given DOB is not valid time",
			})
			return
		}
	}

	delete(updateFields, "_id")
	delete(updateFields, "email")
	delete(updateFields, "createdAt")
	delete(updateFields, "emailVerified")
	delete(updateFields, "chatRooms")
	delete(updateFields, "likePosts")
	delete(updateFields, "likeComments")
	delete(updateFields, "friends")
	delete(updateFields, "savedPosts")

	if err != nil {
		errStr := fmt.Sprintf("Cannot bind the given data with UpdateUserInfo: %+v", err)

		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": errStr,
			"msg": "Cannot bind the given data with UpdateUserInfo",
		})
		return
	}

	if updateFields["password"] != nil {

		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": "Using change-password node to change password",
			"msg": "Using change-password node to change password",
		})
		return
	}

	user := utils.GetUserFromContext(c)

	if user == nil {
		return
	}
	// Admin should be able to do

	UpsertedID, err := models.UpdateUserByOID(user.ID, updateFields)

	fmt.Printf("error is %+v\n", err)

	if err != nil {
		errStr := fmt.Sprintf("Cannot update this user: %+v", err)

		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"err":            errStr,
				"msg":            "Cannot update this user",
				"user":           user,
				"UpdateUserInfo": updateFields,
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"UpsertedID":     UpsertedID,
			"UpdateUserInfo": updateFields,
		},
	)

}

func ToggleFunc(field string) gin.HandlerFunc {
	return func(c *gin.Context) {
		condition := c.Param("codition")
		id := c.Param("id")
		OID := utils.GetOID(id, c)
		if OID != nil {
			return
		}

		var adding bool

		if condition == "1" {
			adding = true
		} else if condition == "0" {
			adding = false
		} else {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"err": "You didn't give a compatible condition",
			})
			return
		}

		user := utils.GetUserFromContext(c)

		if user == nil {
			return
		}

		result, err := models.ToggleElementToUserArray(field, adding, *OID, user.ID)

		if err != nil {
			errStr := fmt.Sprintf("Cannot toggle the condition: %+v", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"err": errStr,
			})
			return
		}

		c.JSON(
			http.StatusOK, gin.H{
				"result": result,
			},
		)
	}

}

func SubscribeUser(c *gin.Context) {

	var upGrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	user := utils.GetUserFromContext(c)

	if user == nil {
		return
	}

	stream, err := models.WatchUserByOID(user.ID)

	if err != nil {
		errStr := fmt.Sprintf("Cannot get the stream: %+v \n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errStr,
		})
		return
	}

	defer stream.Close(context.TODO())

	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)

	defer ws.Close()

	if err != nil {
		errStr := fmt.Sprintf("The websocket is not working due to the error: %+v \n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errStr,
		})
		return
	}

	for {
		ok := stream.Next(context.TODO())
		if ok {
			next := stream.Current

			var m map[string]interface{}

			err := bson.Unmarshal(next, &m)
			if err != nil {
				log.Print(err)
				break
			}

			err = ws.WriteJSON(m["fullDocument"].(map[string]interface{}))
			if err != nil {
				log.Print(err)
				break
			}
		}
	}
}
