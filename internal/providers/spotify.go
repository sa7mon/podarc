package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/librespot-org/librespot-golang/librespot"
	"github.com/librespot-org/librespot-golang/librespot/core"
	"github.com/sa7mon/podarc/internal/interfaces"
	"github.com/zmb3/spotify/v2"
	"github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"
)

func GetSpotifyPodcastFeed(showID string) (interfaces.Podcast, error) {
	feed := SpotifyFeed{}

	ctx := context.Background()
	config := &clientcredentials.Config{
		ClientID:     os.Getenv("SPOTIFY_ID"),
		ClientSecret: os.Getenv("SPOTIFY_SECRET"),
		TokenURL:     spotifyauth.TokenURL,
	}
	apiToken, err := config.Token(ctx)
	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
	}

	feed.Token = apiToken

	httpClient := spotifyauth.New().Client(ctx, apiToken)
	client := spotify.New(httpClient)

	show, err := client.GetShow(ctx, spotify.ID(showID), spotify.Market(spotify.CountryUSA))
	if err != nil {
		log.Fatal(err)
	}

	totalEps := show.Episodes.Total
	pageSize := 50
	pages := int(math.Ceil(float64(totalEps) / float64(pageSize)))
	items := make([]spotify.EpisodePage, totalEps)
	for i := 0; i < pages; i++ {
		fmt.Printf("Requesting page %v of %v. Limit: %v, Offset: %v\n", i, pages, pageSize, i*pageSize)
		// Request pages of episodes
		page, err := client.GetShowEpisodes(ctx, showID,
			spotify.Market(spotify.CountryUSA),
			spotify.Limit(pageSize),
			spotify.Offset(i*pageSize))
		if err != nil {
			fmt.Println(err)
		}

		// We assign the episodes to the array we've properly sized to avoid the memory reallocation of append()
		for j, ep := range page.Episodes {
			items[(i*pageSize)+j] = ep
		}
	}
	fmt.Printf("Found %v episodes\n", len(items))

	// Map API response to SpotifyFeed
	feed.Channel.Description = show.Description
	feed.Channel.Title = show.Name
	feed.Channel.Link.Href = show.Href
	feed.Channel.Image.URL = getLargestImageURL(show.Images)
	feed.Channel.Language = show.Languages[0] // Assumption: use first language listed in feed
	feed.Channel.Explicit = show.Explicit
	feed.Channel.Author = "Spotify"
	feed.Channel.Link.Text = "Spotify"
	feed.Channel.Link.Href = show.ExternalURLs["spotify"]

	// TODO: Category

	feed.Channel.Items = make([]interfaces.PodcastEpisode, len(items))

	for i, ep := range items {
		// TODO: Enclosure
		ge := SpotifyEpisode{GenericEpisode{
			Title:       ep.Name,
			Description: ep.Description,
			PubDate:     ep.ReleaseDate,
			Duration:    strconv.Itoa(ep.Duration_ms / 1000),
			Link:        ep.Href,
			//Author:      "Spotify",
		}}

		ge.Image.Href = getLargestImageURL(ep.Images)
		ge.Image.Text = getLargestImageURL(ep.Images)
		ge.GUID.Text = ep.ID.String()
		ge.GUID.IsPermaLink = "false"

		feed.Channel.Items[i] = ge
	}

	return feed, nil
}

//type SpotifyFeed struct {
//	Channel struct {
//		Title       string `xml:"title"`
//		Description string `xml:"description"`
//		Image       struct {
//			Text  string `xml:",chardata"`
//			Href  string `xml:"href,attr"`
//			URL   string `xml:"url"`
//			Title string `xml:"title"`
//			Link  string `xml:"link"`
//		} `xml:"image"`
//		Language string `xml:"language"`
//		Category []struct {
//			Text     string `xml:",chardata"`
//			AttrText string `xml:"text,attr"`
//			Category []struct {
//				Text     string `xml:",chardata"`
//				AttrText string `xml:"text,attr"`
//			} `xml:"category"`
//		} `xml:"category"`
//		Explicit bool   `xml:"explicit"`
//		Author   string `xml:"author"`
//		Text     string `xml:",chardata"`
//		Link     struct {
//			Text string `xml:",chardata"`
//			Href string `xml:"href,attr"`
//			Rel  string `xml:"rel,attr"`
//			Type string `xml:"type,attr"`
//		} `xml:"link"`
//		Owner struct {
//			Text  string `xml:",chardata"`
//			Name  string `xml:"name"`
//			Email string `xml:"email"`
//		} `xml:"owner"`
//		Copyright string `xml:"copyright"`
//
//		//Type     string `xml:"type"`
//		//Summary  string `xml:"summary"`
//		Items []SpotifyEpisode `xml:"item"`
//	} `xml:"channel"`
//}

type SpotifyFeed struct {
	GenericPodcast
}

func getLargestImageURL(images []spotify.Image) string {
	largestImageWidth := 0
	largestImageURL := ""
	for _, image := range images {
		if image.Width > largestImageWidth {
			largestImageWidth = image.Width
			largestImageURL = image.URL
		}
	}
	return largestImageURL
}

type SpotifyEpisode struct {
	GenericEpisode
}

//type SpotifyEpisode struct {
//	Title     string `xml:"title"`
//	Enclosure struct {
//		Text   string `xml:",chardata"`
//		URL    string `xml:"url,attr"`
//		Length string `xml:"length,attr"`
//		Type   string `xml:"type,attr"`
//	} `xml:"enclosure"`
//	GUID struct {
//		Text        string `xml:",chardata"`
//		IsPermaLink string `xml:"isPermaLink,attr"`
//	} `xml:"guid"`
//	PubDate     string `xml:"pubDate"`
//	Description string `xml:"description"`
//
//	// Duration of episode in seconds
//	Duration int    `xml:"duration"`
//	Link     string `xml:"link"`
//	Image    struct {
//		Text string `xml:",chardata"`
//		Href string `xml:"href,attr"`
//	} `xml:"image"`
//	Explicit string `xml:"explicit"`
//
//	// itunes:title
//	// itunes:episode
//	EpisodeType string `xml:"episodeType"`
//	Episode     string `xml:"episode"`
//
//	//Text        string `xml:",chardata"`
//	//Author      string `xml:"author"`
//	//Subtitle    string `xml:"subtitle"`
//	//Summary     string `xml:"summary"`
//	//Encoded     string `xml:"encoded"`
//}

func (feed SpotifyFeed) GetEpisodes() []interfaces.PodcastEpisode {
	return feed.Channel.Items
}

func (feed SpotifyFeed) NumEpisodes() int {
	return len(feed.Channel.Items)
}

func (feed SpotifyFeed) GetTitle() string {
	return feed.Channel.Title
}

func (feed SpotifyFeed) GetDescription() string {
	return feed.Channel.Description
}

func (SpotifyFeed) GetPublisher() string {
	return "Spotify"
}

func (e SpotifyEpisode) GetTitle() string {
	return e.Title
}

func (e SpotifyEpisode) GetDescription() string {
	return e.Description
}

func (e SpotifyEpisode) GetURL() string {

	return fmt.Sprintf(`https://api-partner.spotify.com/pathfinder/v1/query?operationName=getEpisode&variables={"uri":"spotify:episode:%v"}&extensions={"persistedQuery":{"version":1,"sha256Hash":"224ba0fd89fcfdfb3a15fa2d82a6112d3f4e2ac88fba5c6713de04d1b72cf482"}}`,
		e.GUID.Text)
}

func (e SpotifyEpisode) GetPublishedDate() string {
	return e.PubDate
}

func (e SpotifyEpisode) GetParsedPublishedDate() (time.Time, error) {
	layout := "Mon, 02 Jan 2006 15:04:05 -0700" // Mon Jan 2 15:04:05 MST 2006
	t, err := time.Parse(layout, e.GetPublishedDate())
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func (e SpotifyEpisode) GetImageURL() string {
	return e.Image.Href
}

func (e SpotifyEpisode) ToString() string {
	return fmt.Sprintf("Title: %s | Description: %s | Url: %s | PublishedDate: "+
		"%s | ImageUrl: %s", e.GetTitle(), e.GetDescription(), e.GetURL(), e.GetPublishedDate(),
		e.GetImageURL())
}

func (e SpotifyEpisode) GetGUID() string {
	return e.GUID.Text
}

func GetSpotifyFileURL(e interfaces.PodcastEpisode, token *oauth2.Token) error {
	filesUrl := fmt.Sprintf(`https://api-partner.spotify.com/pathfinder/v1/query?operationName=getEpisode&variables={"uri":"spotify:episode:%v"}&extensions={"persistedQuery":{"version":1,"sha256Hash":"224ba0fd89fcfdfb3a15fa2d82a6112d3f4e2ac88fba5c6713de04d1b72cf482"}}`,
		e.GetGUID())

	client := &http.Client{}

	// Get the data
	req, err := http.NewRequest("GET", filesUrl, nil)
	if err != nil {
		return err
	}

	token.SetAuthHeader(req)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("app-platform", "WebPlayer")
	req.Header.Add("Accept-Encoding", "gzip, deflate")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	epResp := GetEpisodeResponse{}
	jsonDecoder := json.NewDecoder(resp.Body)
	err = jsonDecoder.Decode(&epResp)

	a, err := librespot.LoginOAuth("", "", "")
	a.Player().LoadTrack()

	core.Session{}.Player().LoadTrack()

	return nil
}

type GetEpisodeResponse struct {
	Data struct {
		Episode struct {
			Audio struct {
				Items []struct {
					ExternallyHosted bool   `json:"externallyHosted"`
					FileID           string `json:"fileId"`
					Format           string `json:"format"`
					URL              string `json:"url"`
				} `json:"items"`
			} `json:"audio"`
			AudioPreview struct {
				Format string `json:"format"`
				URL    string `json:"url"`
			} `json:"audioPreview"`
			ContentRating struct {
				Label string `json:"label"`
			} `json:"contentRating"`
			CoverArt struct {
				Sources []struct {
					Height int    `json:"height"`
					URL    string `json:"url"`
					Width  int    `json:"width"`
				} `json:"sources"`
			} `json:"coverArt"`
			Description string `json:"description"`
			Duration    struct {
				TotalMilliseconds int `json:"totalMilliseconds"`
			} `json:"duration"`
			HTMLDescription string `json:"htmlDescription"`
			ID              string `json:"id"`
			Name            string `json:"name"`
			Playability     struct {
				Playable bool   `json:"playable"`
				Reason   string `json:"reason"`
			} `json:"playability"`
			PlayedState struct {
				PlayPositionMilliseconds int    `json:"playPositionMilliseconds"`
				State                    string `json:"state"`
			} `json:"playedState"`
			Podcast struct {
				CoverArt struct {
					Sources []struct {
						Height int    `json:"height"`
						URL    string `json:"url"`
						Width  int    `json:"width"`
					} `json:"sources"`
				} `json:"coverArt"`
				Name      string      `json:"name"`
				ShowTypes []string    `json:"showTypes"`
				Trailer   interface{} `json:"trailer"`
				URI       string      `json:"uri"`
			} `json:"podcast"`
			ReleaseDate struct {
				IsoString time.Time `json:"isoString"`
			} `json:"releaseDate"`
			Segments struct {
				Segments struct {
					TotalCount int `json:"totalCount"`
				} `json:"segments"`
			} `json:"segments"`
			SharingInfo struct {
				ShareID  string `json:"shareId"`
				ShareURL string `json:"shareUrl"`
			} `json:"sharingInfo"`
			Type string `json:"type"`
			URI  string `json:"uri"`
		} `json:"episode"`
	} `json:"data"`
}
