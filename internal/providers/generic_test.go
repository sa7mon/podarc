package providers

import (
	"fmt"
	"github.com/sa7mon/podarc/test"
	"testing"
	"time"
)

func TestGenericUnmarshal(t *testing.T) {
	feed,err := GetGenericPodcastFeed("https://feeds.feedburner.com/pod-save-america")
	if err != nil {
		t.Error(err)
	}
	if feed.NumEpisodes() < 450 {
		t.Fatalf("Expected at least 450 episodes. Found: %v", feed.NumEpisodes())
	}
	test.AssertEqual(t, feed.GetTitle(), "Pod Save America")
	test.AssertEqual(t, feed.GetPublisher(), "Crooked Media")
	test.AssertNotEmpty(t, "feedDescription", feed.GetDescription())

	firstEp := feed.GetEpisodes()[feed.NumEpisodes()-1]

	test.AssertNotEmpty(t, "description", firstEp.GetDescription())
	test.AssertNotEmpty(t, "title", firstEp.GetTitle())
	test.AssertNotEmpty(t, "url", firstEp.GetURL())
	test.AssertNotEmpty(t, "imageURL", firstEp.GetImageURL())
	test.AssertNotEmpty(t, "getPublishedDate", firstEp.GetPublishedDate())

	// Test that the first episode's parsed published date is after Jan 01 2017

	jan_1_layout := "Mon, 02 Jan 2006 15:04:05 -0700" // Mon Jan 2 15:04:05 MST 2006
	jan_1, err := time.Parse(jan_1_layout, "Sun, 01 Jan 2017 01:01:01 -0000")
	if err != nil {
		t.Error("Couldn't get test date ready: " + err.Error())
	}

	firstEpPubDate, err := firstEp.GetParsedPublishedDate()
	if err != nil {
		t.Error(err)
	}

	if !firstEpPubDate.After(jan_1) {
		t.Error("Latest episode's parsed published date was not after Jan 1 2017: " + firstEpPubDate.String())
	}

}

func TestGetGenericPodcastFeed(t *testing.T) {
	_, err := GetGenericPodcastFeed("https://httpbin.org/status/404")
	if err == nil {
		t.Error("Trying to get 404 podcast didn't return an error")
	}

	_, err = GetGenericPodcastFeed("https://httpbin.org/status/200")
	if err == nil {
		t.Error("Trying to get empty podcast page didn't return an error")
	}
}

func TestGenericEpisode_ToString(t *testing.T) {
	title := "My Cool Episode"
	description := "It's very cool!"
	link := "https://my.cool.podcast/episode1"
	imageLink := "https://my.cool.podcast/episode1.jpeg"
	publishedDate := "2020-01-05 10:35:42"

	g := GenericEpisode{Title: title, Description: description, PubDate: publishedDate}
	g.Enclosure.URL = link
	g.Image.Href = imageLink

	expected := fmt.Sprintf("Title: %s | Description: %s | Url: %s | PublishedDate: " +
		"%s | ImageUrl: %s", title, description, link, publishedDate, imageLink)

	test.AssertString(t, "GenericEpisode_toString", expected, g.ToString())
}