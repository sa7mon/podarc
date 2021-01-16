package archiver

import (
	"fmt"
	"github.com/sa7mon/podarc/internal/id3"
	"github.com/sa7mon/podarc/internal/interfaces"
	"github.com/sa7mon/podarc/internal/utils"
	"github.com/stvp/slug"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

var concurrentDownloads = 2

func ArchivePodcast(podcast interfaces.Podcast, destDirectory string, overwriteExisting bool, renameFiles bool,
	creds utils.Credentials) error {
	var episodesToArchive []interfaces.PodcastEpisode

	for _, episode := range podcast.GetEpisodes() {
		if overwriteExisting {
			episodesToArchive = append(episodesToArchive, episode)
		} else {   // if file does not exist in destDirectory, add to episodesToArchive
			var episodeFileName string
			if renameFiles {
				episodeFileName = GetEpisodeFileName(episode.GetURL(), episode) // Clean/normalize audio file name
			} else {
				episode, err := GetFileNameFromEpisodeURL(episode) // Leave file name as it is in the URL
				if err != nil {
					return err
				}
				episodeFileName = episode
			}
			episodePath := path.Join(destDirectory, episodeFileName)
			if _, err := os.Stat(episodePath); os.IsNotExist(err) {
				episodesToArchive = append(episodesToArchive, episode)
			}
		}
	}
	log.Printf("[%s] [archiver] Found %d total episodes", podcast.GetTitle(), len(podcast.GetEpisodes()))
	log.Printf("[%s] [archiver] Found %d episodes to archive", podcast.GetTitle(), len(episodesToArchive))

	archivedEpisodes := 0
	// For each episode not currently downloaded - download it.

	var channel = make(chan int, concurrentDownloads)

	for _, episode := range episodesToArchive {
		channel <- 1
		go func(episodeToArchive interfaces.PodcastEpisode) error {
			fileURL := episodeToArchive.GetURL()
			fileName, err := GetFileNameFromEpisodeURL(episodeToArchive)
			if err != nil {
				return err
			}
			episodePath := path.Join(destDirectory, fileName)

			headers := make(map[string]string, 1)
			if podcast.GetPublisher() == "Stitcher" {
				valid, reason := utils.IsStitcherTokenValid(creds.StitcherNewToken)
				if !valid {
					log.Fatal("Bad Stitcher token: " + reason)
				}
				headers["Authorization"] = "Bearer " + creds.StitcherNewToken
			}
			log.Printf("[%s] [archiver] Downloading episode '%s'...", podcast.GetTitle(), episodeToArchive.GetTitle())
			err = utils.DownloadFile(episodePath, fileURL, headers, false)
			if err != nil {
				return err
			}
			// Write ID3 tags to file
			err = WriteID3TagsToFile(episodePath, episodeToArchive, podcast)
			if err != nil {
				return err
			}
			if renameFiles {
				err = os.Rename(episodePath, GetEpisodeFileName(episodePath, episodeToArchive))
				if err != nil {
					return err
				}
				return nil
			}
			archivedEpisodes++
			fmt.Printf("\r")
			log.Printf("[%s] [archiver] (%d/%d) archived episode: '%s'", podcast.GetTitle(), archivedEpisodes, len(episodesToArchive), episodeToArchive.GetTitle())
			<-channel
			return nil
		}(episode)
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

/**
	Returns the name of the file this episode should be saved as
 */
func GetEpisodeFileName(episodeFile string, episode interfaces.PodcastEpisode) string {
	oldDate, _ := episode.GetParsedPublishedDate()
	isoDate := oldDate.Format("2006-01-02")

	slug.Replacement = '-'
	cleanTitle := slug.Clean(episode.GetTitle())
	extension := filepath.Ext(episodeFile)
	if strings.ContainsRune(extension, '?') {
		extension = extension[0:strings.Index(extension, "?")]
	}
	newName := isoDate + "_" + cleanTitle + extension
	return newName
}

func GetFileNameFromEpisodeURL(episode interfaces.PodcastEpisode) (string, error) {
	parsed, err := url.Parse(episode.GetURL())
	if err != nil {
		return "", err
	}

	// url.Path returns the path portion of the URL (without query parameters)
	// path.Base() returns everything after the final slash
	return path.Base(parsed.Path), nil
}