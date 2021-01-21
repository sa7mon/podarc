package providers

import (
	"github.com/sa7mon/podarc/internal/interfaces"
	"github.com/sa7mon/podarc/test"
	"testing"
)

func TestGetPatreonPodcastFeed(t *testing.T) {
	_, err := GetPatreonPodcastFeed("https://www.patreon.com/rss/asdfasdf1234")
	if err == nil {
		t.Error("Patreon feed without auth didn't return an error")
	}

	_, err = GetPatreonPodcastFeed("https://www.patreon.com/rss/asdfasdf1234aaaddddd?auth=fakekeyhere")
	if err == nil {
		t.Error("404 Patreon feed didn't return an error")
	}
}

func TestPatreonStruct(t *testing.T) {
	patreonPod := PatreonPodcast{}
	patreonPod.Channel.Title = "My Test Podcast"
	patreonPod.Channel.Description = "This is my cool podcast!"

	patreonEp := PatreonEpisode{Title: "My cool ep",
		Description: "This is a nice ep!"}
	patreonEp.Enclosure.URL = "https://asdf.lol/myep.mp3"
	patreonEp.GUID.Text = "asdf1234"
	patreonEp.ImageURL = "https://asdf.lol/cover.jpeg"

	patreonPod.Episodes = []interfaces.PodcastEpisode{patreonEp}

	test.AssertEqual(t, patreonPod.GetTitle(), "My Test Podcast")
	test.AssertEqual(t, patreonPod.GetDescription(), "This is my cool podcast!")
	test.AssertEqual(t, patreonPod.GetTitle(), "My Test Podcast")

	test.AssertEqual(t, len(patreonPod.GetEpisodes()), 1)
	test.AssertEqual(t, patreonPod.NumEpisodes(), 1)
	test.AssertEqual(t, patreonPod.GetPublisher(), "Patreon")
	test.AssertEqual(t, patreonEp.GetGUID(), "asdf1234")
	test.AssertEqual(t, patreonEp.GetImageURL(), "https://asdf.lol/cover.jpeg")
}