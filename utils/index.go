package utils

import (
	"fmt"
	"net/http"
	"os"
	"quenc/models"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

)

// SetupFindOptions - Setting up the FindOptions for the Query
func SetupFindOptions(findOptions *options.FindOptions, c *gin.Context) error {

	skip, limit, sort := c.Query("skip"), c.Query("limit"), c.Query("sort")

	sortMap := map[string]int{}
	if strings.TrimSpace(sort) != "" { // Expect "createdAt_1,likeCount_-1"
		sortedRequire := strings.Split(sort, ",")
		for _, sr := range sortedRequire {
			if s := strings.Split(sr, "_"); len(s) == 2 {
				sortOrd, err := strconv.Atoi(s[1])
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"err":  err,
						"msg":  "Cannot get the sort order",
						"s[1]": s[1],
						"s[0]": s[0],
					})
					return err
				}
				// fmt.Printf("s is %+v\n", s)
				// fmt.Printf("s[0] is %+v\n", s[0])
				// fmt.Printf("s[1] is %+v\n", s[1])
				// fmt.Printf("sortOrd is %+v\n", sortOrd)

				sortMap[s[0]] = sortOrd
			} else {
				sortMap[sort] = -1
			}
		}

		findOptions.SetSort(sortMap)
	}

	// findOptions := options.Find()
	if strings.TrimSpace(skip) != "" {
		inputSkip, err := strconv.ParseInt(skip, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err":  err,
				"msg":  "Cannot setup skip",
				"skip": skip,
			})
			return err
		}
		findOptions.SetSkip(inputSkip)
	}

	if strings.TrimSpace(limit) != "" {
		inputLimit, err := strconv.ParseInt(limit, 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"err":   err,
				"msg":   "Cannot setup limit",
				"limit": limit,
			})
			return err
		}
		findOptions.SetLimit(inputLimit)
	}

	return nil
}

func GetOID(id string, c *gin.Context) *primitive.ObjectID {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		errStr := fmt.Sprintf("Cannot transfrom the given id to ObjectId: %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": errStr,
			"id":  id,
		})
		return nil
	}
	return &oid
}

// GenerateAuthToken - Generate the Auth token for given id
func GenerateAuthToken(id string) (interface{}, error) {
	/*
		Method for generating the token
	*/
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"_id": id,
		// "exp":   time.Now().Add(time.Hour * 2).Unix(),
	})

	authToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		return nil, err
	}

	return authToken, nil
}

// GetUserFromContext - Return User Object
func GetUserFromContext(c *gin.Context) *models.User {
	var user *models.User
	userStr, ok := c.Get("user")

	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"err": "Cannot retrieve the user after token authorization",
			"msg": "Cannot retrieve the user after token authorization",
		})
		return nil
	} else {
		user = userStr.(*models.User)
	}

	return user
}

func GetDomainFromEmail(email string) string {
	emailParts := strings.Split(email, "@")
	if len(emailParts) > 2 {
		return ""
	}
	return emailParts[1]
}

func GetSkipLimitSortFromContext(c *gin.Context) (*int, *int, *string, error) {
	skipStr := c.Query("skip")
	limitStr := c.Query("limit")
	skip, err := strconv.Atoi(skipStr)
	if err != nil {
		errStr := fmt.Sprintf("Cannot convert the given skip: %+v", err)
		c.AbortWithStatusJSON(
			http.StatusBadRequest, gin.H{
				"err":     errStr,
				"skipStr": skipStr,
			},
		)
		return nil, nil, nil, err

	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		errStr := fmt.Sprintf("Cannot convert the given limit: %+v", err)
		c.AbortWithStatusJSON(
			http.StatusBadRequest, gin.H{
				"err":      errStr,
				"limitStr": limitStr,
			},
		)
		return nil, nil, nil, err
	}

	var sort *string
	sortStr := c.Query("sort")
	if strings.TrimSpace(sortStr) == "" {
		sort = nil
	} else {
		sort = &sortStr
	}

	return &skip, &limit, sort, nil
}
