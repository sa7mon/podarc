package providers

import (
	"github.com/sa7mon/podarc/internal/utils"
	"github.com/sa7mon/podarc/test"
	"testing"
)

func TestStitcherUnmarshal(t *testing.T) {
	creds := utils.ReadCredentials("../../creds.json")

	if len(creds.SessionToken) < 20 {
		t.Errorf("Loaded session token missing or invalid.")
	}

	pc := GetStitcherPodcastFeed("467097", creds.SessionToken)

	test.AssertString(t, "Podcast Title", "Office Ladies", pc.GetTitle())
	test.AssertString(t, "Podcast Description", "The Office co-stars and best friends, Jenna " +
		"Fischer and Angela Kinsey, are doing the ultimate The Office re-watch podcast for you. Each week Jenna and " +
		"Angela will break down an episode of The Office and give exclusive behind the scene stories that only two " +
		"people who were there, can tell you.", pc.GetDescription())
	test.AssertString(t, "Publisher", "Stitcher", pc.GetPublisher())

	firstEpisode := pc.GetEpisodes()[pc.NumEpisodes()-1]
	test.AssertString(t, "Episode Title", "Office Ladies Trailer", firstEpisode.GetTitle())
	test.AssertString(t, "Episode Description", "Join Jenna Fischer and Angela Kinsey as " +
		"they give you a sneak peak of what's to come. " +
		"Office Ladies premieres October 16th!", firstEpisode.GetDescription())
	test.AssertString(t, "Episode URL", "https://cloudfront.wolfpub.io/OL-000.2-20190913-TrailerFinished.mp3", firstEpisode.GetURL())
	test.AssertString(t, "Episode Published Date", "2019-09-25 04:00:33", firstEpisode.GetPublishedDate())
	test.AssertString(t, "Episode Image URL", "https://secureimg.stitcher.com/feedimageswide/480x270_467097.jpg", firstEpisode.GetImageURL())
}

func TestFetchSticherPodcastFromUrl(t *testing.T) {
	creds := utils.ReadCredentials("../../creds.json")

	if len(creds.SessionToken) < 20 {
		t.Errorf("Loaded session token missing or invalid.")
	}
	p, err := FetchPodcastFromURL(" https://app.stitcher.com/browse/feed/467097/details", creds)
	if err != nil {
		t.Errorf(err.Error())
	}

	test.AssertTypesAreEqual(t, p, &StitcherPodcast{})
}