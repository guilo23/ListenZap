package sqlite

import (
	"log"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

type Contacts struct {
	Jid     string `json:"jid,omitempty"`
	Message string `json:"message,omitempty"`
	Nome    string `json:"name"`
	Number  string `json:"number"`
}

func ExtractNumber(jid string) string {
	if strings.Contains(jid, "@s.whatsapp.net") {
		return strings.Split(jid, "@")[0]
	}
	return ""
}

func WatchMessagesDb(ch chan<- *Contacts) {
	db := Db_Connect()
	defer db.Close()

	var lastId string
	err := db.QueryRow(QueryMaxId()).Scan(&lastId)
	if err != nil {
		log.Fatal("erro ao buscar maxId inicial:", err)
	}

	log.Println("Iniciando watch a partir do ID:", lastId)

	for {
		rows, err := db.Query(QueryLastMessage(), lastId)
		if err != nil {
			log.Println("erro na query, aguardando...", err)
			time.Sleep(2 * time.Second)
			continue
		}

		for rows.Next() {
			var message, jid, nome, currentId string

			if err := rows.Scan(&message, &jid, &nome, &currentId); err != nil {
				log.Println("erro ao ler linha:", err)
				continue
			}

			lastId = currentId

			ch <- &Contacts{
				Jid:     jid,
				Nome:    nome,
				Number:  ExtractNumber(jid),
				Message: message,
			}
		}

		rows.Close()
		time.Sleep(2 * time.Second)
	}
}
