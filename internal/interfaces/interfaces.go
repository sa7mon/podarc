package interfaces

type Podcast interface {
	numEpisodes() int
	GetEpisodes() []PodcastEpisode
	getTitle() string
	getDescription() string
	getPublisher() string
}

type PodcastEpisode interface {
	GetTitle() string
	getDescription() string
	getUrl() string
	getPublishedDate() string
	getImageUrl() string
}
