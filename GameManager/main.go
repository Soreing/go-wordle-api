package main

import (
	"GameManager/helper/database"
	"time"
	"github.com/joho/godotenv"
)

// Gets the current game from the database
func getCurrentGameExpiry() (int64, error) {
	row := database.Pool().QueryRow("SELECT \"expiresAt\" FROM games ORDER BY id DESC LIMIT 1")

	var expiry int64
	err := row.Scan(&expiry)
	return expiry, err
}

// Gets the current game from the database
func createNewGameInDatabase() (int64, error) {
	var wordId int
	row := database.Pool().QueryRow("SELECT id FROM words WHERE solution=true ORDER BY RANDOM() LIMIT 1")
	err1 := row.Scan(&wordId)

	if err1 != nil {
		return 0, err1
	}

	beg := time.Now().UnixNano()/1000000
	end := beg + 7200 * 1000
	_, err2 := database.Pool().Exec("INSERT INTO games(\"wordId\", \"createdAt\", \"expiresAt\") VALUES($1, $2, $3)", wordId, beg, end)
	return int64(end), err2
}

func main() {
	if err := godotenv.Load(); err != nil {
    	panic("Error loading .env file")
  	}
	
	var err error
	var expiry int64

	database.Connect()
	defer database.Disconnect()

	expiry, err = getCurrentGameExpiry()
	if err != nil { panic(err) }

	for {
		now := time.Now().UnixNano()/1000000
		dif := expiry - now

		if dif <= 0 {
			expiry, err = createNewGameInDatabase()
		} else {
			time.Sleep(time.Duration(dif/1000 + 1) * time.Second)
		}
	}

}

