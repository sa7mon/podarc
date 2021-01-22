package interfaces

import "time"

type Podcast interface {
	NumEpisodes() int
	GetEpisodes() []PodcastEpisode
	GetTitle() string
	GetDescription() string
	GetPublisher() string
}

type PodcastEpisode interface {
	GetTitle() string
	GetDescription() string
	GetURL() string
	GetPublishedDate() string
	GetParsedPublishedDate() (time.Time, error)
	GetImageURL() string
	ToString() string
	GetGUID() string
}
