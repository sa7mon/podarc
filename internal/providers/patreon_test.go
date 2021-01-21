package providers

import (
	"github.com/sa7mon/podarc/internal/interfaces"
	"github.com/sa7mon/podarc/test"
	"testing"
	"time"
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
	patreonEp.PubDate = "Tue, 03 Jan 2006 11:04:05 CST"

	patreonPod.Episodes = []interfaces.PodcastEpisode{patreonEp}

	test.AssertEqual(t, patreonPod.GetTitle(), "My Test Podcast")
	test.AssertEqual(t, patreonPod.GetDescription(), "This is my cool podcast!")
	test.AssertEqual(t, patreonPod.GetTitle(), "My Test Podcast")

	test.AssertEqual(t, len(patreonPod.GetEpisodes()), 1)
	test.AssertEqual(t, patreonPod.NumEpisodes(), 1)
	test.AssertEqual(t, patreonPod.GetPublisher(), "Patreon")
	test.AssertEqual(t, patreonEp.GetGUID(), "asdf1234")
	test.AssertEqual(t, patreonEp.GetImageURL(), "https://asdf.lol/cover.jpeg")
	test.AssertEqual(t, patreonEp.GetTitle(), "My cool ep")
	test.AssertEqual(t, patreonEp.GetDescription(), "This is a nice ep!")
	test.AssertEqual(t, patreonEp.GetURL(), "https://asdf.lol/myep.mp3")
	test.AssertEqual(t, patreonEp.GetPublishedDate(), "Tue, 03 Jan 2006 11:04:05 CST")

	expectedPubTime, err := time.Parse("Mon, 02 Jan 2006 15:04:05 MST", "Tue, 03 Jan 2006 11:04:05 CST")
	if err != nil {
		t.Error("Error creating test time: " + err.Error())
	}

	actualPubTime, err := patreonEp.GetParsedPublishedDate()
	if err != nil {
		t.Error(err)
	}

	test.AssertEqual(t, expectedPubTime, actualPubTime)
	test.AssertString(t, "toString()", "Title: My cool ep | Description: This is a nice ep! | Url: https://asdf.lol/myep.mp3 | PublishedDate: Tue, 03 Jan 2006 11:04:05 CST | ImageUrl: https://asdf.lol/cover.jpeg", patreonEp.ToString())
}