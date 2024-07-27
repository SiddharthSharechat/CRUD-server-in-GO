package Controllers

import (
	"fmt"
	"github.com/SiddharthSharechat/CRUDGo/Initializers"
	"github.com/SiddharthSharechat/CRUDGo/Mappers"
	"github.com/SiddharthSharechat/CRUDGo/Models"
	"github.com/SiddharthSharechat/CRUDGo/Repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"strconv"
	"time"
)

func UsersCreate(c *gin.Context) {

	var body struct {
		Name     string
		Email    string
		Location string
	}

	c.Bind(&body)

	user := Models.User{
		Id:        uuid.New().String(),
		Name:      body.Name,
		Email:     body.Email,
		Location:  body.Location,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result := Initializers.Db.Create(&user)

	if result.Error != nil {
		c.Status(400)
		return
	}

	userResponse := Mappers.UserResponseMapper(user)
	c.JSON(200, gin.H{
		"user": userResponse,
	})

}

func UsersGet(c *gin.Context) {
	id := c.Param("id")
	var userResCache Models.UserResponse
	res := Repository.GetValue(id, &userResCache)
	if res {
		c.JSON(200, gin.H{
			"user":       userResCache,
			"from_cache": true,
		})
		return
	}

	var user Models.User
	Initializers.Db.First(&user, "id = ?", id)

	userResponse := Mappers.UserResponseMapper(user)
	Repository.SetValue(id, userResponse)
	c.JSON(200, gin.H{
		"user":       userResponse,
		"from_cache": false,
	})
}

func UsersUpdate(c *gin.Context) {

	id := c.Param("id")
	var body struct {
		Name     string
		Email    string
		Location string
	}

	c.Bind(&body)

	var user Models.User
	Initializers.Db.First(&user, "id = ?", id)

	Initializers.Db.Model(&user).Updates(Models.User{
		Name:      body.Name,
		Email:     body.Email,
		Location:  body.Location,
		UpdatedAt: time.Now(),
	},
	)

	userResponse := Mappers.UserResponseMapper(user)

	var userResCache Models.UserResponse
	res := Repository.GetValue(id, userResCache)
	if res {
		Repository.SetValue(id, userResponse)
	}

	c.JSON(200, gin.H{
		"user": userResponse,
	})
}

func UsersDelete(c *gin.Context) {
	id := c.Param("id")

	var userResCache Models.UserResponse
	res := Repository.GetValue(id, userResCache)
	if res {
		Repository.Expire(id)
	}

	var user Models.User
	Initializers.Db.First(&user, "id = ?", id)

	Initializers.Db.Delete(&user)

	userResponse := Mappers.UserResponseMapper(user)
	c.JSON(200, gin.H{
		"deleted_record": userResponse,
	})
}

func UsersGetPaginated(c *gin.Context) {
	location := c.Query("location")
	limitStr := c.DefaultQuery("limit", "5")
	pageStr := c.DefaultQuery("page", "1")

	limit, _ := strconv.Atoi(limitStr)
	page, _ := strconv.Atoi(pageStr)

	key := fmt.Sprintf("%s:%s:%s", location, page, limit)

	var storedUsers []Models.UserResponse
	isCached := Repository.LGet(key, &storedUsers)
	if isCached {
		c.JSON(200, gin.H{
			"page":       page,
			"limit":      limit,
			"users":      storedUsers,
			"from cache": true,
		})
		return
	}

	offset := (page - 1) * limit

	var users []Models.User
	Initializers.Db.Where("location = ?", location).Order("created_at ASC").Offset(offset).Limit(limit).Find(&users)

	var userResponses []Models.UserResponse
	for _, user := range users {
		userResponse := Mappers.UserResponseMapper(user)
		Repository.RPush(key, userResponse)
		userResponses = append(userResponses, userResponse)
	}

	c.JSON(200, gin.H{
		"page":  page,
		"limit": limit,
		"users": userResponses,
	})
}
