package providers

import (
	"errors"
	"fmt"
	"github.com/sa7mon/podarc/internal/interfaces"
	"github.com/sa7mon/podarc/internal/utils"
	"regexp"
)

func FetchPodcastFromUrl(feedUrl string, creds utils.Credentials) (interfaces.Podcast, error) {
	stitcherR := regexp.MustCompile(`https://app\.stitcher\.com/browse/feed/(?P<feedId>\d+)`)
	libsynR := regexp.MustCompile(`\S+\.libsynpro.com/rss`)
	libSynMatches := libsynR.MatchString(feedUrl)
	stitcherMatches := stitcherR.FindStringSubmatch(feedUrl)

	if len(stitcherMatches) > 0 {
		stitcherPod := GetStitcherPodcastFeed(stitcherMatches[1], creds.SessionToken)
		return stitcherPod, nil
	} else if libSynMatches {
		libsynPod := GetLibsynProPodcastFeed(feedUrl)
		return libsynPod, nil
	}
	return nil, errors.New(fmt.Sprintf("Unsupported feed URL '%s'", feedUrl))
}