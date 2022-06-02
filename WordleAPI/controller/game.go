package controller

import (
	"WordleAPI/helper/database"
	"WordleAPI/helper/wordle"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Fetches the current game from the database
// The response returns the details as a JSON
func FetchCurrentGame(c *gin.Context){
	row := database.Pool().QueryRow(
		"SELECT \"id\", \"expiresAt\" "+
		"FROM games "+
		"ORDER BY id DESC LIMIT 1",
	)

	game := struct{
		Id int
		ExpiresAt uint64
	} {}

	err := row.Scan(&game.Id, &game.ExpiresAt)

	if err == nil {
		c.JSON(http.StatusOK, game)
	} else {
		c.Status(500)
	}
}

// Fetches one player's status from a specific game
// The response returns the colored letters of all guesses as a JSON
func FetchPlayersGame(c *gin.Context){
	gameId := c.Param("gameId")
	playerCode := c.Param("playerCode")
	
	// Fetch the solution of the given game
	game := database.Pool().QueryRow(
		"SELECT word "+
		"FROM games "+
		"INNER JOIN words on games.\"wordId\"=words.id "+
		"WHERE games.id=$1", 
		gameId,
	)
	
	var solution string
	if err := game.Scan(&solution); err != nil{
		c.String(404, "Game not found")
	} else {
		// Fetch the guesses of the user in the given game
		guesses, err := database.Pool().Query(
			"SELECT \"guessNum\", words.word " +  
			"FROM guesses " +
			"INNER JOIN words ON guesses.\"wordId\"=words.id " +
			"INNER JOIN players ON guesses.\"playerId\"=players.id " +
			"WHERE \"gameId\"=$1 AND code=$2 " +
			"ORDER BY \"guessNum\" DESC", 
			gameId, playerCode,
		)

		status := struct {
			Solved bool
			GuessCount int
			Guesses[]string
		}{}

		// Srotre the first row (The last guess)
		guesses.Next()
		var guessNum int
		var guessWord string
		err = guesses.Scan(&guessNum, &guessWord)

		// If there was no first row, initialize an empty structure
		if err != nil {
			status.Solved = false
			status.GuessCount = 0
			status.Guesses = make([]string, 0)
		// Else initialize the structure and process guesses
		} else {
			status.GuessCount = guessNum
			status.Guesses = make([]string, guessNum)

			for i := guessNum-1; i >= 0; i-- {
				status.Guesses[i] = wordle.ColorWord(guessWord, solution)
				if guessWord == solution {
					status.Solved = true
				}

				guesses.Next()
				err = guesses.Scan(&guessNum, &guessWord)
			}
		}

		// Respond with the game status
		c.JSON(200, status)
	}
}