package providers

import (
	"encoding/xml"
	"fmt"
	"github.com/sa7mon/podarc/internal/interfaces"
	"log"
	"net/http"
	"time"
)

/*************************
*
*		Structs
*
**************************/

type StitcherPodcast struct {
	Feed            StitcherFeed      `xml:"feed"`
	ShowDescription string            `xml:"show_description,attr"`
	Episodes        []StitcherEpisode `xml:"episodes>episode"`
}

type StitcherFeed struct { 		  // TODO: Possibly move these properties to StitcherPodcast and delete this struct
	Name          string          `xml:"name"`
	Description   string          `xml:"description"`
	LatestEpisode StitcherEpisode `xml:"episode"`
	Premium       bool            `xml:"premium,attr"`
	EpisodeCount  int             `xml:"episodeCount,attr"`
}

type StitcherEpisode struct {
	Id          string `xml:"id,attr"`
	Image       string `xml:"episodeImage,attr"`
	Published   string `xml:"published,attr"`
	Title       string `xml:"title"`
	Description string `xml:"description"`
	URL         string `xml:"url,attr"`
}

/*************************
*
*    Getters / Setters
*
**************************/


func (s StitcherPodcast) NumEpisodes() int {
	return s.Feed.EpisodeCount
}

func (s StitcherPodcast) GetEpisodes() []interfaces.PodcastEpisode {
	// TODO: Might be more efficient to store these values rather than do a for loop every time the getter is called
	// Golang doesn't allow you to directly return a slice of a type as a slice of an interface
	// https://golang.org/doc/faq#convert_slice_of_interface
	intEpisodes := make([]interfaces.PodcastEpisode, len(s.Episodes))
	for i, elem := range s.Episodes {
		intEpisodes[i] = elem
	}
	return intEpisodes
}

func (s StitcherPodcast) GetTitle() string {
	return s.Feed.Name
}

func (s StitcherPodcast) GetDescription() string {
	return s.ShowDescription
}

func (s StitcherPodcast) GetPublisher() string {
	return "Stitcher"
}

func (s StitcherEpisode) GetTitle() string {
	return s.Title
}

func (s StitcherEpisode) GetDescription() string {
	return s.Description
}

func (s StitcherEpisode) GetUrl() string {
	return s.URL
}

func (s StitcherEpisode) GetPublishedDate() string {
	return s.Published
}

func (s StitcherEpisode) GetParsedPublishedDate() (time.Time, error) {
	layout := "2006-01-02 15:04:05"
	t, err := time.Parse(layout, s.GetPublishedDate())
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func (s StitcherEpisode) GetImageUrl() string {
	return s.Image
}

func (s StitcherEpisode) ToString() string {
	return fmt.Sprintf("Title: %s | Description: %s | URL: %s | PublishedDate: " +
		"%s | ImageURL: %s", s.GetTitle(), s.GetDescription(), s.GetUrl(), s.GetPublishedDate(),
		s.GetImageUrl())
}

func GetStitcherPodcastFeed(feedID string, sess string) *StitcherPodcast {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://app.stitcher.com/Service/GetFeedDetailsWithEpisodes.php?" +
		"mode=webApp&fid=%s&max_epi=5000&sess=%s", feedID, sess), nil)
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

	podcast := &StitcherPodcast{}

	xmlDecoder := xml.NewDecoder(resp.Body)
	err = xmlDecoder.Decode(podcast)
	if err != nil {
		log.Fatal(err)
	}
	return podcast
}