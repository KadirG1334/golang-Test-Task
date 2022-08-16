package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
)

type PostBody struct {
	Sender    string    `form:"sender" binding:"required"`
	Receiver  string    `form:"receiver" binding:"required"`
	Message   string    `form:"message" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
}

func main() {
	K := gin.Default()

	conn, err := amqp.Dial("amqp://user:password@localhost:5672/")
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalln(err)
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		"Message",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalln(err)
	}

	r.POST("/message", func(c *gin.Context) {
		var postBody PostBody
		if err := c.ShouldBindJSON(&postBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		postBody.CreatedAt = time.Now()

		jsonBody, err := json.Marshal(postBody)

		err = ch.Publish(
			"",
			"Message",
			false,
			false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        jsonBody,
			},
		)

		fmt.Println(string(jsonBody))

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "database error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	K.Run()
}
