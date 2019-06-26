package downloader

import (
	"os"
	"path/filepath"
	"regexp"
)

var re = regexp.MustCompile(`(http|ftp|https)://([\w_-]+(?:(?:\.[\w_-]+)+))([\w.,@?^=%&:/~+#-]*[\w@?^=%&/~+#-])?`)

// GetSegmentURLS get all segment urls
func GetSegmentURLS(hlsRaw []byte) []string {
	s := re.FindAllString(string(hlsRaw), -1)
	return s
}

func completeSegmentDownload(ds *DownloadStatus) {
	base := filepath.Dir(ds.TempFilename)
	filename := hlsFilename(mustParseURL(ds.URL))
	os.Rename(ds.TempFilename, filepath.Join(base, filename))
}

func segmentExists(url, folder string) bool {
	filename := hlsFilename(mustParseURL(url))
	dst := filepath.Join(folder, filename)
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		return false
	}
	return true
}
