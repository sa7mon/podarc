package archiver

import (
	"fmt"
	"github.com/sa7mon/podarc/internal/interfaces"
	"github.com/sa7mon/podarc/internal/utils"
	"log"
	"net/url"
	"os"
	"path"
)

func ArchivePodcast(podcast interfaces.Podcast, destDirectory string, overwriteExisting bool) error {
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
	for _, episode := range episodesToArchive {
		fileUrl := episode.GetUrl()
		episodePath := path.Join(destDirectory, GetFileNameFromEpisodeURL(episode.GetUrl()))
		err := utils.DownloadFile(episodePath, fileUrl, false)
		if err != nil {
			return err
		}
		fmt.Printf("[%s] Downloaded %s", podcast.GetTitle(),GetFileNameFromEpisodeURL(episode.GetUrl()))
	}
	return nil
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