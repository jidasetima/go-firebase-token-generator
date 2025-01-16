package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2/google"
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

	// Read the file contents into memory
	fileContent, err := os.ReadFile(filenamePayload.Filename)
	if err != nil {
		log.Print("Error reading credentials file: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "File contents not readable"})
		return
	}

	// Retrieve the token using Google OAuth2 credentials
	credentials, err := google.CredentialsFromJSON(context.Background(), fileContent, "https://www.googleapis.com/auth/firebase.messaging")
	if err != nil {
		log.Print("Error loading credentials: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "File contents not readable"})
		return
	}

	// Get token from source
	token, err := credentials.TokenSource.Token()
	if err != nil {
		log.Print("Error retrieving OAuth2 token: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not able to retrieve OAuth2 token"})
		return
	}

	var tokenPayload tokenPayload
	tokenPayload.Token = token.AccessToken

	c.JSON(http.StatusOK, tokenPayload)
}
