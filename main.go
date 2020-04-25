package main

import (
	"flag"
	"fmt"
	"github.com/sa7mon/podarc/internal/providers"
	"github.com/sa7mon/podarc/internal/utils"
	"log"
	"os"
)

func main() {
	feedUrl := flag.String("feedUrl", "", "URL of podcast feed to archive (Required)")
	flag.Parse()

	if *feedUrl == "" || !utils.IsValidUrl(*feedUrl){
		fmt.Printf("Error - Invalid feedUrl: '%s'\n", *feedUrl)
		flag.PrintDefaults()
		os.Exit(1)
	}

	creds := utils.ReadCredentials("creds.json")
	creds = creds

	//feedUrl := "http://mates.nerdistind.libsynpro.com/rss"
	//feedUrl := "https://app.stitcher.com/browse/feed/467097/details"
	fetchedPodcast := providers.FetchPodcastFromUrl(*feedUrl, creds)

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