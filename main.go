package main

import (
	firebase "firebase.google.com/go/v4"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"os"
)

type tokenPayload struct {
	Token string `json:"token"`
}

type filenamePayload struct {
	Filename string `json:"filename"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file", err)
	}

	serverUrl, serverUrlIsSet := os.LookupEnv("SERVER_URL")

	if serverUrlIsSet != true {
		log.Fatalf("Error loading SERVER_URL environment variable")
	}

	router := gin.Default()
	router.POST("/token", getToken)

	router.Run(serverUrl)
}

func getToken(c *gin.Context) {
	var filenamePayload filenamePayload
	if err := c.BindJSON(&filenamePayload); err != nil {
		log.Print("Bad Request. Filename not provided")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Filename required"})
		return
	}

	credentials := option.WithCredentialsFile(filenamePayload.Filename)
	application, err := firebase.NewApp(c, nil, credentials)
	if err != nil {
		log.Fatalf("error initializing app: $v\n", err)
	} else {
		log.Print("successful connection to firebase")
	}

	client, err := application.Auth(c)

	if err != nil {
		log.Print("error getting auth: $v\n", err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	token, err := client.CustomToken(c, filenamePayload.Filename)

	if err != nil {
		log.Print("error getting token: $v\n", err)
		c.JSON(http.StatusBadGateway, gin.H{"Firebase token minting failed with error": err.Error()})
		return
	}
	log.Printf("Got custom token: $v\n", token)

	var tokenPayload tokenPayload
	tokenPayload.Token = token

	c.JSON(http.StatusOK, tokenPayload)
}
