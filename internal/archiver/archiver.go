package archiver

import (
	"github.com/sa7mon/podarc/internal/interfaces"
	"log"
	"net/url"
	"os"
	"path"
)

func ArchivePodcast(podcast interfaces.Podcast, destDirectory string) {
	// Determine which episodes have not already been downloaded

	var episodesToArchive []interfaces.PodcastEpisode

	for _, episode := range podcast.GetEpisodes() {
		episodeFileName := GetFileNameFromEpisodeURL(episode.GetUrl())

		// if file does not exist in destDirectory, add to episodesToArchive
		episodePath := path.Join(destDirectory, episodeFileName)
		log.Println("DEBUG: Checking if file exists: " + episodePath)
		if _, err := os.Stat(episodePath); os.IsNotExist(err) {
			episodesToArchive = append(episodesToArchive, episode)
		}
	}

	log.Printf("Found %d episodes to archive", len(episodesToArchive))

	// For each episode not currently downloaded - download it.

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

func GetFileNameFromEpisodeURL(fullUrl string) string {
	parsed, err := url.Parse(fullUrl)
	if err != nil {
		log.Println(err)
	}

	// url.Path returns the path portion of the URL (without query parameters)
	// path.Base() returns everything after the final slash
	return path.Base(parsed.Path)
}