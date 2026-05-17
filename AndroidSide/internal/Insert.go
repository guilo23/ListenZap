package internal

import (
	"AndroidSide/sqlite"
	"fmt"
	"log"
	"os/exec"
	"regexp"
)

func extractContactId(out string) string {
	re := regexp.MustCompile(`_id=(\d+)`)
	matches := re.FindStringSubmatch(out)
	if len(matches) < 2 {
		fmt.Println("nao foi possivel extrair o Id do contato do output")
		return ""
	}
	return matches[1]
}
func InsertContact(contact sqlite.Contacts) error {

	rawCmd := exec.Command("content", "insert",
		"--uri", "content://com.android.contacts/raw_contacts",
		"--bind", "account_type:s:vnd.sec.contact.phone",
		"--bind", "account_name:s:phone",
	)
	if _, err := rawCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("erro ao criar raw_contact: %w", err)
	}

	queryCmd := exec.Command("content", "query",
		"--uri", "content://com.android.contacts/raw_contacts",
		"--projection", "_id",
		"--where", "account_type='vnd.sec.contact.phone'",
		"--sort", "_id DESC",
	)

	out, err := queryCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("erro ao buscar raw_contact_id: %w", err)
	}
	extractId := extractContactId(string(out))

	if extractId == "" {
		return fmt.Errorf("interrompendo insercao: id do contato nao foi encontrado no output")
	}

	nameCmd := exec.Command("content", "insert",
		"--uri", "content://com.android.contacts/data",
		"--bind", "raw_contact_id:i:"+extractId,
		"--bind", "mimetype:s:vnd.android.cursor.item/name",
		"--bind", "data1:s:"+contact.Nome,
	)
	if out, err := nameCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("erro ao inserir nome: %w - output: %s", err, out)
	}
	phoneCmd := exec.Command("content", "insert",
		"--uri", "content://com.android.contacts/data",
		"--bind", "raw_contact_id:i:"+extractId,
		"--bind", "mimetype:s:vnd.android.cursor.item/phone_v2",
		"--bind", "data1:s:+"+contact.Number,
		"--bind", "data2:i:2",
	)
	if out, err := phoneCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("erro ao inserir numero: %w - output: %s", err, out)
	}
	waCmd := exec.Command("content", "insert",
		"--uri", "content://com.android.contacts/data",
		"--bind", "raw_contact_id:i:"+extractId,
		"--bind", "mimetype:s:vnd.android.cursor.item/vnd.com.whatsapp.profile",
		"--bind", "data1:s:+"+contact.Number+"@s.whatsapp.net",
		"--bind", "data3:s:"+contact.Nome,
	)

	if out, err := waCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("erro ao vincular whatsapp: %w — output: %s", err, out)
	}
	log.Printf("Contato inserido com sucesso: %s (%s)\n", contact.Nome, contact.Number)
	return nil
}
