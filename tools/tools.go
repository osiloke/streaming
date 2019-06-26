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
	log.Debug.Println("done")
	// }()
	return "done"
}
