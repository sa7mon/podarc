package main

import (
	"github.com/sa7mon/podarc/internal/id3-go"
	"github.com/sa7mon/podarc/internal/interfaces"
	"log"
)

func WriteID3TagsToFile(filePath string, episode interfaces.PodcastEpisode, podcast interfaces.Podcast) error {


	return nil
}

func main() {
	/*
		Hacky workaround because the library doesn't support deleting ID3v1 tags.
		We need to use ID3v2 because v1 has a 30-character limit on the title field (and likely others).
		If the file has v1 tags, re-open forcing v2 tags which effectively erases all existing tags
		that we don't set here.
	*/

	filePath := "./temp/ep-1.mp3"

	file, err := id3.Open(filePath, false)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if file.Version()[0:1] == "1" {
		log.Println("ID3v1 detected. Re-opening file and forcing ID3v2...")
		file.Close()
		file, err = id3.Open(filePath, true)
		if err != nil {
			log.Println(err)
		}
		defer file.Close()
	}

	file.SetArtist("This is this artist!")
	file.SetTitle("This is the title!")
	file.SetGenre("Podcast")

	file.SetYear("2020")
	file.SetDate("2020-07-01")
	file.SetReleaseYear("2020-07-01")

	file.Close()

	//f, err := os.Open(filePath)
	//if err != nil {
	//	panic(err)
	//}
	//defer f.Close()
	//
	//tagger := v2.NewTag(2)
	//ft := v2.V23FrameTypeMap["COMM"]
	//at := v2.V23FrameTypeMap["TIT2"]
	//dateFrame := v2.NewUnsynchTextFrame(ft, "Comment", "This is my comment")
	//titleFrame := v2.NewTextFrame(at, "This is my title!")
	//tagger.AddFrames(dateFrame, titleFrame)
	//
	//_, err = f.Write(tagger.Bytes())
	//f.Close()

	// Set date recorded
	// Save podcast publisher to one of the tags
	// Set cover image
}