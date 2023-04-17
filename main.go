package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"openai-request/connection"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Connect to database
	connection.DatabaseConnect()

	// Create new Echo instance
	e := echo.New()

	// Serve static files from "/public" directory
	e.Static("/public", "public")

	// Routing
	e.GET("/openai", openaiChat)

	port := os.Getenv("PORT")

	// Start server
	println("Server running on port " + port)
	e.Logger.Fatal(e.Start(":" + port))
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func openaiChat(c echo.Context) error {
	role := c.QueryParam("role")
	content := c.QueryParam("content")

	messages := []Message{
		{
			Role:    role,
			Content: content,
		},
	}

	requestBody := map[string]interface{}{
		"model":    "gpt-3.5-turbo",
		"messages": messages,
	}

	postBody, err := json.Marshal(requestBody)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(postBody)) // Print the JSON request body

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewReader(postBody))
	if err != nil {
		log.Fatal(err)
	}

	// Set the authorization header
	req.Header.Set("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))

	// Set the Content-Type header
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var jsonData interface{}
	err = json.Unmarshal(responseData, &jsonData)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(jsonData)
	return c.JSON(http.StatusOK, jsonData)
}
