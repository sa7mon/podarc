package main

import (
	"flag"
	"fmt"
	"github.com/sa7mon/podarc/internal/archiver"
	"github.com/sa7mon/podarc/internal/providers"
	"github.com/sa7mon/podarc/internal/utils"
	"log"
	"os"
)

func main2() {

	//http.HandleFunc("/pod-save-america", func(w http.ResponseWriter, r *http.Request) {
	//	feed, err := providers.GetGenericPodcastFeed("https://feeds.feedburner.com/pod-save-america")
	//	if err != nil {
	//		log.Println(err)
	//	}
	//	pf := feed.ToPodarcFile("https://example.com/pod-save-america")
	//	w.Header().Set("Content-Type", "application/xml")
	//	if err := pf.Podcast.Encode(w); err != nil {
	//		http.Error(w, err.Error(), http.StatusInternalServerError)
	//	}
	//})
	//log.Fatal(http.ListenAndServe(":8000", nil))

	feed, err := providers.GetGenericPodcastFeed("https://feeds.feedburner.com/pod-save-america")
	if err != nil {
		log.Println(err)
	}

	//err = feed.SaveToFile("pod-save-america.xml")
	err = feed.SaveToFile2("pod-save-america2.xml")
	if err != nil {
		log.Println(err)
	}
}

func main() {
	feedURL := flag.String("feedUrl", "", "URL of podcast feed to archive. (Required)")
	destDirectory := flag.String("outputDir", "", "Directory to save the files into. (Required)")
	overwrite := flag.Bool("overwrite", false, "Overwrite episodes already downloaded. Default: false")
	renameFiles := flag.Bool("renameFiles", true, "Rename downloaded files to friendly names.")
	threads := flag.Int("threads", 2, "Number of threads to use when downloading")
	flag.Parse()

	if *feedURL == "" || !utils.IsValidURL(*feedURL) {
		fmt.Printf("Error - Invalid feedUrl: '%s'\n", *feedURL)
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *destDirectory == "" {
		fmt.Printf("Error - Invalid outputDir: '%s'\n", *destDirectory)
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *threads == 0 {
		fmt.Println("Error - threads must be larger than 0")
		flag.PrintDefaults()
		os.Exit(1)
	}

	credentials, err := utils.ReadCredentials("creds.json")
	if err != nil {
		log.Println("Error reading creds file: " + err.Error())
	}

	fetchedPodcast, err := providers.FetchPodcastFromURL(*feedURL, credentials)
	if err != nil {
		log.Println("Error fetching podcast from URL - " + err.Error())
		os.Exit(1)
	}

	err = archiver.ArchivePodcast(fetchedPodcast, *destDirectory, *overwrite, *renameFiles, credentials, *threads)
	if err != nil {
		log.Println("[main] Error: " + err.Error())
	}
}
