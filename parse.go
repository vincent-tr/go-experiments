package main

import (
	"io"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/goodsign/monday"
)

type item struct {
	date   time.Time
	active bool
}

var data map[string]*item = make(map[string]*item)

func parse(body io.Reader) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		log.Fatal((err))
	}

	doc.Find("td[style='line-height:18px; color:#000000;padding-top:10px; padding-bottom:4px; font-family:Arial, Helvetica, sans-serif; font-size:18px;']").Each(func(i int, s *goquery.Selection) {
		parseLine(s.Text())
	})
}

func parseLine(line string) {
	r := regexp.MustCompile(`(.*)\(Maison\)\s+\((.*)\)`)
	m := r.FindStringSubmatch(line)
	if len(m) != 3 {
		log.Fatalf("Cannot parse line: '%s'\n", line)
	}

	active := true
	label := m[1]
	label = strings.TrimSpace(label)
	const prefix = "RESTAURATION"
	if strings.HasPrefix(label, prefix) {
		active = false
		label = strings.TrimSpace(label[len(prefix):])
	}

	sdate := m[2]

	date, err := monday.ParseInLocation("15:04:05 02/Jan/06", sdate, time.Local, monday.LocaleFrFR)
	if err != nil {
		log.Fatal((err))
	}

	//log.Printf("%s\n", line)
	//log.Printf("%s %s %t\n", label, date, active)

	it := data[label]
	if it != nil {
		if it.date.Before(date) {
			it.date = date
			it.active = active
		}
	} else {
		it = &item{
			active: active,
			date:   date,
		}

		data[label] = it
	}
}
