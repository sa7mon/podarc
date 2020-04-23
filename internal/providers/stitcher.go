package providers

import (
	"podarc/interfaces"
)

/*************************
*
*		Structs
*
**************************/

type StitcherPodcast struct {
	Feed            StitcherFeed      `xml:"feed"`
	ShowDescription string            `xml:"show_description,attr"`
	Episodes        []StitcherEpisode `xml:"episodes>episode"`
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

/*************************
*
*    Getters / Setters
*
**************************/


func (s StitcherPodcast) NumEpisodes() int {
	panic("implement me")
}

func (s StitcherPodcast) GetEpisodes() []interfaces.PodcastEpisode {
	panic("implement me")
}

func (s StitcherPodcast) GetTitle() string {
	panic("implement me")
}

func (s StitcherPodcast) GetDescription() string {
	panic("implement me")
}

func (s StitcherPodcast) GetPublisher() string {
	panic("implement me")
}

func (s StitcherEpisode) GetTitle() string {
	panic("implement me")
}

func (s StitcherEpisode) GetDescription() string {
	panic("implement me")
}

func (s StitcherEpisode) GetUrl() string {
	panic("implement me")
}

func (s StitcherEpisode) GetPublishedDate() string {
	panic("implement me")
}

func (s StitcherEpisode) GetImageUrl() string {
	panic("implement me")
}