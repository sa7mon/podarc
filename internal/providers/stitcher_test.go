package providers

import (
	"github.com/sa7mon/podarc/internal/utils"
	"github.com/sa7mon/podarc/test"
	"testing"
)

func TestStitcherUnmarshal(t *testing.T) {
	testCreds := utils.Credentials{SessionToken: "asdf1234_session_token", StitcherNewToken: "asdf_stitcher_new_token"}

	if len(testCreds.SessionToken) < 20 {
		t.Errorf("Loaded session token missing or invalid.")
	}

	pc := GetStitcherPodcastFeed("467097", testCreds.SessionToken)

	test.AssertString(t, "Podcast Title", "Office Ladies", pc.GetTitle())
	test.AssertString(t, "Podcast Description", "The Office co-stars and best friends, Jenna " +
		"Fischer and Angela Kinsey, are doing the ultimate The Office re-watch podcast for you. Each week Jenna and " +
		"Angela will break down an episode of The Office and give exclusive behind the scene stories that only two " +
		"people who were there, can tell you.", pc.GetDescription())
	test.AssertString(t, "Publisher", "Stitcher", pc.GetPublisher())

	firstEpisode := pc.GetEpisodes()[pc.NumEpisodes()-1]
	test.AssertString(t, "Episode Title", "Happy New Year!", firstEpisode.GetTitle())
	test.AssertString(t, "Episode Description", "The Office Ladies are taking this week off," +
		" but we have a special New Years Eve memory shared by the one and only, Creed Bratton. We'll be back " +
		"January 8th with The Fire.&nbsp;", firstEpisode.GetDescription())
	test.AssertString(t, "Episode URL", "https://s3.amazonaws.com/stitcher.assets/audio/premium_required.mp3", firstEpisode.GetURL())
	test.AssertString(t, "Episode Published Date", "2019-12-31 13:04:00", firstEpisode.GetPublishedDate())
	test.AssertString(t, "Episode Image URL", "https://secureimg.stitcher.com/feedimageswide/480x270_467097.jpg", firstEpisode.GetImageURL())
}

func TestFetchSticherPodcastFromUrl(t *testing.T) {
	testCreds := utils.Credentials{SessionToken: "asdf1234_session_token", StitcherNewToken: "asdf_stitcher_new_token"}

	if len(testCreds.SessionToken) < 20 {
		t.Errorf("Loaded session token missing or invalid.")
	}
	p, err := FetchPodcastFromURL(" https://app.stitcher.com/browse/feed/467097/details", testCreds)
	if err != nil {
		t.Errorf(err.Error())
	}

	test.AssertTypesAreEqual(t, p, &StitcherPodcast{})
}