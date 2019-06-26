package downloader

import (
	"regexp"
)

var re = regexp.MustCompile(`(http|ftp|https)://([\w_-]+(?:(?:\.[\w_-]+)+))([\w.,@?^=%&:/~+#-]*[\w@?^=%&/~+#-])?`)

// GetSegmentURLS get all segment urls
func GetSegmentURLS(hlsRaw []byte) []string {
	s := re.FindAllString(string(hlsRaw), -1)
	return s
}
