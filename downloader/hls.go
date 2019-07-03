package downloader

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var re2 = regexp.MustCompile(`((?:(?:https?|ftp|file)))`)

func hashKey(data string) string {
	sha := sha256.Sum256([]byte(data))
	return hex.EncodeToString(sha[:])
}

// ProxyHLSUrls replace hls urls
func ProxyHLSUrls(hlsRaw []byte, proxyServerURL string) ([]byte, error) {
	s := re2.ReplaceAllString(string(hlsRaw), fmt.Sprintf(`%s$1`, proxyServerURL))
	return []byte(s), nil
}

var re = regexp.MustCompile(`(http|ftp|https)://([\w_-]+(?:(?:\.[\w_-]+)+))([\w.,@?^=%&:/~+#-]*[\w@?^=%&/~+#-])?`)

func idAndFile(url *url.URL) []string {
	pt := strings.Split(url.Path, "/")
	id := pt[len(pt)-2]
	hlspart := pt[len(pt)-1]
	return []string{strings.Split(id, "_trd")[0], hlspart}
}
func hlsFilename(url *url.URL) string {
	pt := strings.Split(url.Path, "/")
	filename := fmt.Sprintf("%s%s", pt[len(pt)-2], pt[len(pt)-1])
	return hashKey(filename)
}
func PrefixedHlsFilename(prefix string, url *url.URL) string {
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

func randate() time.Time {
	min := time.Date(1987, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Now().Unix()
	delta := max - min

	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0)
}

func completeSegmentDownload(ds *DownloadStatus) {
	base := filepath.Dir(ds.TempFilename)
	filename := PrefixedHlsFilename(ds.Prefix, mustParseURL(ds.URL))
	name := filepath.Join(base, filename)
	// os.Rename(ds.TempFilename, name)
	currenttime := randate()
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
