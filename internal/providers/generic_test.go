package providers

import (
	"github.com/sa7mon/podarc/test"
	"testing"
)

func TestGenericUnmarshal(t *testing.T) {
	feed := GetGenericPodcastFeed("https://feeds.feedburner.com/pod-save-america")
	if feed.NumEpisodes() < 450 {
		t.Fatalf("Expected at least 450 episodes. Found: %v", feed.NumEpisodes())
	}
	test.AssertEqual(t, feed.GetTitle(), "Pod Save America")
	test.AssertEqual(t, feed.GetPublisher(), "Crooked Media")

	firstEp := feed.GetEpisodes()[feed.NumEpisodes()-1]

	test.AssertNotEmpty(t, "description", firstEp.GetDescription())
	test.AssertNotEmpty(t, "title", firstEp.GetTitle())
	test.AssertNotEmpty(t, "url", firstEp.GetURL())
	test.AssertNotEmpty(t, "imageURL", firstEp.GetImageURL())
	test.AssertNotEmpty(t, "getPublishedDate", firstEp.GetPublishedDate())
}