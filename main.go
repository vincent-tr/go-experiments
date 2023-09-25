package main

import (
	"log"
	"math"
	"mime/quotedprintable"
	"net/mail"
	"os"
	"strings"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

func main() {
	log.Println("Connecting to server...")

	// Connect to server
	c, err := client.DialTLS("imap.gmail.com:993", nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected")

	// Don't forget to logout
	defer c.Logout()

	// Login
	if err := c.Login(os.Getenv("USER"), os.Getenv("PASS")); err != nil {
		log.Fatal(err)
	}
	log.Println("Logged in")

	// Select INBOX
	mbox, err := c.Select("INBOX", true)
	if err != nil {
		log.Fatal(err)
	}

	section := &imap.BodySectionName{}
	//items := []imap.FetchItem{section.FetchItem()}

	// Get the last 4 messages
	from := uint32(1)
	to := mbox.Messages
	if mbox.Messages > 300 {
		// We're using unsigned integers here, only subtract if the result is > 0
		from = mbox.Messages - 300
	}
	seqset := new(imap.SeqSet)
	seqset.AddRange(from, to)

	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope, section.FetchItem()}, messages)
	}()

	for msg := range messages {
		if !strings.Contains(msg.Envelope.Subject, "Bentel") {
			continue
		}

		log.Printf("* %s\n", msg.Envelope.Date.Local().String())

		r := msg.GetBody(section)
		if r == nil {
			log.Fatal("Server didn't returned message body")
		}

		m, err := mail.ReadMessage(r)
		if err != nil {
			log.Fatal(err)
		}

		hdate, err := time.Parse("Mon, 02 Jan 2006 15:04:05 -0700 (MST)", strings.TrimSpace(strings.Split(m.Header["Received"][0], ";")[1]))
		if err != nil {
			log.Fatal(err)
		}

		diff := hdate.Sub(msg.Envelope.Date).Seconds()

		if math.Abs(diff) > 60 {
			log.Printf("Header received date %s\n", hdate.Local())
			log.Printf("Header diff %fsecs\n", diff)
			log.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
		}

		r2 := quotedprintable.NewReader(m.Body)
		/*
			body, err := io.ReadAll(r2)
			if err != nil {
				log.Fatal(err)
			}

			log.Println(string(body))
		*/
		parse(r2)
	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}

	log.Println("Done!")

	for label, item := range data {
		var a string
		if item.active {
			a = "ON"
		} else {
			a = "OFF"
		}
		log.Printf("%s -> %s (%s)\n", label, a, item.date.String())
	}
}
