package downloader

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cavaliercoder/grab"
	"github.com/cskr/pubsub"
)

//EventBus sends events to a listener
type EventBus interface {
	SendEvent(channel, message string)
}

// DownloadStatusChannel is the pubsub channel
const DownloadStatusChannel = "download_status"

// DownloadStatus represents the status of a download
type DownloadStatus struct {
	URL          string `json:"url"`
	TempFilename string `json:"tempFilename"`
	Progress     string `json:"progress"`
	Status       string `json:"status"`
	Error        string `json:"error"`
}

func hashKey(data string) string {
	sha := sha256.Sum256([]byte(data))
	return hex.EncodeToString(sha[:])
}

func mustParseURL(urlSt string) *url.URL {
	u, _ := url.Parse(urlSt)
	return u
}

func hlsFilename(url *url.URL) string {
	pt := strings.Split(url.Path, "/")
	filename := fmt.Sprintf("%s%s", pt[len(pt)-2], pt[len(pt)-1])
	return hashKey((filename))
}
func prefixedHlsFilename(prefix string, url *url.URL) string {
	pt := strings.Split(url.Path, "/")
	filename := fmt.Sprintf("%s%s", pt[len(pt)-2], pt[len(pt)-1])
	return hashKey((prefix + filename))
}

// DownloadHLSURL download a url to a file
func DownloadHLSURL(url *url.URL, filename, folder string, ps *pubsub.PubSub) ([]byte, error) {
	start := time.Now()
	// done := make(chan int64)
	resp, err := http.Get(url.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	dst := filepath.Join(folder, filename)
	out, err := os.Create(dst)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		out.Close()
		return nil, err
	}
	out.Close()

	buf, err := ioutil.ReadFile(dst)
	if err != nil {
		return nil, err
	}
	// done <- n
	elapsed := time.Since(start)
	log.Printf("downloaded %s to %s - %s", url, dst, elapsed)

	return []byte(strings.TrimSpace(string(buf))), err
}

// DownloadSegmentURLs takes an array of urls to be downloaded
func DownloadSegmentURLs(urls []string, folder string, ps *pubsub.PubSub, client *grab.Client) error {
	reqs := make([]*grab.Request, 0)
	for i := 0; i < len(urls); i++ {
		filename := hlsFilename(mustParseURL(urls[i]))
		dst := filepath.Join(folder, filename)
		req, err := grab.NewRequest(dst, urls[i])
		if err != nil {
			return err
		}
		reqs = append(reqs, req)
	}
	respCh := client.DoBatch(4, reqs...)

	for resp := range respCh {
		url := resp.Request.URL().String()
		filename := hlsFilename(mustParseURL(url))
		dst := filepath.Join(folder, filename)
		if err := resp.Err(); err != nil {
			ps.Pub(DownloadStatus{URL: url, TempFilename: dst, Progress: fmt.Sprintf("%v", resp.Progress()), Status: "error", Error: err.Error()}, DownloadStatusChannel)
			return err
		}
		ds := DownloadStatus{URL: url, TempFilename: dst, Progress: fmt.Sprintf("%v", resp.Progress()), Status: "done", Error: ""}
		completeSegmentDownload(&ds)
		ps.Pub(ds, DownloadStatusChannel)
	}
	return nil
}

// DownloadHLSPlaylist download an HLS playlist
func DownloadHLSPlaylist(url, storage string, ps *pubsub.PubSub) error {
	client := grab.NewClient()
	sourceURL := mustParseURL(url)
	filename := hlsFilename(sourceURL)
	content, err := DownloadHLSURL(sourceURL, filename, storage, ps)
	if err != nil {
		log.Printf("DownloadHLSPlaylist %v", err)
		return err
	}
	urls := GetSegmentURLS(content)
	return DownloadSegmentURLs(urls, storage, ps, client)
}

// IsHSLPlaylistDownloaded checks if an hls file has downloaded
func IsHSLPlaylistDownloaded(url, folder string) bool {
	sourceURL := mustParseURL(url)
	filename := hlsFilename(sourceURL)
	dst := filepath.Join(folder, filename)
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		return false
	}
	data, err := ioutil.ReadFile(dst)
	if err != nil {
		return false
	}
	urls := GetSegmentURLS(data)
	for _, url := range urls {
		if !segmentExists(url, folder) {
			return false
		}
	}
	return true
}
