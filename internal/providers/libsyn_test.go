package providers

import (
	"github.com/sa7mon/podarc/internal/providers"
	"github.com/sa7mon/podarc/internal/utils"
	"github.com/sa7mon/podarc/test"
	"testing"
)

func TestLibsynUnmarshal(t *testing.T) {
	feedURL := "http://mates.nerdistind.libsynpro.com/rss"
	fetchedPodcast := providers.GetLibsynProPodcastFeed(feedURL)

	test.AssertString(t, "Podcast Title", "Mike and Tom Eat Snacks", fetchedPodcast.GetTitle())
	test.AssertString(t, "Podcast Description","Michael Ian Black and Tom Cavanagh eat snacks and talk about it!", fetchedPodcast.GetDescription())
	test.AssertString(t, "Publisher", "Libsyn", fetchedPodcast.GetPublisher())

	firstEpisode := fetchedPodcast.GetEpisodes()[fetchedPodcast.NumEpisodes()-1]
	test.AssertString(t, "Episode Title", "Episode 51- Racist Peruvian Snacks", firstEpisode.GetTitle())
	test.AssertString(t, "Episode Description", "<p>Michael Ian Black and Tom Cavanagh eat snacks and talk about it!</p>", firstEpisode.GetDescription())
	test.AssertString(t, "Episode URL", "http://traffic.libsyn.com/mates/MATES51_Peruvian_Snacks.mp3?dest-id=50920", firstEpisode.GetURL())
	test.AssertString(t, "Episode Published Date", "Mon, 05 Mar 2012 08:00:00 +0000", firstEpisode.GetPublishedDate())
	test.AssertString(t, "Episode Image URL", "http://static.libsyn.com/p/assets/8/d/b/d/8dbd7e032866e1a8/MATES_logo.jpg", firstEpisode.GetImageURL())
}

func TestFetchPodcastFromUrl(t *testing.T) {
	blankCreds := utils.Credentials{}
	p, err := providers.FetchPodcastFromURL("http://mates.nerdistind.libsynpro.com/rss", blankCreds)
	if err != nil {
		t.Errorf(err.Error())
	}

	test.AssertTypesAreEqual(t, p, &providers.LibsynPodcast{})
}