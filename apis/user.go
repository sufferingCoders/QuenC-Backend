package apis

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"quenc/models"
	"quenc/utils"

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
	var user models.User
	err := c.ShouldBindJSON(&user)

	fmt.Printf("err %+v", err)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": err,
			"msg": "Cannot bind the given LoginInfo",
		})
		return
	}

	if foundUser, err := models.FindUserByEmail(user.Email); foundUser != nil || err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err":       err,
			"msg":       "The email has been used",
			"user":      user,
			"foundUser": foundUser,
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"err":  err,
			"msg":  "Cannot hash the Password",
			"user": user,
		})
		return
	}
	user.Password = string(hashedPassword)
	user.EmailVerified = false

	InsertedID, err := models.AddUser(&user)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"err": err,
			"msg": "Unable to add the user to Database",
		})
		return
	}

	user.ID = InsertedID.(primitive.ObjectID)
	user.Password = ""

	authToken, err := utils.GenerateAuthToken(user.ID.Hex())

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"err":  err,
			"msg":  "Cannot Generate the Auth Token",
			"user": user,
		})
		return
	}

	// Send verification email here

	err = models.SendingVerificationEmail(&user)

	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadGateway, gin.H{
				"err": err,
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
		c.AbortWithStatusJSON(
			http.StatusBadGateway, gin.H{
				"err": err,
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
			"uid":   uid,
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
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": err,
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
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err":       err,
			"msg":       "Email or Passwrod is not correct",
			"loginInfo": loginInfo,
		})
		return
	}

	authToken, err := utils.GenerateAuthToken(user.ID.Hex())

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"err":  err,
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
	var updateUserInfo UpdateUserInfo
	err := c.ShouldBindJSON(&updateUserInfo)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": err,
			"msg": "Cannot bind the given data with UpdateUserInfo",
		})
		return
	}

	if updateUserInfo.UpdateDetail["password"] != nil {
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

	UpsertedID, err := models.UpdateUserByOID(user.ID, updateUserInfo.UpdateDetail)

	fmt.Printf("error is %+v\n", err)

	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"err":            err,
				"msg":            "Cannot update this user",
				"user":           user,
				"UpdateUserInfo": updateUserInfo,
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"UpsertedID":     UpsertedID,
			"UpdateUserInfo": updateUserInfo,
		},
	)

}

func SubscribeUser(c *gin.Context) {

	var upGrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)

	defer ws.Close()

	if err != nil {
		errStr := fmt.Sprintf("The websocket is not working due to the error: %+v \n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errStr,
		})
		return
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
