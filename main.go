package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"time"
)

type StitcherPodcast struct {
	Feed 			StitcherFeed 	  `xml:"feed"`
	ShowDescription string 	   		  `xml:"show_description,attr"`
	Episodes 		[]StitcherEpisode `xml:"episodes>episode"`
}

type StitcherEpisodes struct {
	Episodes []StitcherEpisode
}
type StitcherFeed struct {
	Name 		  string 		  `xml:"name"`
	Description   string 		  `xml:"description"`
	LatestEpisode StitcherEpisode `xml:"episode"`
	Premium 	  bool 			  `xml:"premium,attr"`
	EpisodeCount  int 			  `xml:"episodeCount,attr"`
}

type StitcherEpisode struct {
	Id          string    `xml:"id,attr"`
	Image       string    `xml:"episodeImage,attr"`
	Published   time.Time `xml:"published,attr"`
	Title       string    `xml:"title"`
	Description string    `xml:"description"`
	Url 		string 	  `xml:"url, attr"`
}

func main() {
	officeLadies := getStitcherPodcastFeed("","")
	log.Println(officeLadies.ShowDescription)

	for _, element := range officeLadies.Episodes {
		log.Println(element.Title)
	}
}

func getStitcherPodcastFeed(feedId string, sess string) *StitcherPodcast {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://app.stitcher.com/Service/GetFeedDetailsWithEpisodes.php?" +
								"mode=webApp&fid=%s&max_epi=5000&sess=%s", feedId, sess), nil)
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