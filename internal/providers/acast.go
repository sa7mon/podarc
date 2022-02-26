package providers

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/sa7mon/podarc/internal/archiver"
	"github.com/sa7mon/podarc/internal/interfaces"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type AcastPodcast struct {
	XMLName     xml.Name `xml:"rss"`
	Text        string   `xml:",chardata"`
	Version     string   `xml:"version,attr"`
	Atom        string   `xml:"atom,attr"`
	Googleplay  string   `xml:"googleplay,attr"`
	Itunes      string   `xml:"itunes,attr"`
	ItunesXmlns string   `xml:"xmlns:itunes,attr"`
	Media       string   `xml:"media,attr"`
	Podaccess   string   `xml:"podaccess,attr"`
	Acast       string   `xml:"acast,attr"`
	Channel     struct {
		Text      string `xml:",chardata"`
		Ttl       string `xml:"ttl"`
		Generator string `xml:"generator"`
		Title     string `xml:"title"`
		Link      struct {
			Text string `xml:",chardata"`
			Href string `xml:"href,attr"`
			Rel  string `xml:"rel,attr"`
			Type string `xml:"type,attr"`
		} `xml:"link"`
		Language    string `xml:"language"`
		Copyright   string `xml:"copyright"`
		Keywords    string `xml:"keywords"`
		Author      string `xml:"author"`
		Subtitle    string `xml:"subtitle"`
		Summary     string `xml:"summary"`
		Description string `xml:"description"`
		Explicit    string `xml:"explicit"`
		Owner       struct {
			Text  string `xml:",chardata"`
			Name  string `xml:"name"`
			Email string `xml:"email"`
		} `xml:"owner"`
		ShowId    string `xml:"showId"`
		ShowUrl   string `xml:"showUrl"`
		Signature struct {
			Text      string `xml:",chardata"`
			Key       string `xml:"key,attr"`
			Algorithm string `xml:"algorithm,attr"`
		} `xml:"signature"`
		Settings string `xml:"settings"`
		Type     string `xml:"type"`
		Image    struct {
			Text  string `xml:",chardata"`
			Href  string `xml:"href,attr"`
			URL   string `xml:"url"`
			Link  string `xml:"link"`
			Title string `xml:"title"`
		} `xml:"image"`
		Category []struct {
			Text     string `xml:",chardata"`
			AttrText string `xml:"text,attr"`
			Category struct {
				Text     string `xml:",chardata"`
				AttrText string `xml:"text,attr"`
			} `xml:"category"`
		} `xml:"category"`
		Items []AcastEpisode `xml:"item"`
		Auth  struct {
			Text  string `xml:",chardata"`
			Login struct {
				Text     string `xml:",chardata"`
				Type     string `xml:"type,attr"`
				Provider string `xml:"provider"`
				URL      string `xml:"url"`
				ShowId   string `xml:"showId"`
			} `xml:"login"`
			User struct {
				Text     string `xml:",chardata"`
				Status   string `xml:"status,attr"`
				Username string `xml:"username"`
				Provider string `xml:"provider"`
			} `xml:"user"`
		} `xml:"auth"`
		Block string `xml:"block"`
	} `xml:"channel"`
}

type AcastEpisode struct {
	Text           string `xml:",chardata"`
	Locked         string `xml:"locked,attr"`
	Ads            string `xml:"ads,attr"`
	Spons          string `xml:"spons,attr"`
	AttrPremium    string `xml:"premium,attr"`
	Title          string `xml:"title"`
	PubDate        string `xml:"pubDate"`
	Duration       string `xml:"duration"`
	ItunesDuration string `xml:"itunes:duration,omitempty"`
	Enclosure      struct {
		Text   string `xml:",chardata"`
		URL    string `xml:"url,attr"`
		Length string `xml:"length,attr"`
		Type   string `xml:"type,attr"`
	} `xml:"enclosure"`
	GUID struct {
		Text        string `xml:",chardata"`
		IsPermaLink string `xml:"isPermaLink,attr"`
	} `xml:"guid"`
	Explicit    string `xml:"explicit"`
	Link        string `xml:"link"`
	EpisodeId   string `xml:"episodeId"`
	EpisodeUrl  string `xml:"episodeUrl"`
	Settings    string `xml:"settings"`
	EpisodeType string `xml:"episodeType"`
	Episode     string `xml:"episode"`
	Image       struct {
		Text   string `xml:",chardata"`
		Locked string `xml:"locked,attr"`
		Ads    string `xml:"ads,attr"`
		Spons  string `xml:"spons,attr"`
	} `xml:"image"`
	Summary     string `xml:"summary"`
	Description string `xml:"description"`
	Premium     struct {
		Text   string `xml:",chardata"`
		Locked string `xml:"locked,attr"`
	} `xml:"premium"`
}

func (g AcastPodcast) NumEpisodes() int {
	return len(g.Channel.Items)
}

func (g AcastPodcast) GetEpisodes() []interfaces.PodcastEpisode {
	// Can't return a slice of GenericEpisode. Instead we create a slice of the PodcastEpisode interface
	episodes := make([]interfaces.PodcastEpisode, len(g.Channel.Items))
	for a, _ := range g.Channel.Items {
		episodes[a] = g.Channel.Items[a]
	}
	return episodes
}

func (g AcastPodcast) GetTitle() string {
	return g.Channel.Title
}

func (g AcastPodcast) GetDescription() string {
	return g.Channel.Description
}

func (g AcastPodcast) GetPublisher() string {
	return g.Channel.Author
}

func (e AcastEpisode) GetTitle() string {
	return strings.TrimSpace(e.Title)
}

func (e AcastEpisode) GetDescription() string {
	return e.Description
}

func (e AcastEpisode) GetURL() string {
	return e.Enclosure.URL
}

func (e AcastEpisode) GetPublishedDate() string {
	return e.PubDate
}

func (e AcastEpisode) GetParsedPublishedDate() (time.Time, error) {
	layout := "Mon, 02 Jan 2006 15:04:05 MST" // Mon Jan 2 15:04:05 MST 2006
	t, err := time.Parse(layout, e.GetPublishedDate())
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func (e AcastEpisode) GetImageURL() string {
	return e.Image.Text
}

func (e AcastEpisode) ToString() string {
	return fmt.Sprintf("Title: %s | Description: %s | Url: %s | PublishedDate: "+
		"%s | ImageUrl: %s", e.GetTitle(), e.GetDescription(), e.GetURL(), e.GetPublishedDate(),
		e.GetImageURL())
}

func (e AcastEpisode) GetGUID() string {
	return e.GUID.Text
}

func (e AcastEpisode) GetDuration() int64 {
	i, err := strconv.ParseInt(e.Duration, 10, 64)
	if err != nil {
		fmt.Printf("Couldn't parse duration '%v' for episode '%v'", e.Duration, e.GetTitle())
		return -1
	}
	return i
}

func GetAcastPodcastFeed(url string) (*AcastPodcast, error) {
	podcast := AcastPodcast{}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return &podcast, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return &podcast, err
	}
	if resp.StatusCode != 200 {
		return &podcast, errors.New("Bad status code while getting podcast - " + resp.Status)
	}

	xmlDecoder := xml.NewDecoder(resp.Body)
	err = xmlDecoder.Decode(&podcast)
	if err != nil {
		return &podcast, err
	}

	podcast.ItunesXmlns = "http://www.itunes.com/dtds/podcast-1.0.dtd"

	episodes := make([]AcastEpisode, len(podcast.Channel.Items))
	for i, elem := range podcast.Channel.Items {
		elem.ItunesDuration = elem.Duration // Manually set the <itunes:duration> tag
		episodes[i] = elem
	}
	podcast.Channel.Items = episodes

	return &podcast, nil
}

func (gp *AcastPodcast) SaveToFile(filename string) error {
	// Starting with the scraped feed:
	//     Replace enclosure URL
	//     Update enclosure length

	for i, ep := range gp.Channel.Items {
		remoteFileName, err := archiver.GetFileNameFromEpisodeURL(ep)
		if err != nil {
			return err
		}
		localFileName := archiver.GetEpisodeFileName(remoteFileName, ep)
		gp.Channel.Items[i].Enclosure.URL = fmt.Sprintf("{PODARC_BASE_URL}/%v", localFileName)
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	xmlWriter := XmlWriter{File: file}

	enc := xml.NewEncoder(xmlWriter)
	enc.Indent(" ", " ")
	if err := enc.Encode(gp); err != nil {
		return err
	}

	return nil
}

// Custom FileWriter to replace &#xA; with newline characters
// TODO: Put this somewhere common

type XmlWriter struct {
	File *os.File
}

func (w XmlWriter) Close() error                   { return w.File.Close() }
func (w XmlWriter) CloseWithError(err error) error { return nil }
func (w XmlWriter) Write(data []byte) (n int, err error) {
	n = len(data)
	data = bytes.Replace(data, []byte("&#xA;"), []byte("\n"), -1)
	_, err = w.File.Write(data)
	return
}
