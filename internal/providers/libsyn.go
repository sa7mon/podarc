package providers

import (
	"encoding/xml"
	"fmt"
	"github.com/sa7mon/podarc/internal/interfaces"
	"log"
	"net/http"
	"time"
)

type LibsynPodcast struct {
	Title 			string 				`xml:"channel>title"`
	ShowDescription string            	`xml:"channel>summary"` // itunes:summary
	Episodes        []LibsynEpisode 	`xml:"channel>item"`
}

type LibsynEpisode struct {
	Title       string          `xml:"title"`
	Guid        string          `xml:"guid"`
	Image       LibsynImage     `xml:"image"`
	Description string          `xml:"description"`
	Published   string          `xml:"pubDate"`
	Enclosure   LibsynEnclosure `xml:"enclosure"`
}

type LibsynEnclosure struct {
	Url 	string `xml:"url,attr"`
}

type LibsynImage struct {
	ImageUrl string `xml:"href,attr"`
}

func (l LibsynPodcast) NumEpisodes() int {
	return len(l.Episodes)
}

func (l LibsynPodcast) GetEpisodes() []interfaces.PodcastEpisode {
	// TODO: Might be more efficient to store these values rather than do a for loop every time the getter is called
	// Golang doesn't allow you to directly return a slice of a type as a slice of an interface
	// https://golang.org/doc/faq#convert_slice_of_interface
	intEpisodes := make([]interfaces.PodcastEpisode, len(l.Episodes))
	for i, elem := range l.Episodes {
		intEpisodes[i] = elem
	}
	return intEpisodes
}

func (l LibsynPodcast) GetTitle() string {
	return l.Title
}

func (l LibsynPodcast) GetDescription() string {
	return l.ShowDescription
}

func (l LibsynPodcast) GetPublisher() string {
	return "Libsyn"
}

func (l LibsynEpisode) GetTitle() string {
	return l.Title
}

func (l LibsynEpisode) GetDescription() string {
	return l.Description
}

func (l LibsynEpisode) GetUrl() string {
	return l.Enclosure.Url
}

func (l LibsynEpisode) GetPublishedDate() string {
	return l.Published
}

func (l LibsynEpisode) GetParsedPublishedDate() (time.Time, error) {
	layout := "Mon, 02 Jan 2006 15:04:05 -0700"
	t, err := time.Parse(layout, l.GetPublishedDate())
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func (l LibsynEpisode) GetImageUrl() string {
	return l.Image.ImageUrl
}

func (l LibsynEpisode) ToString() string {
	return fmt.Sprintf("Title: %s | Description: %s | Url: %s | PublishedDate: " +
		"%s | ImageUrl: %s", l.GetTitle(), l.GetDescription(), l.GetUrl(), l.GetPublishedDate(),
		l.GetImageUrl())
}

func GetLibsynProPodcastFeed(rssUrl string) *LibsynPodcast {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest("GET", rssUrl, nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 200 {
		log.Fatal("Bad status code while getting podcast - " + resp.Status)
	}

	podcast := &LibsynPodcast{}

	xmlDecoder := xml.NewDecoder(resp.Body)
	err = xmlDecoder.Decode(podcast)
	if err != nil {
		log.Fatal(err)
	}
	return podcast
}