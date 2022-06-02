package main

import (
	"WordleAPI/helper/database"
	"WordleAPI/controller"
	"os"
	"fmt"
	"time"
	"math/rand"
	"github.com/joho/godotenv"
	"github.com/gin-gonic/gin"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	if err := godotenv.Load(); err != nil {
    	panic("Error loading .env file")
  	}

	database.Connect()
	defer database.Disconnect()
	
	r := gin.Default()
	r.GET("/game/current", controller.FetchCurrentGame)
	r.GET("/game/:gameId/:playerCode", controller.FetchPlayersGame)
	r.GET("/player/:playerCode", controller.FetchOnePlayerByCode)
	r.POST("/player", controller.CreateOnePlayer)
	r.POST("/guess", controller.MakeOneGuessInAGame)
	
	if os.Getenv("SSL") != "disable" {
		r.RunTLS(fmt.Sprintf(":%s", os.Getenv("HTTPS_PORT")), os.Getenv("CERTFILE"), os.Getenv("KEYFILE"))
	} else {
		r.Run(fmt.Sprintf(":%s", os.Getenv("HTTP_PORT")))
	}
}

