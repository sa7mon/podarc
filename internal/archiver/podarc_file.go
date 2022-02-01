package archiver

import "github.com/eduncan911/podcast"

// TEMP PROBABLY
type PodarcFeed struct {
	Podcast podcast.Podcast
	podcast.Item
}

type PodarcFile struct {
	FeedName string        `json:"feed_name"`
	Episodes []EpisodeFile `json:"episodes"`
}

type EpisodeFile struct {
	GUID     string `json:"guid"`
	FileName string `json:"file_name"`
}
