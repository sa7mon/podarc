/*
	Much of the consumer-producer pattern code copied from: https://medium.com/hdac/producer-consumer-pattern-implementation-with-golang-6ac412cf941c
*/

package archiver

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/sa7mon/podarc/internal/id3"
	"github.com/sa7mon/podarc/internal/interfaces"
	"github.com/sa7mon/podarc/internal/utils"
	"github.com/stvp/slug"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"unicode"
)

/*
	State to pass to each consumer
	Only access after locking the sync.Mutex to ensure thread safety
*/
type ArchiveState struct {
	archivedCount  uint32
	toArchiveCount uint32
}

type State struct {
	mutex               sync.Mutex
	queue               Queue
	err                 error
	archivedCount       uint32
	totalToArchiveCount uint32
}

/*
	FIFO queue
*/
type Queue struct {
	items []interfaces.PodcastEpisode
}

func (q Queue) String() string {
	return fmt.Sprintf("%v", q.items)
}

/*
	Constructor
*/
func NewQueue(initial []interfaces.PodcastEpisode) Queue {
	return Queue{items: initial}
}

func (q *Queue) Add(item interfaces.PodcastEpisode) {
	q.items = append(q.items, item)
}

func (q *Queue) Get() interfaces.PodcastEpisode {
	if len(q.items) == 0 {
		return nil
	}
	popped := q.items[0]  // Get top element
	q.items = q.items[1:] // Remove top element
	return popped
}

func (q *Queue) Length() int {
	return len(q.items)
}

func Work(state *State, wg *sync.WaitGroup, workerID int, podcast interfaces.Podcast, destDirectory string,
	renameFiles bool, creds utils.Credentials) {
	for {
		// Get work to do
		state.mutex.Lock()
		if state.err != nil {
			// Error happened somewhere - kill this thread
			state.mutex.Unlock()
			wg.Done()
			return
		}
		episode := state.queue.Get()
		state.mutex.Unlock()
		if episode == nil {
			//fmt.Printf("[worker%v] No work to do. Exiting\n", workerID)
			wg.Done()
			return
		}

		// If an episode has no download link (like a feed item that's just a note to the listener), skip it
		if strings.TrimSpace(episode.GetURL()) == "" {
			log.Printf("[%s] [archiver] Episode '%v' has no download link. Skipping...", podcast.GetTitle(),
				episode.GetTitle())
			state.mutex.Lock()
			state.archivedCount++
			log.Printf("[%s] [archiver] (%d/%d) archived episode: '%s'", podcast.GetTitle(), state.archivedCount,
				state.totalToArchiveCount, episode.GetTitle())
			state.mutex.Unlock()
			continue
		}

		fileURL := episode.GetURL()
		fileName, err := GetFileNameFromEpisodeURL(episode)
		if err != nil {
			state.mutex.Lock()
			state.err = err
			state.mutex.Unlock()
			wg.Done()
			return
		}

		// Generate a UUID and append a weird sanitized version of the episode name to the end
		// this suffix is just for debugging purposes as the file will get immediately renamed after downloading
		fileName = uuid.New().String() + SanitizeFileName(episode.GetTitle()) + filepath.Ext(fileName)

		// Patreon allows identical file names across episodes.
		// Use the GUID for the filename.
		// TODO: Move this logic elsewhere or use GUID file names for all providers
		//if podcast.GetPublisher() == "Patreon" {
		//	fileName = fmt.Sprintf("%v_%v", episode.GetGUID(), fileName)
		//}
		episodePath := path.Join(destDirectory, fileName)

		headers := make(map[string]string, 1)
		if podcast.GetPublisher() == "Stitcher" {
			valid, reason := utils.IsStitcherTokenValid(creds.StitcherNewToken)
			if !valid {
				state.mutex.Lock()
				state.err = errors.New("Bad Stitcher token: " + reason)
				state.mutex.Unlock()
				wg.Done()
				return
			}
			headers["Authorization"] = "Bearer " + creds.StitcherNewToken
		}
		log.Printf("[%s] [archiver] Downloading episode '%s'...", podcast.GetTitle(), episode.GetTitle())
		err = utils.DownloadFile(episodePath, fileURL, headers, false)
		if err != nil {
			state.mutex.Lock()
			state.err = err
			state.mutex.Unlock()
			wg.Done()
			return
		}
		// Write ID3 tags to file
		err = WriteID3TagsToFile(episodePath, episode, podcast)
		if err != nil {
			state.mutex.Lock()
			state.err = err
			state.mutex.Unlock()
			wg.Done()
			return
		}

		if renameFiles {
			err = os.Rename(episodePath, path.Join(destDirectory, GetEpisodeFileName(episodePath, episode)))
			if err != nil {
				state.mutex.Lock()
				state.err = err
				state.mutex.Unlock()
				wg.Done()
				return
			}
		}

		state.mutex.Lock()
		state.archivedCount++
		log.Printf("[%s] [archiver] (%d/%d) archived episode: '%s'", podcast.GetTitle(), state.archivedCount,
			state.totalToArchiveCount, episode.GetTitle())
		state.mutex.Unlock()
	}

}

func ArchivePodcast(podcast interfaces.Podcast, destDirectory string, overwriteExisting bool, renameFiles bool,
	creds utils.Credentials, threads int) error {
	var episodesToArchive []interfaces.PodcastEpisode

	log.Printf("[%s] [archiver] Found %d total episodes", podcast.GetTitle(), len(podcast.GetEpisodes()))
	for _, episode := range podcast.GetEpisodes() {
		if strings.TrimSpace(episode.GetURL()) == "" {
			log.Printf("[%s] [archiver] Episode '%v' has no download link. Ignoring episode", podcast.GetTitle(), episode.GetTitle())
			continue
		}
		if overwriteExisting {
			episodesToArchive = append(episodesToArchive, episode)
		} else { // if file does not exist in destDirectory, add to episodesToArchive
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
	log.Printf("[%s] [archiver] Found %d episodes to archive", podcast.GetTitle(), len(episodesToArchive))

	// Setup producers/consumers
	nConsumers := threads
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Instantiate a thread-safe state object
	archiveState := &State{err: nil, queue: NewQueue(episodesToArchive),
		totalToArchiveCount: uint32(len(episodesToArchive))}

	wg := &sync.WaitGroup{}
	wg.Add(nConsumers)
	for i := 0; i < nConsumers; i++ {
		go Work(archiveState, wg, i, podcast, destDirectory, renameFiles, creds)
	}

	wg.Wait()
	feedFile := path.Join(destDirectory, "feed.xml")
	err := podcast.SaveToFile(feedFile)
	if err != nil {
		return errors.New("[archiver] error saving podcast to file: " + err.Error())
	}
	return archiveState.err
}

/*
	Contains a hacky workaround because the library doesn't support deleting ID3v1 tags.
	We need to use ID3v2 because v1 has a 30-character limit on the title field (and likely others).
	If the file has v1 tags, re-open forcing v2 tags which effectively erases all existing tags
	that we don't set here.

	TODO:
		- Set date recorded
		- Save podcast publisher to one of the tags
		- Set cover image
*/
func WriteID3TagsToFile(filePath string, episode interfaces.PodcastEpisode, podcast interfaces.Podcast) error {

	file, err := id3.Open(filePath, false)
	if err != nil {
		return err
	}
	//defer file.Close()

	if file.Version()[0:1] == "1" { // Re-open the file, forcing v2
		log.Println("ID3v1 detected. Re-opening file and forcing ID3v2...")

		file.Close()
		file, err = id3.Open(filePath, true)
		if err != nil {
			return err
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

/*
	Returns the file name from an episode URL.

	Example: "https://my.site/podcast/episode1.mp3?asdf=1" -> "episode1.mp3"
*/
func GetFileNameFromEpisodeURL(episode interfaces.PodcastEpisode) (string, error) {
	parsed, err := url.Parse(episode.GetURL())
	if err != nil {
		return "", err
	}

	// url.Path returns the path portion of the URL (without query parameters)
	// path.Base() returns everything after the final slash
	return path.Base(parsed.Path), nil
}

func SanitizeFileName(dirty string) string {
	clean := ""
	for _, char := range dirty {
		if unicode.IsLetter(char) || unicode.IsNumber(char) {
			clean = clean + string(char)
		} else {
			clean = clean + "_"
		}
	}
	return clean
}
