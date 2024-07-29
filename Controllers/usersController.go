package Controllers

import (
	"fmt"
	"github.com/SiddharthSharechat/CRUDGo/Initializers"
	"github.com/SiddharthSharechat/CRUDGo/Mappers"
	"github.com/SiddharthSharechat/CRUDGo/Models"
	"github.com/SiddharthSharechat/CRUDGo/Repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"strconv"
	"time"
)

const paginationKey = "pagination"

var body struct {
	Name     string `form:"name" binding:"required"`
	Email    string `form:"email" binding:"required,email"`
	Location string `form:"location" binding:"required"`
}

func UsersCreate(c *gin.Context) {

	err := c.Bind(&body)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

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
		c.JSON(400, gin.H{
			"message": result.Error.Error(),
		})
		return
	}

	Repository.ClearPaginationCache(paginationKey)

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
	result := Initializers.Db.First(&user, "id = ?", id)
	log.Printf("DB Get response : %v", result)

	userResponse := Mappers.UserResponseMapper(user)
	Repository.SetValue(id, userResponse)
	c.JSON(200, gin.H{
		"user":       userResponse,
		"from_cache": false,
	})
}

func UsersUpdate(c *gin.Context) {

	id := c.Param("id")

	err := c.Bind(&body)
	if err != nil {
		c.Status(500)
		return
	}

	var user Models.User
	res := Initializers.Db.First(&user, "id = ?", id)
	log.Printf("DB Get in update response : %v", res)

	Initializers.Db.Model(&user).Updates(Models.User{
		Name:      body.Name,
		Email:     body.Email,
		Location:  body.Location,
		UpdatedAt: time.Now(),
	},
	)

	userResponse := Mappers.UserResponseMapper(user)

	Repository.Expire(id)
	Repository.ClearPaginationCache(paginationKey)

	c.JSON(200, gin.H{
		"user": userResponse,
	})
}

func UsersDelete(c *gin.Context) {
	id := c.Param("id")

	Repository.Expire(id)

	var user Models.User
	Initializers.Db.First(&user, "id = ?", id)

	res := Initializers.Db.Delete(&user)
	log.Printf("DB Delete response : %v", res)

	Repository.ClearPaginationCache(paginationKey)

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
	if len(location) == 0 {
		c.JSON(400, gin.H{
			"message": "location is a required Parameter",
		})
		return
	}

	if limit <= 0 || page <= 0 {
		c.JSON(400, gin.H{
			"message": "Invalid limit or page",
		})
		return
	}

	key := fmt.Sprintf("%s:%d:%d", location, page, limit)

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
	res := Initializers.Db.Where("location = ?", location).Order("created_at ASC").Offset(offset).Limit(limit).Find(&users)
	log.Printf("Pagination DB response: %v", res)

	var userResponses []Models.UserResponse
	Repository.RPush(paginationKey, key)
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
