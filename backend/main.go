package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	openai "github.com/sashabaranov/go-openai"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Question struct {
	Question string `json:question`
}

type QuestionAnswer struct {
	gorm.Model
	Question string
	Answer   string
}

func main() {
	var apiToken, databaseUri string

	apiToken = os.Getenv("OPENAI_TOKEN")
	fmt.Println(apiToken)

	databaseUri, exists := os.LookupEnv("DATABASE_URI")
	if !exists {
		databaseUri = "sqlite://file::memory:?cache=shared"
	}

	client := openai.NewClient(apiToken)

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/ask", func(c *gin.Context) {
		var question Question
		err := c.BindJSON(&question)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})

			return
		}

		response, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model: openai.GPT3Dot5Turbo,
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleUser,
						Content: string(question.Question),
					},
				},
			},
		)

		if err != nil {
			fmt.Printf("ChatCompletion error: %v\n", err)

			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		db, err := openDatabaseConnection(databaseUri)
		if err != nil {
			log.Fatalln(err)
		}

		migrate(db)

		db.Create(&QuestionAnswer{Question: string(question.Question), Answer: response.Choices[0].Message.Content})

		c.JSON(http.StatusOK, gin.H{
			"result": response.Choices[0].Message.Content,
		})
	})

	r.Run()
}

func openDatabaseConnection(databaseUri string) (*gorm.DB, error) {
	if strings.HasPrefix(databaseUri, "postgres://") {
		return gorm.Open(postgres.Open(databaseUri), &gorm.Config{})
	} else if strings.HasPrefix(databaseUri, "sqlite://") {
		return gorm.Open(sqlite.Open(databaseUri[9:]), &gorm.Config{})
	} else {
		return nil, fmt.Errorf("unsupported database type")
	}
}

func migrate(db *gorm.DB) {
	db.AutoMigrate(&QuestionAnswer{})
}
