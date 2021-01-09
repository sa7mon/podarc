package utils

import "regexp"

func IsValidURL(url string) bool {
	// Taken from https://regexr.com/3ajfi

	r := regexp.MustCompile(`([--:\w?@%&+~#=]*\.[a-z]{2,4}/{0,2})((?:[?&](?:\w+)=(?:\w+))+|[--:\w?@%&+~#=]+)?`)
	return r.MatchString(url)
}