package providers

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/sa7mon/podarc/internal/interfaces"
	"html"
	"net/http"
	"time"
)

type PatreonPodcast struct {
	XMLName    xml.Name `xml:"rss"`
	Text       string   `xml:",chardata"`
	Version    string   `xml:"version,attr"`
	Itunes     string   `xml:"itunes,attr"`
	Atom       string   `xml:"atom,attr"`
	Googleplay string   `xml:"googleplay,attr"`
	Channel    struct {
		Text  string `xml:",chardata"`
		Title string `xml:"title"`
		Link  struct {
			Text string `xml:",chardata"`
			Href string `xml:"href,attr"`
			Rel  string `xml:"rel,attr"`
			Type string `xml:"type,attr"`
		} `xml:"link"`
		Description string `xml:"description"`
		Owner       struct {
			Text  string `xml:",chardata"`
			Name  string `xml:"name"`
			Email string `xml:"email"`
		} `xml:"owner"`
		Author string `xml:"author"`
		Image  struct {
			Text  string `xml:",chardata"`
			Href  string `xml:"href,attr"`
			URL   string `xml:"url"`
			Title string `xml:"title"`
			Link  string `xml:"link"`
		} `xml:"image"`
		Block         string `xml:"block"`
		Language      string `xml:"language"`
		PubDate       string `xml:"pubDate"`
		LastBuildDate string `xml:"lastBuildDate"`
		Items          []PatreonEpisode `xml:"item"`
	} `xml:"channel"`

	Episodes []interfaces.PodcastEpisode
}

type PatreonEpisode struct {
	Text        string `xml:",chardata"`
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Enclosure   struct {
		Text   string `xml:",chardata"`
		URL    string `xml:"url,attr"`
		Length string `xml:"length,attr"`
		Type   string `xml:"type,attr"`
	} `xml:"enclosure"`
	Guid struct {
		Text        string `xml:",chardata"`
		IsPermaLink string `xml:"isPermaLink,attr"`
	} `xml:"guid"`
	PubDate string `xml:"pubDate"`
	ImageURL string // Patreon doesn't add an image for every episode. Return feed image
}

func (p PatreonEpisode) GetTitle() string {
	return p.Title
}

func (p PatreonEpisode) GetDescription() string {
	return p.Description
}

func (p PatreonEpisode) GetURL() string {
	return p.Enclosure.URL
}

func (p PatreonEpisode) GetPublishedDate() string {
	return p.PubDate
}

func (p PatreonEpisode) GetParsedPublishedDate() (time.Time, error) {
	layout := "Mon, 02 Jan 2006 15:04:05 MST" // Mon Jan 2 15:04:05 MST 2006
	t, err := time.Parse(layout, p.GetPublishedDate())
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func (p PatreonEpisode) GetImageURL() string {
	return p.ImageURL
}

func (p PatreonEpisode) ToString() string {
	return fmt.Sprintf("Title: %s | Description: %s | Url: %s | PublishedDate: " +
		"%s | ImageUrl: %s", p.GetTitle(), p.GetDescription(), p.GetURL(), p.GetPublishedDate(),
		p.GetImageURL())

}

func (p PatreonPodcast) NumEpisodes() int {
	return len(p.Episodes)
}

func (p PatreonPodcast) GetEpisodes() []interfaces.PodcastEpisode {
	return p.Episodes
}

func (p PatreonPodcast) GetTitle() string {
	return p.Channel.Title
}

func (p PatreonPodcast) GetDescription() string {
	return p.Channel.Description
}

func (p PatreonPodcast) GetPublisher() string {
	return "Patreon"
}

func GetPatreonPodcastFeed(feedURL string) (*PatreonPodcast, error) {
	podcast := &PatreonPodcast{}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest("GET", feedURL, nil)
	if err != nil {
		return podcast, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return podcast, err
	}
	if resp.StatusCode != 200 {
		return podcast, errors.New("Bad status code while getting podcast - " + resp.Status)
	}

	xmlDecoder := xml.NewDecoder(resp.Body)
	err = xmlDecoder.Decode(podcast)
	if err != nil {
		return podcast, err
	}

	// Set episode image URL and unescape download URL
	for _, ep := range podcast.Channel.Items {
		ep.ImageURL = podcast.Channel.Image.Href
		ep.Enclosure.URL = html.UnescapeString(ep.Enclosure.URL)
	}

	episodes := make([]interfaces.PodcastEpisode, len(podcast.Channel.Items))
	for i, elem := range podcast.Channel.Items {
		episodes[i] = elem
	}
	podcast.Episodes = episodes

	return podcast, nil
}
