package providers

import (
	"github.com/sa7mon/podarc/test"
	"testing"
)

func TestGetStitcherPodcastFeed(t *testing.T) {
	stitcherPod := GetStitcherPodcastFeed("comedy-bang-bang-the-podcast", "")
	test.AssertEqual(t, stitcherPod.GetPublisher(), "Stitcher")
	test.AssertEqual(t, stitcherPod.GetTitle(), "Comedy Bang Bang: The Podcast")
	if stitcherPod.NumEpisodes() < 681 {
		t.Errorf("Expected podcast to have at least 681 episodes. Got: %v", stitcherPod.NumEpisodes())
	}


}