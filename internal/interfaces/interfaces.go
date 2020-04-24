package interfaces

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
	GetUrl() string
	GetPublishedDate() string
	GetImageUrl() string
	ToString() string
}
