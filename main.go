package main

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type User struct {
	Id               string         `json:"id"`
	Address          string         `json:"address"`
	Upline           string         `json:"upline"`
	TelegramUsername sql.NullString `json:"telegramUsername"`
	BlockNumber      uint           `json:"blockNumber"`
	Rank             sql.NullInt16  `json:"rank"`
}

var users = []User{
	{Id: "0", Address: "0xf4faaa84d20968f0cecd1ec7a0d1e10cf20b552c", Upline: "0x0", BlockNumber: 0},
	{Id: "1", Address: "0x976c30089cd5da634dde85432b336251c2634c3e", Upline: "0xf4faaa84d20968f0cecd1ec7a0d1e10cf20b552c", BlockNumber: 0},
	{Id: "2", Address: "0x4320cabe04a8308de7f813fd318758ef32655030", Upline: "0x976c30089cd5da634dde85432b336251c2634c3e", BlockNumber: 0},
}

// DB set up
func setupDB() *sql.DB {
	// dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	connStr := "postgres://postgres:beritaunik1234@localhost:5432/valhalla?sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		panic(err)
	}

	return db
}

func getUsers(c *gin.Context) {
	db := setupDB()
	defer db.Close() // Close the database connection after the function returns

	rows, err := db.Query(`SELECT * FROM "User"`)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Address, &user.TelegramUsername, &user.Upline, &user.BlockNumber, &user.Rank)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		user.TelegramUsername.String = "apa boss"
		user.Rank.Int16 = 0
		users = append(users, user)
	}

	c.IndentedJSON(http.StatusOK, users)
}

func addUser(c *gin.Context) {
	var newUser User
	if err := c.BindJSON(&newUser); err != nil {
		return
	}

	// create id from last index
	lastUserId := len(users)
	lastUserId++
	// assigne id
	newUser.Id = strconv.Itoa(lastUserId)

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
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"})
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
	// this will be printed in the terminal, confirming the connection to the database
	router := gin.Default()
	router.GET("api/get_users", getUsers)
	router.GET("api/get_users/:id", getUserById)
	router.POST("api/add_user", addUser)
	router.PUT("api/dell_users/:id", deletUserById)
	router.Run("localhost:8080")
}
