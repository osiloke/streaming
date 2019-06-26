package downloader

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
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
	Prefix       string `json:"prefix"`
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

// DownloadHLSURL download a url to a file
func DownloadHLSURL(url *url.URL, filename, folder, segmentURLPrefix string, ps *pubsub.PubSub) ([]byte, error) {
	start := time.Now()
	// done := make(chan int64)
	resp, err := http.Get(url.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	dst := filepath.Join(folder, filename)
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	body := buf.Bytes()
	proxiedBody, _ := ProxyHLSUrls(body, segmentURLPrefix)
	err = ioutil.WriteFile(dst, proxiedBody, 0644)
	if err != nil {
		return nil, err
	}
	// done <- n
	elapsed := time.Since(start)
	log.Printf("downloaded %s to %s - %s", url, dst, elapsed)

	return []byte(strings.TrimSpace(string(proxiedBody))), err
}

// DownloadSegmentURLs takes an array of urls to be downloaded
func DownloadSegmentURLs(urls []string, folder, segmentURLPrefix string, ps *pubsub.PubSub, client *grab.Client) error {
	reqs := make([]*grab.Request, 0)
	for i := 0; i < len(urls); i++ {
		filename := prefixedHlsFilename(segmentURLPrefix, mustParseURL(urls[i]))
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
		filename := prefixedHlsFilename(segmentURLPrefix, mustParseURL(url))
		dst := filepath.Join(folder, filename)
		if err := resp.Err(); err != nil {
			ps.Pub(DownloadStatus{URL: url, Prefix: segmentURLPrefix, TempFilename: dst, Progress: fmt.Sprintf("%v", resp.Progress()), Status: "error", Error: err.Error()}, DownloadStatusChannel)
			return err
		}
		ds := DownloadStatus{URL: url, Prefix: segmentURLPrefix, TempFilename: dst, Progress: fmt.Sprintf("%v", resp.Progress()), Status: "done", Error: ""}
		completeSegmentDownload(&ds)
		ps.Pub(ds, DownloadStatusChannel)
	}
	return nil
}

// DownloadHLSPlaylist download an HLS playlist
func DownloadHLSPlaylist(url, storage, segmentURLPrefix string, ps *pubsub.PubSub) error {
	client := grab.NewClient()
	sourceURL := mustParseURL(url)
	filename := prefixedHlsFilename(segmentURLPrefix, sourceURL)
	content, err := DownloadHLSURL(sourceURL, filename, storage, segmentURLPrefix, ps)
	if err != nil {
		log.Printf("DownloadHLSPlaylist %v", err)
		return err
	}
	urls := GetSegmentURLS(content, segmentURLPrefix)
	return DownloadSegmentURLs(urls, storage, segmentURLPrefix, ps, client)
}

// IsHSLPlaylistDownloaded checks if an hls file has downloaded
func IsHSLPlaylistDownloaded(url, folder, segmentURLPrefix string) bool {
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
	urls := GetSegmentURLS(data, segmentURLPrefix)
	for _, url := range urls {
		if !segmentExists(url, folder) {
			return false
		}
	}
	return true
}
