package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/sa7mon/podarc/internal/archiver"
	"github.com/sa7mon/podarc/internal/providers"
	"github.com/sa7mon/podarc/internal/utils"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

func FolderExists(dirPath string) bool {
	info, err := os.Stat(dirPath)
	if os.IsNotExist(err) { // Path does not exist
		return false
	} else {
		if !info.IsDir() { // Path exists but it's a file
			return false
		}
		return true
	}
}

func FileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) { // Path does not exist
		return false
	} else {
		if info.IsDir() { // Path exists but it's a directory
			return false
		}
		return true
	}
}

func serve(podDir *string, baseUrl *string, bindPort uint16) {
	fs := http.FileServer(http.Dir(*podDir))

	http.Handle("/files/", http.StripPrefix("/files", fs))
	http.HandleFunc("/feed/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(405)
			return
		}

		log.Println("GET " + r.URL.Path)
		requestedFeed := strings.TrimLeft(r.URL.Path, "/feed/")

		if FolderExists(path.Join(*podDir, requestedFeed)) {
			feedFile := path.Join(*podDir, requestedFeed, "feed.xml")
			if FileExists(feedFile) {
				feedFileBytes, err := ioutil.ReadFile(feedFile)
				if err != nil {
					fmt.Println(err)
				}

				w.Header().Set("Content-Type", "application/xml")
				w.WriteHeader(200)
				w.Write(bytes.Replace(feedFileBytes, []byte("{PODARC_BASE_URL}"), []byte(*baseUrl+"/files/"+strings.TrimPrefix(requestedFeed, "/")), -1))
				return
			}
		}
		w.WriteHeader(404)
		w.Write([]byte("404 feed not found"))
	})

	serveAddress := fmt.Sprintf("localhost:%v", bindPort)
	fmt.Println("Server listening on: " + serveAddress)
	log.Fatal(http.ListenAndServe(serveAddress, nil))
}

func main() {

	archiveCmd := flag.NewFlagSet("archive", flag.ExitOnError)

	feedURL := archiveCmd.String("feedUrl", "", "URL of podcast feed to archive. (Required)")
	destDirectory := archiveCmd.String("outputDir", "", "Directory to save the files into. (Required)")
	overwrite := archiveCmd.Bool("overwrite", false, "Overwrite episodes already downloaded. Default: false")
	renameFiles := archiveCmd.Bool("renameFiles", true, "Rename downloaded files to friendly names.")
	threads := archiveCmd.Int("threads", 2, "Number of threads to use when downloading")

	serveCmd := flag.NewFlagSet("serve", flag.ExitOnError)
	podDir := serveCmd.String("dir", "", "Directory containing podcasts to serve.")
	baseUrl := serveCmd.String("baseUrl", "http://localhost:8282", "Base URL at which to serve podcasts. Default: http://localhost:8282")
	bindPort := serveCmd.Uint("bindPort", 8282, "Local port to bind to")

	if len(os.Args) < 2 {
		fmt.Println("expected 'archive' or 'serve' subcommand")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "archive":
		archiveCmd.Parse(os.Args[2:])
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

	case "serve":
		serveCmd.Parse(os.Args[2:])
		if *bindPort > 65535 {
			fmt.Printf("Error - bind port must be within 0 - 65535\n")
			flag.PrintDefaults()
			os.Exit(1)
		}
		serve(podDir, baseUrl, uint16(*bindPort))
	default:
		fmt.Println("expected 'archive' or 'serve' subcommand")
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
