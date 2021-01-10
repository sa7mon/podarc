package providers

import (
	"github.com/sa7mon/podarc/test"
	"testing"
)

func TestGenericUnmarshal(t *testing.T) {
	feed := GetGenericPodcastFeed("https://feeds.feedburner.com/pod-save-america")
	if feed.NumEpisodes() <= 450 {
		t.Fatalf("Expected at least 450 episodes. Found: %v", feed.NumEpisodes())
	}
	test.AssertEqual(t, feed.GetTitle(), "Pod Save America")
	test.AssertEqual(t, feed.GetPublisher(), "Crooked Media")

	firstEp := feed.GetEpisodes()[0]

	test.AssertNotEmpty(t, firstEp.GetDescription())
	test.AssertNotEmpty(t, firstEp.GetTitle())
	test.AssertNotEmpty(t, firstEp.GetURL())
	test.AssertNotEmpty(t, firstEp.GetImageURL())
	test.AssertNotEmpty(t, firstEp.GetPublishedDate())
}