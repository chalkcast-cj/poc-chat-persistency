package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
	"strings"
)

func main() {
	db, err := sqlx.Connect("mysql", "root:password@/chat")
	if err != nil {
		log.Fatalln(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	subject := "chat:room_id:1"
	consumerGroup := "chat-ingestion"
	consumer := "my-sql-ingestion"

	if err := rdb.XGroupCreateMkStream(context.TODO(), subject, consumerGroup, "0").Err(); err != nil {
		// a hack
		if !strings.Contains(err.Error(), "BUSYGROUP Consumer Group name already exists") {
			log.Fatalln(err)
		}
	}
	for {
		rdbResult, err := rdb.XReadGroup(context.TODO(), &redis.XReadGroupArgs{
			Group:    consumerGroup,
			Consumer: consumer,
			Streams:  []string{subject, ">"},
			NoAck:    true,
			Count:    1,
		}).Result()
		if err != nil {
			log.Fatalln(err)
		}
		message := rdbResult[0].Messages[0]
		values := message.Values
		query := `INSERT INTO messages (message_id, user_id, room_id, message) VALUES (?, ?, ?, ?)`
		dbResult, err := db.ExecContext(context.TODO(), query, message.ID, values["user_id"], values["room_id"], values["message"])
		if err != nil {
			log.Fatalln(err)
		}
		affected, err := dbResult.RowsAffected()
		if err != nil {
			log.Fatalln(err)
		}
		if affected == 0 {
			log.Fatalln(err)

		}
		fmt.Println(message)
		rdb.XAck(context.TODO(), subject, consumerGroup, message.ID)
	}
}
