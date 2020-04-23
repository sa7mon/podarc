package providers

import (
	"podarc/interfaces"
)

type LibsynPodcast struct {
	Title 			string 				`xml:"channel>title"`
	ShowDescription string            	`xml:"channel>itunes:summary"`
	Episodes        []LibsynEpisode 	`xml:"channel>item"`
}

type LibsynEpisode struct {
	Title 		string `xml:"title"`
	Guid        string    `xml:"guid"`
	Image       string    `xml:"itunes:image>href, attr"`
	Description string    `xml:"description"`
	Published   string	  `xml:"pubDate"`
	Url 		string 	  `xml:"enclosure>url, attr"`
}

func (l LibsynPodcast) NumEpisodes() int {
	return len(l.Episodes)
}

func (l LibsynPodcast) GetEpisodes() []interfaces.PodcastEpisode {
	// TODO: Might be more efficient to store these values rather than do a for loop every time the getter is called
	// Golang doesn't allow you to directly return a slice of a type as a slice of an interface
	// https://golang.org/doc/faq#convert_slice_of_interface
	intEpisodes := make([]interfaces.PodcastEpisode, len(l.Episodes))
	for i, elem := range l.Episodes {
		intEpisodes[i] = elem
	}
	return intEpisodes
}

func (l LibsynPodcast) GetTitle() string {
	return l.Title
}

func (l LibsynPodcast) GetDescription() string {
	return l.ShowDescription
}

func (l LibsynPodcast) GetPublisher() string {
	return "Libsyn"
}

func (l LibsynEpisode) GetTitle() string {
	return l.Title
}

func (l LibsynEpisode) GetDescription() string {
	return l.Description
}

func (l LibsynEpisode) GetUrl() string {
	return l.Url
}

func (l LibsynEpisode) GetPublishedDate() string {
	return l.Published
}

func (l LibsynEpisode) GetImageUrl() string {
	return l.Image
}
