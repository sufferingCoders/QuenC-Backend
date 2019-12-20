package utils

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SetupFindOptions - Setting up the FindOptions for the Query
func SetupFindOptions(findOptions *options.FindOptions, c *gin.Context) {

	skip, limit, sort := c.Query("skip"), c.Query("limit"), c.Query("sort")

	// findOptions := options.Find()
	if strings.TrimSpace(skip) != "" {
		inputSkip, err := strconv.ParseInt(skip, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err":  err,
				"msg":  "Cannot setup skip",
				"skip": skip,
			})
			return
		}
		findOptions.SetSkip(inputSkip)
	}

	if strings.TrimSpace(limit) != "" {
		inputLimit, err := strconv.ParseInt(limit, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err":   err,
				"msg":   "Cannot setup limit",
				"limit": limit,
			})
			return
		}
		findOptions.SetLimit(inputLimit)
	}

	sortMap := map[string]int{}
	if strings.TrimSpace(sort) != "" {
		if s := strings.Split(sort, "_"); len(s) == 2 {
			sortOrd, err := strconv.Atoi(s[1])
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"err":  err,
					"msg":  "Cannot get the sort order",
					"s[1]": s[1],
					"s[0]": s[0],
				})
				return
			}
			// fmt.Printf("s is %+v\n", s)
			// fmt.Printf("s[0] is %+v\n", s[0])
			// fmt.Printf("s[1] is %+v\n", s[1])
			// fmt.Printf("sortOrd is %+v\n", sortOrd)

			sortMap[s[0]] = sortOrd
		} else {
			sortMap[sort] = -1
		}

		findOptions.SetSort(sortMap)
	}
}
