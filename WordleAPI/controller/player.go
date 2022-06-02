package controller

import (
	"WordleAPI/helper/database"
	"WordleAPI/helper/hexcode"
	"github.com/gin-gonic/gin"
)

// Creates one player in the database
// The response returns the code of the new player
func CreateOnePlayer(c *gin.Context){
	GetCode: code := hexcode.Generate(16)
	_, err := database.Pool().Exec(
		"INSERT INTO players(code) "+
		"VALUES($1)", 
		code,
	)

	if err == nil {
		c.String(200, code)
	} else {
		// On duplicate key error, try a different code
		if database.GetErrorCode(err) == "23505" {
			goto GetCode
		} else {
			c.Status(500)
		}
	}
}

// Fetches one player's details from the database by a player code
// The response returns the details as a JSON
func FetchOnePlayerByCode(c *gin.Context){
	playerCode := c.Param("playerCode")
	row := database.Pool().QueryRow(
		"SELECT id, code "+
		"FROM players "+
		"WHERE code=$1", 
		playerCode,
	)

	player := struct {
		Id int
		Code string
	}{}

	err := row.Scan(&player.Id, &player.Code)

	if err == nil {
		c.JSON(200, player)
	} else {
		c.Status(404)
	}
}