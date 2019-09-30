package downloader

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cavaliercoder/grab"
	"github.com/cskr/pubsub"
	"github.com/osiloke/streaming/log"
)

//EventBus sends events to a listener
type EventBus interface {
	SendEvent(channel, message string)
}

// DownloadStatusChannel is the pubsub channel
const DownloadStatusChannel = "download_status"

// RemoveStatusChannel is the pubsub channel for remove status
const RemoveStatusChannel = "remove_status"

// DownloadStatus represents the status of a download
type DownloadStatus struct {
	URL          string `json:"url"`
	ID           string `json:"id"`
	Segment      string `json:"segment"`
	TempFilename string `json:"tempFilename"`
	Prefix       string `json:"prefix"`
	Progress     string `json:"progress"`
	Status       string `json:"status"`
	Error        string `json:"error"`
}

// RemoveStatus status sent while removing an HLS url from cache
type RemoveStatus struct {
	URL          string `json:"url"`
	ID           string `json:"id"`
	Segment      string `json:"segment"`
	TempFilename string `json:"tempFilename"`
	Prefix       string `json:"prefix"`
	Progress     string `json:"progress"`
	Status       string `json:"status"`
	Error        string `json:"error"`
}

func mustParseURL(urlSt string) *url.URL {
	u, _ := url.Parse(urlSt)
	return u
}

// GetHLSURLSPath get all paths related to an HLS url
func GetHLSURLSPath(url *url.URL, folder, segmentURLPrefix string) ([]string, error) {
	start := time.Now()
	hlsFilename := PrefixedHlsFilename(segmentURLPrefix, url)
	dst := filepath.Join(folder, hlsFilename)
	defer func() {
		elapsed := time.Since(start)
		log.Debug.Printf("GetHLSURLSPath - elapsed - %s", elapsed)
	}()
	f, err := os.Open(dst)
	if err != nil {
		return nil, err
	}
	hlsBody, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	urls := []string{dst}
	for _, match := range re.FindAllString(string(hlsBody), -1) {
		url := PrefixedHlsFilename(segmentURLPrefix, mustParseURL(match))
		urls = append(urls, filepath.Join(folder, url))
	}
	return urls, nil
}

// GetHLSSegments get segments un HLS file
func GetHLSSegments(url *url.URL, folder, segmentURLPrefix string) ([]string, error) {
	start := time.Now()
	hlsFilename := PrefixedHlsFilename(segmentURLPrefix, url)
	dst := filepath.Join(folder, hlsFilename)
	defer func() {
		elapsed := time.Since(start)
		log.Debug.Printf("GetHLSSegments - elapsed - %s", elapsed)
	}()
	f, err := os.Open(dst)
	if err != nil {
		return nil, err
	}
	hlsBody, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	urls := []string{dst}
	for _, match := range re.FindAllString(string(hlsBody), -1) {
		url := mustParseURL(match)
		urls = append(urls, filepath.Join(folder, hashKey(url.String())))
	}
	return urls, nil
}

// DownloadHLSURL download a url to a file
func DownloadHLSURL(url *url.URL, filename, folder, segmentURLPrefix string, ps *pubsub.PubSub) ([]byte, error) {
	start := time.Now()
	// done := make(chan int64)
	resp, err := http.Get(url.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("unable to retrieve hls")
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
	log.Debug.Printf("downloaded %s to %s - %s", url, dst, elapsed)

	return []byte(strings.TrimSpace(string(body))), err
}

// DownloadSegmentURLs takes an array of urls to be downloaded
func DownloadSegmentURLs(urls []string, folder, segmentURLPrefix string, ps *pubsub.PubSub, client *grab.Client) error {
	reqs := make([]*grab.Request, 0)
	for i := 0; i < len(urls); i++ {
		filename := PrefixedHlsFilename(segmentURLPrefix, mustParseURL(urls[i]))
		dst := filepath.Join(folder, filename)
		if _, err := os.Stat(dst); !os.IsNotExist(err) {
			continue
		}
		req, err := grab.NewRequest(dst, urls[i])
		if err != nil {
			return err
		}
		reqs = append(reqs, req)
	}
	respCh := client.DoBatch(4, reqs...)

	for resp := range respCh {
		idf := idAndFile(resp.Request.URL())
		url := resp.Request.URL().String()
		filename := PrefixedHlsFilename(segmentURLPrefix, mustParseURL(url))
		dst := filepath.Join(folder, filename)
		if err := resp.Err(); err != nil {
			os.Remove(filename)
			ps.Pub(DownloadStatus{URL: url, ID: idf[0], Segment: idf[1], Prefix: segmentURLPrefix, TempFilename: dst, Progress: fmt.Sprintf("%v", resp.Progress()), Status: "error downloading segment", Error: err.Error()}, DownloadStatusChannel)
			return err
		}
		ds := DownloadStatus{URL: url, Prefix: segmentURLPrefix, TempFilename: dst, Progress: fmt.Sprintf("%v", resp.Progress()), Status: "downloaded segment", Error: ""}
		completeSegmentDownload(&ds)
		ps.Pub(ds, DownloadStatusChannel)
	}
	log.Debug.Printf("Downloaded %v segments\n", len(reqs))
	return nil
}

// DownloadHLSPlaylist download an HLS playlist
func DownloadHLSPlaylist(url, storage, segmentURLPrefix string, ps *pubsub.PubSub) error {
	client := grab.NewClient()
	sourceURL := mustParseURL(url)
	idf := idAndFile(sourceURL)
	filename := PrefixedHlsFilename(segmentURLPrefix, sourceURL)
	content, err := DownloadHLSURL(sourceURL, filename, storage, segmentURLPrefix, ps)
	if err != nil {
		ds := DownloadStatus{URL: url, ID: idf[0], Segment: idf[1], Prefix: segmentURLPrefix, TempFilename: filename, Progress: "0", Status: "failed hls", Error: err.Error()}
		ps.Pub(ds, DownloadStatusChannel)

		log.Debug.Printf("DownloadHLSPlaylist %v", err)
		return err
	}
	ds := DownloadStatus{URL: url, ID: idf[0], Segment: idf[1], Prefix: segmentURLPrefix, TempFilename: filename, Progress: "1", Status: "downloaded index", Error: ""}
	ps.Pub(ds, DownloadStatusChannel)
	urls := GetSegmentURLS(content, segmentURLPrefix)
	err = DownloadSegmentURLs(urls, storage, segmentURLPrefix, ps, client)
	if err != nil {
		ds = DownloadStatus{URL: url, ID: idf[0], Segment: idf[1], Prefix: segmentURLPrefix, TempFilename: filename, Progress: "0", Status: "failed hls", Error: err.Error()}
		ps.Pub(ds, DownloadStatusChannel)

		return err
	}

	ds = DownloadStatus{URL: url, ID: idf[0], Segment: idf[1], Prefix: segmentURLPrefix, TempFilename: filename, Progress: "1", Status: "downloaded hls", Error: ""}
	ps.Pub(ds, DownloadStatusChannel)
	return nil
}

// RemoveHLSPlaylist removes a cached HLS playlist
func RemoveHLSPlaylist(url, storage, segmentURLPrefix string, ps *pubsub.PubSub) error {
	urls, err := GetHLSSegments(mustParseURL(url), storage, segmentURLPrefix)
	if err != nil {
		if !strings.Contains(err.Error(), "no such file or directory") {
			return err
		}
		return nil
	}
	for _, url := range urls[1:] {
		if err := os.Remove(url); err != nil {
			if !strings.Contains(err.Error(), "no such file or directory") {
				ds := RemoveStatus{URL: url, Prefix: segmentURLPrefix, TempFilename: "", Progress: "1", Status: "failed segment", Error: err.Error()}
				ps.Pub(ds, RemoveStatusChannel)
				return err
			}
		}
		ds := RemoveStatus{URL: url, Prefix: segmentURLPrefix, TempFilename: "", Progress: "1", Status: "remove segment", Error: ""}
		ps.Pub(ds, RemoveStatusChannel)
	}
	err = os.Remove(urls[0])
	if err != nil {
		if !strings.Contains(err.Error(), "no such file or directory") {
			ds := RemoveStatus{URL: url, Prefix: segmentURLPrefix, TempFilename: "", Progress: "1", Status: "failed segment", Error: err.Error()}
			ps.Pub(ds, RemoveStatusChannel)
			return err
		}
	}
	ds := RemoveStatus{URL: urls[0], Prefix: segmentURLPrefix, TempFilename: "", Progress: "1", Status: "remove index", Error: ""}
	ps.Pub(ds, RemoveStatusChannel)
	return nil
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
