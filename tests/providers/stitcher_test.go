package providers

import (
	"github.com/sa7mon/podarc/internal/providers"
	"github.com/sa7mon/podarc/internal/testutils"
	"github.com/sa7mon/podarc/internal/utils"
	"testing"
)

func TestStitcherUnmarshal(t *testing.T) {
	creds := utils.ReadCredentials("../../creds.json")

	if len(creds.SessionToken) < 20 {
		t.Errorf("Loaded session token missing or invalid.")
	}

	pc := providers.GetStitcherPodcastFeed("467097", creds.SessionToken)

	testutils.AssertString(t, "Podcast Title", "Office Ladies", pc.GetTitle())
	testutils.AssertString(t, "Podcast Description", "The Office co-stars and best friends, Jenna " +
		"Fischer and Angela Kinsey, are doing the ultimate The Office re-watch podcast for you. Each week Jenna and " +
		"Angela will break down an episode of The Office and give exclusive behind the scene stories that only two " +
		"people who were there, can tell you.", pc.GetDescription())
	testutils.AssertString(t, "Publisher", "Stitcher", pc.GetPublisher())

	//firstEpisode := fetchedPodcast.GetEpisodes()[fetchedPodcast.NumEpisodes()-1]
	//testutils.AssertString(t, "Episode Title", "Episode 51- Racist Peruvian Snacks", firstEpisode.GetTitle())
	//testutils.AssertString(t, "Episode Description", "<p>Michael Ian Black and Tom Cavanagh eat snacks and talk about it!</p>", firstEpisode.GetDescription())
	//testutils.AssertString(t, "Episode URL", "http://traffic.libsyn.com/mates/MATES51_Peruvian_Snacks.mp3?dest-id=50920", firstEpisode.GetUrl())
	//testutils.AssertString(t, "Episode Published Date", "Mon, 05 Mar 2012 08:00:00 +0000", firstEpisode.GetPublishedDate())
	//testutils.AssertString(t, "Episode Image URL", "http://static.libsyn.com/p/assets/8/d/b/d/8dbd7e032866e1a8/MATES_logo.jpg", firstEpisode.GetImageUrl())
}
