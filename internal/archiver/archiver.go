package archiver

import (
	"fmt"
	"github.com/sa7mon/podarc/internal/id3-go"
	"github.com/sa7mon/podarc/internal/interfaces"
	"github.com/sa7mon/podarc/internal/utils"
	"github.com/stvp/slug"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
)

func ArchivePodcast(podcast interfaces.Podcast, destDirectory string, overwriteExisting bool, renameFiles bool,
	creds utils.Credentials) error {
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

		headers := make(map[string]string, 1)
		if podcast.GetPublisher() == "Stitcher" {
			headers["Authorization"] = "Bearer " + creds.StitcherNewToken
		}
		err := utils.DownloadFile(episodePath, fileUrl, headers, true)
		if err != nil {
			return err
		}
		// Write ID3 tags to file
		err = WriteID3TagsToFile(episodePath, episode, podcast)
		if err != nil {
			return err
		}
		if renameFiles {
			err := RenameFile(episodePath, episode)
			if err != nil {
				return err
			}
		}
		archivedEpisodes += 1
		fmt.Printf("\r")
		log.Printf("[%s] (%d/%d) archived episode: '%s'", podcast.GetTitle(), archivedEpisodes, len(episodesToArchive), episode.GetTitle())
	}
	return nil
}

func WriteID3TagsToFile(filePath string, episode interfaces.PodcastEpisode, podcast interfaces.Podcast) error {
	/*
		Contains a hacky workaround because the library doesn't support deleting ID3v1 tags.
		We need to use ID3v2 because v1 has a 30-character limit on the title field (and likely others).
		If the file has v1 tags, re-open forcing v2 tags which effectively erases all existing tags
		that we don't set here.
	*/

	file, err := id3.Open(filePath, false)
	if err != nil {
		return err
	}
	//defer file.Close()

	if file.Version()[0:1] == "1" {  // Re-open the file, forcing v2
		log.Println("ID3v1 detected. Re-opening file and forcing ID3v2...")

		file.Close()
		file, err = id3.Open(filePath, true)
		if err != nil {
			log.Println(err)
		}
		defer file.Close()
	}

	file.SetArtist(podcast.GetTitle())
	file.SetTitle(episode.GetTitle())
	file.SetGenre("Podcast")

	publishedDate, err := episode.GetParsedPublishedDate()
	if err != nil {
		return err
	}
	file.SetYear(strconv.Itoa(publishedDate.Year()))
	file.SetDate(episode.GetPublishedDate())
	file.SetReleaseYear(episode.GetPublishedDate())

	// TODO:
	// Set date recorded
	// Save podcast publisher to one of the tags
	// Set cover image

	err = file.Close()
	if err != nil {
		return err
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

func RenameFile(episodeFile string, episode interfaces.PodcastEpisode) error {
	oldDate, _ := episode.GetParsedPublishedDate()
	isoDate := oldDate.Format("2006-01-02")

	slug.Replacement = '-'
	cleanTitle := slug.Clean(episode.GetTitle())
	newName := isoDate + "_" + cleanTitle + filepath.Ext(episodeFile)
	err := os.Rename(episodeFile, filepath.Dir(episodeFile) + "/" + newName)
	if err != nil {
		return err
	}
	return nil
}