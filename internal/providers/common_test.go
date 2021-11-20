package providers

import (
	"github.com/sa7mon/podarc/internal/utils"
	"github.com/sa7mon/podarc/test"
	"testing"
)

func TestFetchPodcastFromURL(t *testing.T) {
	// Generic podcast URL
	generic, err := FetchPodcastFromURL("https://feeds.feedburner.com/pod-save-america", utils.Credentials{})
	if err != nil {
		t.Error(err)
	}
	test.AssertEqual(t, generic.GetTitle(), "Pod Save America")
	test.AssertEqual(t, generic.GetPublisher(), "Crooked Media")

	libsyn, err := FetchPodcastFromURL("http://mates.nerdistind.libsynpro.com/rss", utils.Credentials{})
	if err != nil {
		t.Error(err)
	}
	test.AssertString(t, "Podcast Title", "Mike and Tom Eat Snacks", libsyn.GetTitle())
	test.AssertString(t, "Publisher", "Libsyn", libsyn.GetPublisher())

	stitcher, err := FetchPodcastFromURL("https://www.stitcher.com/show/comedy-bang-bang-the-podcast", utils.Credentials{})
	if err != nil {
		t.Error(err)
	}
	test.AssertEqual(t, stitcher.GetPublisher(), "Stitcher")
	test.AssertEqual(t, stitcher.GetTitle(), "Comedy Bang Bang: The Podcast")

	firstEp := stitcher.GetEpisodes()[0]
	test.AssertNotEmpty(t, "GUID", firstEp.GetGUID())

	_, err = FetchPodcastFromURL("https://www.patreon.com/rss/asdfasdf1234", utils.Credentials{})
	if err == nil {
		t.Error("404 Patreon feed didn't return an error")
	}

}