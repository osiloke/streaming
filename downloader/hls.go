package downloader

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var re2 = regexp.MustCompile(`((?:(?:https?|ftp|file)))`)

// ProxyHLSUrls replace hls urls
func ProxyHLSUrls(hlsRaw []byte, proxyServerURL string) ([]byte, error) {
	s := re2.ReplaceAllString(string(hlsRaw), fmt.Sprintf(`%s$1`, proxyServerURL))
	return []byte(s), nil
}

var re = regexp.MustCompile(`(http|ftp|https)://([\w_-]+(?:(?:\.[\w_-]+)+))([\w.,@?^=%&:/~+#-]*[\w@?^=%&/~+#-])?`)

func hlsFilename(url *url.URL) string {
	pt := strings.Split(url.Path, "/")
	filename := fmt.Sprintf("%s%s", pt[len(pt)-2], pt[len(pt)-1])
	return hashKey(filename)
}
func prefixedHlsFilename(prefix string, url *url.URL) string {
	pt := strings.Split(url.Path, "/")
	filename := fmt.Sprintf("%s%s", pt[len(pt)-2], pt[len(pt)-1])
	return hashKey(prefix + filename)
}

// GetSegmentURLS get all segment urls
func GetSegmentURLS(hlsRaw []byte, segmentURLPrefix string) []string {
	s := re.FindAllString(string(hlsRaw), -1)
	found := make([]string, 0)
	for _, v := range s {
		found = append(found, v)
	}
	return found
}

func completeSegmentDownload(ds *DownloadStatus) {
	base := filepath.Dir(ds.TempFilename)
	filename := prefixedHlsFilename(ds.Prefix, mustParseURL(ds.URL))
	name := filepath.Join(base, filename)
	os.Rename(ds.TempFilename, name)
	currenttime := time.Now().Local()
	os.Chtimes(name, currenttime, currenttime)
}

func segmentExists(url, folder string) bool {
	filename := hlsFilename(mustParseURL(url))
	dst := filepath.Join(folder, filename)
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		return false
	}
	return true
}
