package rss

import (
	"time"
)

type Date string

const wordpressDateFormat = "Mon, 02 Jan 2006 15:04 EDT"

// Channel struct for RSS
type Channel struct {
	Title         string `xml:"title"`
	Link          string `xml:"link"`
	Description   string `xml:"description"`
	Language      string `xml:"language"`
	LastBuildDate Date   `xml:"lastBuildDate"`
	Item          []Item `xml:"item"`
}

// ItemEnclosure struct for each Item Enclosure
type ItemEnclosure struct {
	URL  string `xml:"url,attr"`
	Type string `xml:"type,attr"`
}

// Item struct for each Item in the Channel
type Item struct {
	Title       string          `xml:"title"`
	Link        string          `xml:"link"`
	Comments    string          `xml:"comments"`
	PubDate     Date            `xml:"pubDate"`
	GUID        string          `xml:"guid"`
	Category    []string        `xml:"category"`
	Enclosure   []ItemEnclosure `xml:"enclosure"`
	Description string          `xml:"description"`
	Author      string          `xml:"author"`
	Content     string          `xml:"content"`
	FullText    string          `xml:"full-text"`
}

type RssFeed struct {
	Channel Channel `xml:"channel"`
}

// Parse (Date function) and returns Time, error
func (d Date) Parse() (time.Time, error) {
	t, err := time.Parse(wordpressDateFormat, string(d))
	if err != nil {
		t, err = time.Parse(time.RFC822, string(d)) // RSS 2.0 spec
		if err != nil {
			t, err = time.Parse(time.RFC3339, string(d)) // AtomS
			if err != nil {
				return time.Time{}, err
			}
		}
	}
	return t, nil
}
