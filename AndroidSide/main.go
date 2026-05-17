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

	fmt.Println("Monitorando o ZAPZAP....")

	go redis.ListenContacts()

	for contato := range chanel {
		fmt.Printf("Mensagem Detectada : %s do : %s numero: %s ",
			contato.Message, contato.Nome, contato.Number)
		redis.PublishWithRetry(contato, 5)
	}
}
