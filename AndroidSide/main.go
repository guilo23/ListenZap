package main

import (
	"AndroidSide/redis"
	"AndroidSide/sqlite"
	"fmt"
)

func main() {
	chanel := make(chan *sqlite.Contacts)
	go sqlite.WatchMessagesDb(chanel)

	redis.InitRedis()

	fmt.Println("Monitorando....")

	go redis.ListenContacts()

	for contato := range chanel {
		fmt.Printf("Mensagem Detectada: %s; DE: %s; Número: %s ",
			contato.Message, contato.Nome, contato.Number)
		redis.PublishWithRetry(contato, 5)
	}
}
