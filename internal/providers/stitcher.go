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

type StitcherFeed struct { 		  // TODO: Possibly move these properties to StitcherPodcast and delete this struct
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
	return s.Feed.EpisodeCount
}

func (s StitcherPodcast) GetEpisodes() []interfaces.PodcastEpisode {
	panic("implement me")
}

func (s StitcherPodcast) GetTitle() string {
	return s.Feed.Name
}

func (s StitcherPodcast) GetDescription() string {
	return s.ShowDescription
}

func (s StitcherPodcast) GetPublisher() string {
	return "Stitcher"
}

func (s StitcherEpisode) GetTitle() string {
	return s.Title
}

func (s StitcherEpisode) GetDescription() string {
	return s.Description
}

func (s StitcherEpisode) GetUrl() string {
	return s.Url
}

func (s StitcherEpisode) GetPublishedDate() string {
	return s.Published
}

func (s StitcherEpisode) GetImageUrl() string {
	return s.Image
}