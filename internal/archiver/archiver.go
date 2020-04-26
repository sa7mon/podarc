package archiver

import (
	"github.com/sa7mon/podarc/internal/interfaces"
	"log"
	"net/url"
	"os"
	"path"
)

func ArchivePodcast(podcast interfaces.Podcast, destDirectory string, overwriteExisting bool) {
	var episodesToArchive []interfaces.PodcastEpisode

	for _, episode := range podcast.GetEpisodes() {
		if overwriteExisting {
			episodesToArchive = append(episodesToArchive, episode)
		} else {   // if file does not exist in destDirectory, add to episodesToArchive
			episodeFileName := GetFileNameFromEpisodeURL(episode.GetUrl())
			episodePath := path.Join(destDirectory, episodeFileName)
			if _, err := os.Stat(episodePath); os.IsNotExist(err) {
				episodesToArchive = append(episodesToArchive, episode)
			}
		}
	}

	log.Printf("[%s] Found %d episodes to archive", podcast.GetTitle(), len(episodesToArchive))

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