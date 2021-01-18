package providers

import (
	"github.com/sa7mon/podarc/internal/interfaces"
	"github.com/sa7mon/podarc/internal/utils"
	"regexp"
)

func FetchPodcastFromURL(feedURL string, creds utils.Credentials) (interfaces.Podcast, error) {
	libsynR := regexp.MustCompile(`\S+\.libsynpro.com/rss`)
	stitcherR := regexp.MustCompile(`https://www\.stitcher\.com/show/(?P<slug>[a-zA-Z0-9-]+)`)

	libSynMatches := libsynR.MatchString(feedURL)
	stitcherMatches := stitcherR.FindStringSubmatch(feedURL)

	if len(stitcherMatches) > 0 {
		return GetStitcherPodcastFeed(stitcherMatches[1], creds.StitcherNewToken)
	} else if libSynMatches {
		libsynPod, err := GetLibsynProPodcastFeed(feedURL)
		return libsynPod, err
	} else {
		genericPod, err := GetGenericPodcastFeed(feedURL)
		return genericPod, err
	}
}