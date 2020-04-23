package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"podarc/interfaces"
	"podarc/providers"
	"time"
)

type credentials struct {
	SessionToken string	`json:"session_token"`
}

func main() {
	creds := readCredentials("../../creds.json")

	fetchedPodcast := fetchPodcastFromUrl("http://mates.nerdistind.libsynpro.com/rss", creds)

	for _, episode := range fetchedPodcast.GetEpisodes() {
		log.Println(episode.GetTitle())
	}

	//log.Println("Getting podcast data...")
	//officeLadies := getStitcherPodcastFeed("467097", creds.SessionToken)
	//log.Println(officeLadies.ShowDescription)
	//
	//for _, element := range officeLadies.Episodes {
	//	log.Println(element.Published)
	//}
}

func fetchPodcastFromUrl(feedUrl string, creds credentials) interfaces.Podcast {

}

func readCredentials(file string) credentials {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Print(err)
	}
	var creds credentials
	err = json.Unmarshal(data, &creds)
	if err != nil {
		fmt.Println("Error reading creds file: ", err)
	}
	return creds
}

func getStitcherPodcastFeed(feedId string, sess string) *providers.StitcherPodcast {
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

	podcast := &providers.StitcherPodcast{}

	xmlDecoder := xml.NewDecoder(resp.Body)
	err = xmlDecoder.Decode(podcast)
	if err != nil {
		log.Fatal(err)
	}
	return podcast
}