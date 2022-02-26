package providers

import (
	"github.com/sa7mon/podarc/internal/interfaces"
	"github.com/sa7mon/podarc/internal/utils"
	"regexp"
)

func FetchPodcastFromURL(feedURL string, creds utils.Credentials) (interfaces.Podcast, error) {
	libsynR := regexp.MustCompile(`\S+\.libsynpro.com/rss`)
	stitcherR := regexp.MustCompile(`https://www\.stitcher\.com/show/(?P<slug>[a-zA-Z0-9-]+)`)
	patreonR := regexp.MustCompile(`http(?:s)*://www\.patreon\.com/rss/.+`)
	acastR := regexp.MustCompile(`http(?:s)*://.+\.memberfulcontent\.com/rss/\d+`)

	libSynMatches := libsynR.MatchString(feedURL)
	stitcherMatches := stitcherR.FindStringSubmatch(feedURL)
	patreonMatches := patreonR.MatchString(feedURL)
	acastMatches := acastR.MatchString(feedURL)

	if len(stitcherMatches) > 0 {
		return GetStitcherPodcastFeed(stitcherMatches[1], creds.StitcherNewToken)
	} else if libSynMatches {
		libsynPod, err := GetLibsynProPodcastFeed(feedURL)
		return libsynPod, err
	} else if patreonMatches {
		patreonPod, err := GetPatreonPodcastFeed(feedURL)
		return patreonPod, err
	} else if acastMatches {
		acastPod, err := GetAcastPodcastFeed(feedURL)
		return acastPod, err
	} else {
		genericPod, err := GetGenericPodcastFeed(feedURL)
		return genericPod, err
	}
}
