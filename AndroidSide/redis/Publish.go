package redis

import (
	"AndroidSide/sqlite"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var client *redis.Client

func InitRedis() {
	client = redis.NewClient(&redis.Options{
		Addr: "10.0.3.2:6379",
	})
}

func PublishToRedis(contact *sqlite.Contacts) error {

	data, err := json.Marshal(contact)
	if err != nil {
		return err
	}
	return client.XAdd(ctx, &redis.XAddArgs{
		Stream: "wa:messages",
		ID:     "*",
		Values: map[string]any{
			"payload": string(data),
		},
	}).Err()
}
func PublishWithRetry(contact *sqlite.Contacts, maxAttempts int) {
	backoff := time.Second
	for i := range maxAttempts {
		if err := PublishToRedis(contact); err != nil {
			log.Printf("Tentativa %d/%d falhou: %v - aguardando %s\n", i+1, maxAttempts, err, backoff)
			time.Sleep(backoff)
			backoff = min(backoff*2, 30*time.Second)
			continue
		}
		return
	}
	log.Printf("Falha crítica ao publicar %s após %d tentativas\n", contact.Jid, maxAttempts)
}
