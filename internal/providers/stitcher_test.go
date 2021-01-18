package providers

import (
	"errors"
	"fmt"
	"github.com/sa7mon/podarc/test"
	"strings"
	"testing"
	"time"
)

func TestGetStitcherPodcastFeed(t *testing.T) {
	stitcherPod, err := GetStitcherPodcastFeed("comedy-bang-bang-the-podcast", "")
	if err != nil {
		t.Error(err)
	}
	test.AssertEqual(t, stitcherPod.GetPublisher(), "Stitcher")
	test.AssertEqual(t, stitcherPod.GetTitle(), "Comedy Bang Bang: The Podcast")
	test.AssertNotEmpty(t, "Stitcher description", stitcherPod.GetDescription())

	if stitcherPod.NumEpisodes() < 681 || len(stitcherPod.GetEpisodes()) < 681 {
		t.Errorf("Expected podcast to have at least 681 episodes. Got: %v", stitcherPod.NumEpisodes())
	}

	latestEp := stitcherPod.GetEpisodes()[stitcherPod.NumEpisodes()-1]
	test.AssertNotEmpty(t, "description", latestEp.GetDescription())
	test.AssertNotEmpty(t, "title", latestEp.GetTitle())
	test.AssertNotEmpty(t, "url", latestEp.GetURL())
	if !strings.Contains(latestEp.GetURL(), "mp3") {
		t.Error(errors.New("Latest episode URL didn't contain 'mp3': " + latestEp.GetURL()))
	}

	// Test that the first episode's parsed published date is after Jan 01 2001
	jan1Layout := "Mon, 02 Jan 2006 15:04:05 -0700" // Mon Jan 2 15:04:05 MST 2006
	jan1, err := time.Parse(jan1Layout, "Sun, 01 Jan 2001 01:01:01 -0000")
	if err != nil {
		t.Error("Couldn't get test date ready: " + err.Error())
	}

	firstEpPubDate, err := latestEp.GetParsedPublishedDate()
	if err != nil {
		t.Error(err)
	}

	if !firstEpPubDate.After(jan1) {
		t.Error("Latest episode's parsed published date was not after Jan 1 2001: " + firstEpPubDate.String())
	}


	test.AssertNotEmpty(t, "imageURL", latestEp.GetImageURL())
	test.AssertNotEmpty(t, "getPublishedDate", latestEp.GetPublishedDate())

	_, err = GetStitcherPodcastFeed("^@&%#%#&&#*!@/\\><.,?", "")
	if err == nil {
		t.Error("Bad URL did not return an error")
	}

	_, err = GetStitcherPodcastFeed("asdf-podcast-doesnt-exist", "")
	if err == nil || err.Error() != "no shows found in API response"{
		t.Errorf("Bad URL did not return an error or return an unexpected error: %v", err.Error())
	}
}

func TestStitcherEpisode_ToString(t *testing.T) {
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