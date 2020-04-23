package main

import (
	"encoding/xml"
	"fmt"
	"github.com/sa7mon/podarc/internal/interfaces"
	"github.com/sa7mon/podarc/internal/providers"
	"github.com/sa7mon/podarc/internal/utils"
	"log"
	"net/http"
	"regexp"
	"time"
)


func main() {
	creds := utils.ReadCredentials("creds.json")
	creds = creds

	feedUrl := "http://mates.nerdistind.libsynpro.com/rss"
	//feedUrl := "https://app.stitcher.com/browse/feed/467097/details"
	fetchedPodcast := fetchPodcastFromUrl(feedUrl, creds)

	for _, episode := range fetchedPodcast.GetEpisodes() {
		log.Println(episode.GetUrl())
	}

	//fmt.Println("Download Started")
	//fmt.Println(fetchedPodcast.GetEpisodes()[0].GetImageUrl())
	//
	//fileUrl := fetchedPodcast.GetEpisodes()[0].GetUrl()
	//err := utils.DownloadFile("podcast.mp3", fileUrl)
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Println("Download Finished")
}

func fetchPodcastFromUrl(feedUrl string, creds utils.Credentials) interfaces.Podcast {
	stitcherR := regexp.MustCompile(`https://app\.stitcher\.com/browse/feed/(?P<feedId>\d+)`)
	libsynR := regexp.MustCompile(`\S+\.libsynpro.com/rss`)
	libSynMatches := libsynR.MatchString(feedUrl)
	stitcherMatches := stitcherR.FindStringSubmatch(feedUrl)

	if len(stitcherMatches) > 0 {
		fmt.Println("Stitcher feed detected")
		fmt.Println("Feed ID: " + stitcherMatches[1]) // Capture group names available via: stitcherR.SubexpNames()

		stitcherPod := getStitcherPodcastFeed(stitcherMatches[1], creds.SessionToken)
		return stitcherPod
	} else if libSynMatches {
		fmt.Println("Libsyn Pro feed detected")
		libsynPod := getLibsynProPodcastFeed(feedUrl)
		return libsynPod
	} else {

	}
	panic("Unknown URL!")
}

func getLibsynProPodcastFeed(rssUrl string) *providers.LibsynPodcast {
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

	podcast := &providers.LibsynPodcast{}

	xmlDecoder := xml.NewDecoder(resp.Body)
	err = xmlDecoder.Decode(podcast)
	if err != nil {
		log.Fatal(err)
	}
	return podcast
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