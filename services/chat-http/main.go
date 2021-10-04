package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"log"
	"net/http"
)

type SayCommand struct {
	UserID  uint   `json:"user_id" db:"user_id" binding:"required"`
	RoomID  uint   `json:"room_id" db:"room_id" binding:"required"`
	Message string `json:"message" db:"message" binding:"required"`
}

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	r := gin.Default()
	r.POST("/say", func(c *gin.Context) {
		var sayCommand SayCommand
		if err := c.ShouldBindJSON(&sayCommand); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": err.Error()})
		}
		if err := rdb.XAdd(c.Request.Context(), &redis.XAddArgs{
			Stream: fmt.Sprintf("chat:room_id:%d", sayCommand.RoomID),
			Values: map[string]interface{}{
				"user_id": sayCommand.UserID,
				"room_id": sayCommand.RoomID,
				"message": sayCommand.Message,
			},
		}).Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	})
	// listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	if err := r.Run(); err != nil {
		log.Fatalln(err)
	}
}
