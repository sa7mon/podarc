package providerss

import (
	"podarc/interfaces"
)

type StitcherPodcast struct {
	Feed            StitcherFeed      `xml:"feed"`
	ShowDescription string            `xml:"show_description,attr"`
	Episodes        []StitcherEpisode `xml:"episodes>episode"`
}

type StitcherEpisodes struct {
	Episodes []StitcherEpisode
}
type StitcherFeed struct {
	Name          string          `xml:"name"`
	Description   string          `xml:"description"`
	LatestEpisode StitcherEpisode `xml:"episode"`
	Premium       bool            `xml:"premium,attr"`
	EpisodeCount  int             `xml:"episodeCount,attr"`
}

type StitcherEpisode struct {
	Id          string    `xml:"id,attr"`
	Image       string    `xml:"episodeImage,attr"`
	Published   string	  `xml:"published,attr"`
	Title       string    `xml:"title"`
	Description string    `xml:"description"`
	Url 		string 	  `xml:"url, attr"`
}

func (s StitcherPodcast) numEpisodes() int {
	panic("implement me")
}

func (s StitcherPodcast) getEpisodes() []interfaces.PodcastEpisode {
	panic("implement me")
}

func (s StitcherPodcast) getTitle() string {
	panic("implement me")
}

func (s StitcherPodcast) getDescription() string {
	panic("implement me")
}

func (s StitcherPodcast) getPublisher() string {
	panic("implement me")
}

