package redis

import (
	"AndroidSide/internal"
	"AndroidSide/sqlite"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func ListenContacts() {
	log.Println("Esperando Mensagens")
	lastId := "0"

	for {
		result, err := client.XRead(ctx, &redis.XReadArgs{
			Streams: []string{"wa:contacts", lastId},
			Count:   1,
			Block:   0,
		}).Result()
		if err != nil {
			log.Println("Erro no XRead, Reconectando..", err)
			time.Sleep(2 * time.Second)
			continue
		}
		for _, stream := range result {
			for _, message := range stream.Messages {
				lastId = message.ID
				payload, ok := message.Values["payload"].(string)
				if !ok {
					log.Println("Erro ao extrair payload")
					continue
				}
				var contact sqlite.Contacts
				if err := json.Unmarshal([]byte(payload), &contact); err != nil {
					log.Println("Erro ao deserializar contato", err)
					continue
				}
				if err := internal.InsertContact(contact); err != nil {
					log.Println("Erro ao inserir contato", err)
				}
			}
		}
	}
}
