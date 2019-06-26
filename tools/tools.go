package tools

import (
	"encoding/json"

	"github.com/cskr/pubsub"
	"github.com/osiloke/streaming/downloader"
)

type EventBus interface {
	SendMessageEvent(channel, message string)
}

func GetHLS(url, storage string, dispatcher EventBus) string {
	ps := pubsub.New(1)
	ch := ps.Sub(downloader.DownloadStatusChannel)
	go func() {
		for c := range ch {
			status := c.(downloader.DownloadStatus)
			v, _ := json.Marshal(status)
			dispatcher.SendMessageEvent("DOWNLOAD_STATUS", string(v))
		}
	}()
	downloader.DownloadHLSPlaylist(url, storage, ps)
	return "done"
}
