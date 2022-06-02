package controller

import (
	"WordleAPI/helper/database"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
	"strings"
)

// Makes a new guess from a plaer in a game
// Inserts the new guess in the database
func MakeOneGuessInAGame(c *gin.Context){
	word   := strings.ToLower(c.PostForm("Word"))
	player := c.PostForm("Player")
	game   := c.PostForm("GameId")
	gameId, err := strconv.Atoi(game)

	if word == "" || player == "" || game == "" || err != nil {
		c.String(400, "Invalid parameters")
	} else {
		wordRec := database.Pool().QueryRow(
			"SELECT id "+
			"FROM words "+
			"WHERE word=$1", 
			word,
		)

		playerRec := database.Pool().QueryRow(
			"SELECT id "+
			"FROM players "+
			"WHERE code=$1", 
			player,
		)

		gameRec := database.Pool().QueryRow(
			"SELECT words.word, \"expiresAt\" "+
			"FROM games "+
			"INNER JOIN words ON games.\"wordId\"=words.id "+
			"WHERE games.id=$1",
			gameId,
		)
		
		var wordId int
		var playerId int
		var solution string
		var expiry int64

		now := time.Now().UnixNano() / 1000000
		werr := wordRec.Scan(&wordId)
		perr := playerRec.Scan(&playerId)
		gerr := gameRec.Scan(&solution, &expiry)

		if werr != nil {
			c.String(400, "Word not valid")
		} else if perr != nil {
			c.String(404, "Player not found")
		} else if gerr != nil {
			c.String(404, "Game not found")
		} else if expiry < now {
			c.String(400, "Game already expired")
		} else {
			lastGuess := database.Pool().QueryRow(
				"SELECT \"guessNum\", words.word " +  
				"FROM guesses " +
				"INNER JOIN words ON guesses.\"wordId\"=words.id " +
				"WHERE \"gameId\"=$1 AND \"playerId\"=$2 " +
				"ORDER BY \"guessNum\" DESC LIMIT 1", 
				gameId, playerId,
			)
			
			var guessCount int
			var guessWord string

			if err = lastGuess.Scan(&guessCount, &guessWord); err != nil {
				guessCount = 0
				guessWord = ""
			}


			if guessCount >= 6 {
				c.String(400, "Too many guesses already")
			} else if guessWord == solution {
				c.String(400, "The game is already solved")
			} else {
				_, err = database.Pool().Exec(
					"INSERT INTO guesses "+
					"VALUES($1, $2, $3, $4)", 
					gameId, playerId, wordId, guessCount+1,
				)
				
				if err != nil {
					c.Status(500)
				} else {
					c.Status(200)
				}
			}
		}
	}
}