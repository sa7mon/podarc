package providers

import (
	"github.com/sa7mon/podarc/internal/interfaces"
	"github.com/sa7mon/podarc/internal/utils"
	"regexp"
)

func FetchPodcastFromUrl(feedURL string, creds utils.Credentials) (interfaces.Podcast, error) {
	stitcherR := regexp.MustCompile(`https://app\.stitcher\.com/browse/feed/(?P<feedId>\d+)`)
	libsynR := regexp.MustCompile(`\S+\.libsynpro.com/rss`)
	stitcherNewR := regexp.MustCompile(`https://www\.stitcher\.com/show/(?P<slug>[a-zA-Z0-9-]+)`)

	libSynMatches := libsynR.MatchString(feedURL)
	stitcherMatches := stitcherR.FindStringSubmatch(feedURL)
	stitcherNewMatches := stitcherNewR.FindStringSubmatch(feedURL)

	if len(stitcherMatches) > 0 {
		stitcherPod := GetStitcherPodcastFeed(stitcherMatches[1], creds.SessionToken)
		return stitcherPod, nil
	} else if len(stitcherNewMatches) > 0 {
		stitcherNewPod := GetStitcherNewPodcastFeed(stitcherNewMatches[1], creds.StitcherNewToken)
		return stitcherNewPod, nil
	} else if libSynMatches {
		libsynPod := GetLibsynProPodcastFeed(feedURL)
		return libsynPod, nil
	} else {
		genericPod := GetGenericPodcastFeed(feedURL)
		return genericPod, nil
	}
}