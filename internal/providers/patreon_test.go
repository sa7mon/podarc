package providers

import "testing"

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