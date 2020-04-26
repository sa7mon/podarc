package main

import (
	"github.com/sa7mon/podarc/internal/id3-go"
	"log"
)

// Open file
// Check ID3 version
// If version 1, remove all tags
// Re-open file
// Check ID3 version


func main() {
	file, err := id3.Open("temp/temp.mp3", true)
	if err != nil {
		log.Println(err)
	}
	log.Println("Version: " + file.Version())
	log.Println(len(file.AllFrames()))

	//file.Tagger.DeleteFrames("TDRC")
	//file.Tagger.DeleteFrames("TORY")
	//file.Tagger.DeleteFrames("TYER")
	//file.Tagger.DeleteFrames("TDOR")

	//ft := v2.V23FrameTypeMap["TIT2"]
	//titleFrame := v2.NewTextFrame(ft, "This is the long name of the episode. It's long.")
	//file.AddFrames(titleFrame)
	file.Close()
}