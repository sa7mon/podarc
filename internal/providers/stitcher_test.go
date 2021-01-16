package providers

import (
	"github.com/sa7mon/podarc/test"
	"testing"
)

func TestGetStitcherPodcastFeed(t *testing.T) {
	stitcherPod, err := GetStitcherPodcastFeed("comedy-bang-bang-the-podcast", "")
	if err != nil {
		t.Error(err)
	}
	test.AssertEqual(t, stitcherPod.GetPublisher(), "Stitcher")
	test.AssertEqual(t, stitcherPod.GetTitle(), "Comedy Bang Bang: The Podcast")
	if stitcherPod.NumEpisodes() < 681 {
		t.Errorf("Expected podcast to have at least 681 episodes. Got: %v", stitcherPod.NumEpisodes())
	}

	_, err = GetStitcherPodcastFeed("^@&%#%#&&#*!@/\\><.,?", "")
	if err == nil {
		t.Error("Bad URL did not return an error")
	}

	_, err = GetStitcherPodcastFeed("asdf-podcast-doesnt-exist", "")
	if err == nil || err.Error() != "no shows found in API response"{
		t.Errorf("Bad URL did not return an error or return an unexpected error: %v", err.Error())
	}
}