package archiver

import (
	"github.com/sa7mon/podarc/internal/providers"
	"github.com/sa7mon/podarc/test"
	"testing"
)

func TestGetFileNameFromEpisodeURL(t *testing.T) {
	ep := providers.GenericEpisode{}
	ep.Enclosure.URL = "https://my.site/file.mp3"

	fileName, err := GetFileNameFromEpisodeURL(ep)
	if err != nil {
		t.Error(err)
	}

	test.AssertString(t, "fileName", "file.mp3", fileName)

	ep.Enclosure.URL = "https://my.site/file.mp3?asdf=1"
	test.AssertString(t, "fileName2", "file.mp3", fileName)

	ep2 := providers.GenericEpisode{}
	ep2.Enclosure.URL = "{}[]_=__++!@#$%A^&*()()()"
	fileName, err = GetFileNameFromEpisodeURL(ep2)
	if err == nil {
		t.Error("Bad URL didn't return an error")
	}
}

func TestGetEpisodeFileName(t *testing.T) {
	//layout := "Mon, 02 Jan 2006 15:04:05 -0700" // Mon Jan 2 15:04:05 MST 2006
	//jan02, err := time.Parse(layout, "Mon, 02 Jan 2006 15:04:05 -0700")
	//if err != nil {
	//	t.Error(err)
	//}

	genericEpisode := providers.GenericEpisode{Title: "My Cool Episode", PubDate: "Mon, 02 Jan 2006 15:04:05 -0700"}

	name1 := GetEpisodeFileName("my_cool_episode.mp3", genericEpisode)
	test.AssertString(t, "GetEpisodeFileName_good", "2006-01-02_my-cool-episode.mp3", name1)

	name2 := GetEpisodeFileName("my_cool_episode.mp3?tracking_tag=asdf", genericEpisode)
	test.AssertString(t, "GetEpisodeFileName_good", "2006-01-02_my-cool-episode.mp3", name2)
}