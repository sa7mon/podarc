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

func main() {
	feedUrl := flag.String("feedUrl", "", "URL of podcast feed to archive. (Required)")
	destDirectory := flag.String("outputDir", "", "Directory to save the files into. (Required)")
	overwrite := flag.Bool("overwrite", false, "Overwrite episodes already downloaded. Default: false")
	renameFiles := flag.Bool("renameFiles", true, "Rename downloaded files to friendly names.")
	flag.Parse()

	if *feedUrl == "" || !utils.IsValidUrl(*feedUrl){
		fmt.Printf("Error - Invalid feedUrl: '%s'\n", *feedUrl)
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *destDirectory == "" {
		fmt.Printf("Error - Invalid outputDir: '%s'\n", *destDirectory)
		flag.PrintDefaults()
		os.Exit(1)
	}

	credentials := utils.ReadCredentials("creds.json")

	fetchedPodcast, err := providers.FetchPodcastFromUrl(*feedUrl, credentials)
	if err != nil {
		log.Println("Error fetching podcast from URL - " + err.Error())
		os.Exit(1)
	}

	err = archiver.ArchivePodcast(fetchedPodcast, *destDirectory, *overwrite, *renameFiles, credentials)
	if err != nil {
		log.Println("Error: " + err.Error())
	}
}