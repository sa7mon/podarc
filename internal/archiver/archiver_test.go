package archiver

import (
	"github.com/sa7mon/podarc/internal/providers"
	"github.com/sa7mon/podarc/test"
	"testing"
)

func TestGetFileNameFromEpisodeURL(t *testing.T) {
	ep := providers.GenericEpisode{}
	ep.Enclosure.URL = "https://my.site/file.mp3"

	fileName := GetFileNameFromEpisodeURL(ep)

	test.AssertString(t, "fileName", "file.mp3", fileName)

	ep.Enclosure.URL = "https://my.site/file.mp3?asdf=1"
	test.AssertString(t, "fileName2", "file.mp3", fileName)
}