package providers

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/eduncan911/podcast"
	"github.com/sa7mon/podarc/internal/archiver"
	"github.com/sa7mon/podarc/internal/interfaces"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type GenericPodcast struct {
	XMLName    xml.Name `xml:"rss"`
	Text       string   `xml:",chardata"`
	Version    string   `xml:"version,attr"`
	Itunes     string   `xml:"itunes,attr"`
	Googleplay string   `xml:"googleplay,attr"`
	Atom       string   `xml:"atom,attr"`
	Media      string   `xml:"media,attr"`
	Content    string   `xml:"content,attr"`
	Channel    struct {
		Text string `xml:",chardata"`
		Link struct {
			Text string `xml:",chardata"`
			Href string `xml:"href,attr"`
			Rel  string `xml:"rel,attr"`
			Type string `xml:"type,attr"`
		} `xml:"link"`
		Title       string `xml:"title"`
		Language    string `xml:"language"`
		Copyright   string `xml:"copyright"`
		Description string `xml:"description"`
		Image       struct {
			Text  string `xml:",chardata"`
			Href  string `xml:"href,attr"`
			URL   string `xml:"url"`
			Title string `xml:"title"`
			Link  string `xml:"link"`
		} `xml:"image"`
		Type     string `xml:"type"`
		Subtitle string `xml:"subtitle"`
		Author   string `xml:"author"`
		Summary  string `xml:"summary"`
		Encoded  string `xml:"encoded"`
		Owner    struct {
			Text  string `xml:",chardata"`
			Name  string `xml:"name"`
			Email string `xml:"email"`
		} `xml:"owner"`
		Category []struct {
			Text     string `xml:",chardata"`
			AttrText string `xml:"text,attr"`
			Category []struct {
				Text     string `xml:",chardata"`
				AttrText string `xml:"text,attr"`
			} `xml:"category"`
		} `xml:"category"`
		NewFeedURL string           `xml:"new-feed-url"`
		Items      []GenericEpisode `xml:"item"`
	} `xml:"channel"`
}

type GenericEpisode struct {
	Text        string `xml:",chardata"`
	Title       string `xml:"title"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	EpisodeType string `xml:"episodeType"`
	Episode     string `xml:"episode"`
	Author      string `xml:"author"`
	Subtitle    string `xml:"subtitle"`
	Summary     string `xml:"summary"`
	Encoded     string `xml:"encoded"`
	Duration    string `xml:"duration"`
	GUID        struct {
		Text        string `xml:",chardata"`
		IsPermaLink string `xml:"isPermaLink,attr"`
	} `xml:"guid"`
	Enclosure struct {
		Text   string `xml:",chardata"`
		URL    string `xml:"url,attr"`
		Length string `xml:"length,attr"`
		Type   string `xml:"type,attr"`
	} `xml:"enclosure"`
	Image struct {
		Text string `xml:",chardata"`
		Href string `xml:"href,attr"`
	} `xml:"image"`
	Explicit string `xml:"explicit"`
	Link     string `xml:"link"`
}

func (g GenericPodcast) NumEpisodes() int {
	return len(g.Channel.Items)
}

func (g GenericPodcast) GetEpisodes() []interfaces.PodcastEpisode {
	// Can't return a slice of GenericEpisode. Instead we create a slice of the PodcastEpisode interface
	episodes := make([]interfaces.PodcastEpisode, len(g.Channel.Items))
	for a, _ := range g.Channel.Items {
		episodes[a] = g.Channel.Items[a]
	}
	return episodes
}

func (g GenericPodcast) GetTitle() string {
	return g.Channel.Title
}

func (g GenericPodcast) GetDescription() string {
	return g.Channel.Description
}

func (g GenericPodcast) GetPublisher() string {
	return g.Channel.Author
}

func (e GenericEpisode) GetTitle() string {
	return strings.TrimSpace(e.Title)
}

func (e GenericEpisode) GetDescription() string {
	return e.Description
}

func (e GenericEpisode) GetURL() string {
	return e.Enclosure.URL
}

func (e GenericEpisode) GetPublishedDate() string {
	return e.PubDate
}

func (e GenericEpisode) GetParsedPublishedDate() (time.Time, error) {
	layout := "Mon, 02 Jan 2006 15:04:05 -0700" // Mon Jan 2 15:04:05 MST 2006
	t, err := time.Parse(layout, e.GetPublishedDate())
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func (e GenericEpisode) GetImageURL() string {
	return e.Image.Href
}

func (e GenericEpisode) ToString() string {
	return fmt.Sprintf("Title: %s | Description: %s | Url: %s | PublishedDate: "+
		"%s | ImageUrl: %s", e.GetTitle(), e.GetDescription(), e.GetURL(), e.GetPublishedDate(),
		e.GetImageURL())
}

func (e GenericEpisode) GetGUID() string {
	return e.GUID.Text
}

func (e GenericEpisode) GetDuration() int64 {
	i, err := strconv.ParseInt(e.Duration, 10, 64)
	if err != nil {
		fmt.Printf("Couldn't parse duration '%v' for episode '%v'", e.Duration, e.GetTitle())
		return -1
	}
	return i
}

func GetGenericPodcastFeed(url string) (*GenericPodcast, error) {
	podcast := GenericPodcast{}

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

	episodes := make([]GenericEpisode, len(podcast.Channel.Items))
	for i, elem := range podcast.Channel.Items {
		episodes[i] = elem
	}
	podcast.Channel.Items = episodes

	return &podcast, nil
}

func (gp *GenericPodcast) SaveToFile(filename string) error {
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
	xmlWriter := io.Writer(file)

	enc := xml.NewEncoder(xmlWriter)
	enc.Indent(" ", " ")
	if err := enc.Encode(gp); err != nil {
		return err
	}

	return nil
}

func (gp *GenericPodcast) SaveToFile2(filename string) error {
	feed := podcast.New(gp.GetTitle(), "https://changeme",
		gp.GetDescription(), nil, nil)

	feed.AddSummary(gp.GetDescription())
	feed.AddAuthor(gp.GetPublisher(), "")

	fmt.Println(len(gp.Channel.Items))
	for _, ep := range gp.Channel.Items {
		i := podcast.Item{}
		i.Title = ep.GetTitle()
		i.Description = ep.GetDescription()
		i.AddSummary(ep.GetDescription())
		i.AddImage(ep.GetImageURL())
		i.Link = "https://changeme"
		i.GUID = ep.GetGUID()
		//i.AddDuration(ep.GetDuration())
		pubTime, err := ep.GetParsedPublishedDate()
		if err == nil {
			i.AddPubDate(&pubTime)
		}
		_, err = feed.AddItem(i)
		if err != nil {
			return err
		}
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	xmlWriter := io.Writer(file)

	enc := xml.NewEncoder(xmlWriter)
	enc.Indent("  ", "    ")
	if err := enc.Encode(feed); err != nil {
		return err
	}

	return nil
}
