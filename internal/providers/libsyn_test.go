package providers

import (
	"fmt"
	"github.com/sa7mon/podarc/internal/utils"
	"github.com/sa7mon/podarc/test"
	"testing"
	"time"
)

func TestLibsynUnmarshal(t *testing.T) {
	feedURL := "http://mates.nerdistind.libsynpro.com/rss"
	fetchedPodcast, err := GetLibsynProPodcastFeed(feedURL)
	if err != nil {
		t.Error(err)
	}

	test.AssertString(t, "Podcast Title", "Mike and Tom Eat Snacks", fetchedPodcast.GetTitle())
	test.AssertString(t, "Podcast Description","Michael Ian Black and Tom Cavanagh eat snacks and talk about it!", fetchedPodcast.GetDescription())
	test.AssertString(t, "Publisher", "Libsyn", fetchedPodcast.GetPublisher())

	firstEpisode := fetchedPodcast.GetEpisodes()[fetchedPodcast.NumEpisodes()-1]
	test.AssertString(t, "Episode Title", "Episode 51- Racist Peruvian Snacks", firstEpisode.GetTitle())
	test.AssertString(t, "Episode Description", "<p>Michael Ian Black and Tom Cavanagh eat snacks and talk about it!</p>", firstEpisode.GetDescription())
	test.AssertString(t, "Episode URL", "http://traffic.libsyn.com/mates/MATES51_Peruvian_Snacks.mp3?dest-id=50920", firstEpisode.GetURL())
	test.AssertString(t, "Episode Published Date", "Mon, 05 Mar 2012 08:00:00 +0000", firstEpisode.GetPublishedDate())
	test.AssertString(t, "Episode Image URL", "http://static.libsyn.com/p/assets/8/d/b/d/8dbd7e032866e1a8/MATES_logo.jpg", firstEpisode.GetImageURL())
	test.AssertNotEmpty(t, "GUID", firstEpisode.GetGUID())

	jan1Layout := "Mon, 02 Jan 2006 15:04:05 -0700" // Mon Jan 2 15:04:05 MST 2006
	jan1, err := time.Parse(jan1Layout, "Mon, 01 Jan 2001 01:01:01 -0000")
	if err != nil {
		t.Error("Couldn't get test date ready: " + err.Error())
	}

	firstEpPubDate, err := firstEpisode.GetParsedPublishedDate()
	if err != nil {
		t.Error(err)
	}

	if !firstEpPubDate.After(jan1) {
		t.Error("Latest episode's parsed published date was not after Jan 1 2001: " + firstEpPubDate.String())
	}
}

func TestFetchPodcastFromUrl(t *testing.T) {
	blankCreds := utils.Credentials{}
	p, err := FetchPodcastFromURL("http://mates.nerdistind.libsynpro.com/rss", blankCreds)
	if err != nil {
		t.Errorf(err.Error())
	}

	test.AssertTypesAreEqual(t, p, &LibsynPodcast{})
}

func TestGetLibsynProPodcastFeed(t *testing.T) {
	_, err := GetLibsynProPodcastFeed("https://httpbin.org/status/404")
	if err == nil {
		t.Error("Trying to get 404 podcast didn't return an error")
	}

	_, err = GetLibsynProPodcastFeed("https://httpbin.org/status/200")
	if err == nil {
		t.Error("Trying to get empty podcast page didn't return an error")
	}

}

func TestLibSynEpisode_ToString(t *testing.T) {
	title := "My Cool Episode"
	description := "It's very cool!"
	link := "https://my.cool.podcast/episode1"
	imageLink := "https://my.cool.podcast/episode1.jpeg"
	publishedDate := "2020-01-05 10:35:42"

	g := LibsynEpisode{Title: title, Description: description, Published: publishedDate}
	g.Enclosure.Url = link
	g.Image.ImageURL = imageLink

	expected := fmt.Sprintf("Title: %s | Description: %s | Url: %s | PublishedDate: " +
		"%s | ImageUrl: %s", title, description, link, publishedDate, imageLink)

	test.AssertString(t, "LibsynEpisode_toString", expected, g.ToString())
}
