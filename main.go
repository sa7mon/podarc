package main

import (
	"fmt"
	"github.com/sa7mon/podarc/internal/interfaces"
	"github.com/sa7mon/podarc/internal/providers"
	"github.com/sa7mon/podarc/internal/utils"
	"log"
	"regexp"
)


func main() {
	creds := utils.ReadCredentials("creds.json")
	creds = creds

	//feedUrl := "http://mates.nerdistind.libsynpro.com/rss"
	feedUrl := "https://app.stitcher.com/browse/feed/467097/details"
	fetchedPodcast := fetchPodcastFromUrl(feedUrl, creds)

	for _, episode := range fetchedPodcast.GetEpisodes() {
		log.Println(episode.ToString())
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

		stitcherPod := providers.GetStitcherPodcastFeed(stitcherMatches[1], creds.SessionToken)
		return stitcherPod
	} else if libSynMatches {
		fmt.Println("Libsyn Pro feed detected")
		libsynPod := providers.GetLibsynProPodcastFeed(feedUrl)
		return libsynPod
	}
	panic("Unknown URL!")
}