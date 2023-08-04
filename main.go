package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct {
	Id               string `json:"id"`
	Address          string `json:"address"`
	TelegramUsername string `json:"telegramUsername"`
	Upline           string `json:"upline"`
	BlockNumber      uint   `json:"blockNumber"`
	Rank             uint   `json:"rank"`
}

var users = []User{
	{Id: "0", Address: "0xf4faaa84d20968f0cecd1ec7a0d1e10cf20b552c", TelegramUsername: "", Upline: "0x0", BlockNumber: 0, Rank: 1},
	{Id: "1", Address: "0x976c30089cd5da634dde85432b336251c2634c3e", TelegramUsername: "", Upline: "0xf4faaa84d20968f0cecd1ec7a0d1e10cf20b552c", BlockNumber: 0, Rank: 0},
	{Id: "2", Address: "0x4320cabe04a8308de7f813fd318758ef32655030", TelegramUsername: "", Upline: "0x976c30089cd5da634dde85432b336251c2634c3e", BlockNumber: 0, Rank: 0},
}

func getUsers(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, users)
}

func addUser(c *gin.Context) {
	var newUser User
	if err := c.BindJSON(&newUser); err != nil {
		return
	}

	lastUserId := len(users)
	newUser.Id = fmt.Sprintf("%d", lastUserId)

	users = append(users, newUser)
	c.IndentedJSON(http.StatusCreated, newUser)
}

func getUserById(c *gin.Context) {
	id := c.Param("id")

	for _, a := range users {
		if a.Id == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

func deletUserById(c *gin.Context) {
	id := c.Param("id")

	var updatedUsers []User

	// Iterasi melalui slice dan salin elemen-elemen yang tidak dihapus ke slice baru
	for _, user := range users {
		if user.Id != id {
			updatedUsers = append(updatedUsers, user)
		}
	}
	users = updatedUsers

	c.IndentedJSON(http.StatusOK, users)
}

func main() {
	router := gin.Default()
	router.GET("api/get_users", getUsers)
	router.GET("api/get_users/:id", getUserById)
	router.POST("api/add_user", addUser)
	router.PUT("api/dell_users/:id", deletUserById)
	router.Run("localhost:8080")
}
