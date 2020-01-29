package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"quenc/models"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

)

func UserAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")

		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"err":      "Token is not provided",
				"msg":      "Token is not provided",
				"tokenStr": tokenStr,
			})
			return
		}

		if s := strings.Split(tokenStr, " "); len(s) == 2 {
			tokenStr = s[1]
		}

		token, _ := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"err":      "Invalid Token",
					"msg":      "Cannot parse the given token",
					"token":    token,
					"tokenStr": tokenStr,
				})
				return nil, fmt.Errorf("Invalid Token")
			}

			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		claims := token.Claims.(jwt.MapClaims)

		inputClaim_userID := claims["_id"].(string)

		oid, err := primitive.ObjectIDFromHex(inputClaim_userID)

		if err != nil {
			errStr := fmt.Sprintf("Cannot get the ObejctId: %+v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"err": errStr,
				"id":  inputClaim_userID,
			})
		}

		user, err := models.FindUserByOID(oid)

		if err != nil {
			errStr := fmt.Sprintf("Cannot find the user during authroization checking: %+v", err)

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"err":      errStr,
				"msg":      "Cannot find the user during authroization checking",
				"token":    token,
				"tokenStr": tokenStr,
			})
			return
		}

		c.Set("user", user)

		if !token.Valid {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"token":    token,
					"tokenStr": tokenStr,
					"err":      "The token is not valid",
					"msg":      "The token is not valid",
				},
			)
			return
		}

		c.Next()

	}

}

func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")

		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"err":      "Token is not provided",
				"msg":      "Token is not provided",
				"tokenStr": tokenStr,
			})
			return
		}

		if s := strings.Split(tokenStr, " "); len(s) == 2 {
			tokenStr = s[1]
		}

		token, _ := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"err":      "Invalid Token",
					"msg":      "Cannot parse the given token",
					"token":    token,
					"tokenStr": tokenStr,
				})
				return nil, fmt.Errorf("Invalid Token")
			}

			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		claims := token.Claims.(jwt.MapClaims)

		inputClaim_userID := claims["_id"].(string)

		oid, err := primitive.ObjectIDFromHex(inputClaim_userID)

		if err != nil {
			errStr := fmt.Sprintf("Cannot get the ObejctId: %+v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"err": errStr,
				"id":  inputClaim_userID,
			})
		}

		user, err := models.FindUserByOID(oid)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"err":      err,
				"msg":      "Cannot find the user during authroization checking",
				"token":    token,
				"tokenStr": tokenStr,
			})
			return
		}

		if !user.IsAmin() {
			errStr := fmt.Sprintf("Unauthorised")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"err":  errStr,
				"user": user,
			})
			return
		}

		c.Set("user", user)

		if !token.Valid {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"token":    token,
					"tokenStr": tokenStr,
					"err":      "The token is not valid",
					"msg":      "The token is not valid",
				},
			)
			return
		}

		c.Next()

	}

}
