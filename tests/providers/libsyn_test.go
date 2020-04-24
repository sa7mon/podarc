package providers

import (
	"github.com/sa7mon/podarc/internal/providers"
	"testing"
)

func AssertString(t *testing.T, valueName string, expected string, got string) {
	if expected != got {
		t.Errorf("%s wrong - expected: '%s', got: '%s'", valueName, expected, got)
	}
}

func TestLibsynUnmarshal(t *testing.T) {
	feedUrl := "http://mates.nerdistind.libsynpro.com/rss"
	fetchedPodcast := providers.GetLibsynProPodcastFeed(feedUrl)

	AssertString(t, "Podcast Title", "Mike and Tom Eat Snacks", fetchedPodcast.GetTitle())
	AssertString(t, "Podcast Description","Michael Ian Black and Tom Cavanagh eat snacks and talk about it!", fetchedPodcast.GetDescription())
	AssertString(t, "Publisher", "Libsyn", fetchedPodcast.GetPublisher())

	firstEpisode := fetchedPodcast.GetEpisodes()[fetchedPodcast.NumEpisodes()-1]
	AssertString(t, "Episode Title", "Episode 51- Racist Peruvian Snacks", firstEpisode.GetTitle())
	AssertString(t, "Episode Description", "<p>Michael Ian Black and Tom Cavanagh eat snacks and talk about it!</p>", firstEpisode.GetDescription())
	AssertString(t, "Episode URL", "http://traffic.libsyn.com/mates/MATES51_Peruvian_Snacks.mp3?dest-id=50920", firstEpisode.GetUrl())
	AssertString(t, "Episode Published Date", "Mon, 05 Mar 2012 08:00:00 +0000", firstEpisode.GetPublishedDate())
	AssertString(t, "Episode Image URL", "http://static.libsyn.com/p/assets/8/d/b/d/8dbd7e032866e1a8/MATES_logo.jpg", firstEpisode.GetImageUrl())
}
