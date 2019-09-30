package tools

import (
	"encoding/json"

	"github.com/cskr/pubsub"
	"github.com/osiloke/streaming/downloader"
	"github.com/osiloke/streaming/log"
)

// EventBus event sender
type EventBus interface {
	SendMessageEvent(channel, message string)
}

// GetHLS get hls and store locally
func GetHLS(url, storage, segmentURLPrefix string, dispatcher EventBus) string {
	ps := pubsub.New(1)
	ch := ps.Sub(downloader.DownloadStatusChannel)
	go downloader.DownloadHLSPlaylist(url, storage, segmentURLPrefix, ps)
	// go func() {
	for c := range ch {
		status := c.(downloader.DownloadStatus)
		v, _ := json.Marshal(status)
		dispatcher.SendMessageEvent("DOWNLOAD_STATUS", string(v))
	}
	log.Debug.Printf("Finished storing - %s", url)
	// }()
	return "done"
}

// GetMultipleHLS get urls
func GetMultipleHLS(urls []string, storage, segmentURLPrefix string, dispatcher EventBus) {
	// tools.GetHLS(url, storage, segmentURLPrefix, dispatcher)
	ps := pubsub.New(1)
	ch := ps.Sub(downloader.DownloadStatusChannel)
	go func() {
		for c := range ch {
			status := c.(downloader.DownloadStatus)
			v, _ := json.Marshal(status)
			dispatcher.SendMessageEvent("DOWNLOAD_STATUS", string(v))
		}
	}()
	for _, url := range urls {
		downloader.DownloadHLSPlaylist(url, storage, segmentURLPrefix, ps)
		log.Debug.Printf("Finished storing - %s", url)
	}
	ps.Unsub(ch, downloader.DownloadStatusChannel)
	// close(ch)
}

// RemoveHLS remove hls from local store
func RemoveHLS(url, storage, segmentURLPrefix string, dispatcher EventBus) string {
	ps := pubsub.New(1)
	ch := ps.Sub(downloader.RemoveStatusChannel)
	go downloader.RemoveHLSPlaylist(url, storage, segmentURLPrefix, ps)
	// go func() {
	for c := range ch {
		status := c.(downloader.RemoveStatus)
		v, _ := json.Marshal(status)
		dispatcher.SendMessageEvent("REMOVE_STATUS", string(v))
	}
	log.Debug.Printf("Finished removing - %s", url)
	// }()
	return "done"
}

// RemoveMultipleHLS remove hls urls from local store
func RemoveMultipleHLS(key string, urls []string, storage, segmentURLPrefix string, dispatcher EventBus) {
	// tools.GetHLS(url, storage, segmentURLPrefix, dispatcher)
	ps := pubsub.New(1)
	ch := ps.Sub(downloader.RemoveStatusChannel)
	dispatcher.SendMessageEvent("DOWNLOADER_REMOVE_START", key)
	go func() {
		for c := range ch {
			status := c.(downloader.RemoveStatus)
			v, _ := json.Marshal(status)
			dispatcher.SendMessageEvent("DOWNLOADER_REMOVE_STATUS", string(v))
		}
	}()
	for _, url := range urls {
		err := downloader.RemoveHLSPlaylist(url, storage, segmentURLPrefix, ps)
		if err != nil {
			log.Debug.Printf("Failed removing - %s - %s", url, err.Error())
		} else {
			log.Debug.Printf("Finished removing - %s", url)
		}
		dispatcher.SendMessageEvent("DOWNLOADER_REMOVE_HLS", url)
	}
	dispatcher.SendMessageEvent("DOWNLOADER_REMOVE_COMPLETE", key)
	ps.Unsub(ch, downloader.RemoveStatusChannel)
	// close(ch)
}
