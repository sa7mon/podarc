package providers

import (
	"fmt"
	"github.com/sa7mon/podarc/internal/interfaces"
	"github.com/sa7mon/podarc/internal/utils"
	"regexp"
)

func FetchPodcastFromUrl(feedUrl string, creds utils.Credentials) interfaces.Podcast {
	stitcherR := regexp.MustCompile(`https://app\.stitcher\.com/browse/feed/(?P<feedId>\d+)`)
	libsynR := regexp.MustCompile(`\S+\.libsynpro.com/rss`)
	libSynMatches := libsynR.MatchString(feedUrl)
	stitcherMatches := stitcherR.FindStringSubmatch(feedUrl)

	if len(stitcherMatches) > 0 {
		fmt.Println("Stitcher feed detected")
		fmt.Println("Feed ID: " + stitcherMatches[1]) // Capture group names available via: stitcherR.SubexpNames()

		stitcherPod := GetStitcherPodcastFeed(stitcherMatches[1], creds.SessionToken)
		return stitcherPod
	} else if libSynMatches {
		fmt.Println("Libsyn Pro feed detected")
		libsynPod := GetLibsynProPodcastFeed(feedUrl)
		return libsynPod
	}
	panic("Unknown URL!")
}