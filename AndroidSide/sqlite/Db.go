package sqlite

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

func QueryMaxId() string {
	id := `SELECT MAX(_id) FROM message`
	return id
}
func QueryLastMessage() string {
	sqlStmt := `
				select
					m.text_data as text,
					j.raw_string as jid,
					w.wa_name as Nome,
					m._id as id
				from message m
				join chat c
					on m.chat_row_id = c._id
				join jid_map jm
					on c.jid_row_id = jm.lid_row_id
				join jid j
					on jm.jid_row_id = j._id
				join wadb.wa_contacts w
					on w.jid = j.raw_string
				WHERE m._id > ?
	`
	return sqlStmt
}
func Db_Connect() *sql.DB {
	dsn := "/data/data/com.whatsapp/databases/msgstore.db?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)"

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		log.Fatal("erro ao abrir msgstore.db:", err)
	}

	db.SetMaxOpenConns(1)

	_, err = db.Exec(`ATTACH DATABASE '/data/data/com.whatsapp/databases/wa.db' AS wadb`)
	if err != nil {
		log.Fatal("erro ao fazer attach do wa.db:", err)
	}

	log.Println("msgstore.db conectado e wa.db attached")
	return db
}
