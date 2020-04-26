package archiver

import (
	"github.com/mikkyang/id3-go"
	v2 "github.com/mikkyang/id3-go/v2"
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

	archivedEpisodes := 0
	// For each episode not currently downloaded - download it.
	for _, episode := range episodesToArchive {
		fileUrl := episode.GetUrl()
		episodePath := path.Join(destDirectory, GetFileNameFromEpisodeURL(episode.GetUrl()))
		err := utils.DownloadFile(episodePath, fileUrl, false)
		if err != nil {
			return err
		}
		// Write ID3 tags to file
		err = WriteID3TagsToFile(episodePath, episode, podcast)
		if err != nil {
			return err
		}
		archivedEpisodes += 1
		log.Printf("[%s] (%d/%d) Downloaded %s", podcast.GetTitle(), archivedEpisodes, len(episodesToArchive), episode.GetTitle())
	}
	return nil
}

func WriteID3TagsToFile(filePath string, episode interfaces.PodcastEpisode, podcast interfaces.Podcast) error {
	file, err := id3.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	log.Println("Version: " + file.Version())

	file.SetArtist(podcast.GetTitle())
	file.SetTitle(episode.GetTitle())
	file.SetGenre("Podcast")

	//log.Println("Detected v2 tag. Writing publisher...")
	ft := v2.V23FrameTypeMap["TIT2"]
	titleFrame := v2.NewTextFrame(ft, episode.GetTitle())
	//allFrames = append(allFrames, textFrame)
	file.AddFrames(titleFrame)

	//if file.Tagger.Version()[0:1] == "2" {
	//	log.Println("Detected v2 tag. Writing publisher...")
	//	ft := v2.V23FrameTypeMap["TPUB"]
	//	textFrame := v2.NewTextFrame(ft, podcast.GetPublisher())
	//	//allFrames = append(allFrames, textFrame)
	//	file.AddFrames(textFrame)
	//}

	// Set year recorded
	// Save podcast publisher to one of the tags
	// Set cover image

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