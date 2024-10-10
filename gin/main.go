package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		time.Sleep(20 * time.Second) // or some db operation
		log.Print("Processing")
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run(":10087")
}
