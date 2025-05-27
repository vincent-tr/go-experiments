package main

import (
	"fmt"
	"os"
	"time"

	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
)

type Config struct {
	IMAPServer string
	IMAPUser   string
	IMAPPass   string
	Mailbox    string
	From       string // optional
	Subject    string // optional
	SinceDays  int
}

func main() {

	config := &Config{
		IMAPServer: "imap.gmail.com:993",
		IMAPUser:   readSecret("gmail-user"),
		IMAPPass:   readSecret("gmail-pass"),
		Mailbox:    "Factures",
		From:       "confirmation-commande@amazon.fr",
		SinceDays:  10,
	}

	msgs, err := fetchMails(config)
	if err != nil {
		panic(err)
	}

	for _, msg := range msgs {
		fmt.Printf("Message %d - %s\n", msg.UID, msg.Envelope.Subject)
	}
}

// TODO: context?
func fetchMails(config *Config) ([]*imapclient.FetchMessageBuffer, error) {
	client, err := imapclient.DialTLS(config.IMAPServer, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to IMAP server: %w", err)
	}
	defer client.Close()

	if err := client.Login(config.IMAPUser, config.IMAPPass).Wait(); err != nil {
		return nil, fmt.Errorf("failed to login to IMAP server: %w", err)
	}

	_, err = client.Select(config.Mailbox, &imap.SelectOptions{ReadOnly: true}).Wait()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to IMAP server: %w", err)
	}

	criteria := &imap.SearchCriteria{
		Since:  time.Now().AddDate(0, 0, -config.SinceDays),
		Header: []imap.SearchCriteriaHeaderField{},
	}

	if config.From != "" {
		criteria.Header = append(criteria.Header, imap.SearchCriteriaHeaderField{Key: "FROM", Value: config.From})
	}

	if config.Subject != "" {
		criteria.Header = append(criteria.Header, imap.SearchCriteriaHeaderField{Key: "SUBJECT", Value: config.Subject})
	}

	searchData, err := client.UIDSearch(criteria, nil).Wait()
	if err != nil {
		return nil, fmt.Errorf("failed to search emails: %w", err)
	}

	msgs, err := client.Fetch(imap.UIDSetNum(searchData.AllUIDs()...), &imap.FetchOptions{
		UID:           true,
		Envelope:      true,
		BodyStructure: &imap.FetchItemBodyStructure{Extended: true},
	}).Collect()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch emails: %w", err)
	}

	return msgs, nil
}

func readSecret(name string) string {
	buff, err := os.ReadFile("../secrets/" + name)
	if err != nil {
		panic(err)
	}

	return string(buff)
}
