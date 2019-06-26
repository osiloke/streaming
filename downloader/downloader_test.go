package downloader

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/cskr/pubsub"
)

func hlsResponse() []byte {
	return []byte(`#EXTM3U
#EXT-X-TARGETDURATION:9
#EXT-X-ALLOW-CACHE:YES
#EXT-X-PLAYLIST-TYPE:VOD
#EXT-X-VERSION:3
#EXT-X-MEDIA-SEQUENCE:1
#EXTINF:9.009,
http://localhost:8888/cache?r=1&file=https://audio.udux.com/hls/0fa9a977f15c41508efe788b085751a5/eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NjE0NTg1NTIsImlwIjoiNDUuMjIyLjk4LjI0OSwxMDQuMTU0LjE2Mi4yMzciLCJyb2xlIjoiUzFmVHZtQWlmSXQyaVlRUTBCckIiLCJ0cmFja19pZCI6IkliUHVrUVpLUExRNENZU2VINzU4IiwiZGV2aWNlIjoiaW9zIiwiY291bnRyeSI6IlVTIiwiaWF0IjoxNTYxNDU0OTUyfQ.cuXW8o6liFeoCHvQLSIVSZCyxs_yjEU6hQKO9TiguAM/IbPukQZKPLQ4CYSeH758_trd_preview.mp4/segment-1-a1.ts
#EXTINF:1.150,
http://localhost:8888/cache?r=1&file=https://audio.udux.com/hls/0fa9a977f15c41508efe788b085751a5/eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NjE0NTg1NTIsImlwIjoiNDUuMjIyLjk4LjI0OSwxMDQuMTU0LjE2Mi4yMzciLCJyb2xlIjoiUzFmVHZtQWlmSXQyaVlRUTBCckIiLCJ0cmFja19pZCI6IkliUHVrUVpLUExRNENZU2VINzU4IiwiZGV2aWNlIjoiaW9zIiwiY291bnRyeSI6IlVTIiwiaWF0IjoxNTYxNDU0OTUyfQ.cuXW8o6liFeoCHvQLSIVSZCyxs_yjEU6hQKO9TiguAM/IbPukQZKPLQ4CYSeH758_trd_preview.mp4/segment-2-a1.ts
#EXT-X-ENDLIST`)
}

func TestDownloadHLSURL(t *testing.T) {
	ps := pubsub.New(1)
	ch := ps.Sub(DownloadStatusChannel)
	go func() {
		for v := range ch {
			log.Println(fmt.Sprintf("%+v", v))
		}
	}()
	type args struct {
		url       *url.URL
		filename  string
		folder    string
		proxyAddr string
		ps        *pubsub.PubSub
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			"testdownloadURL",
			args{
				url: mustParseURL(
					"https://audio.udux.com/hls/0fa9a977f15c41508efe788b085751a5/eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NjE0NTg1NTIsImlwIjoiNDUuMjIyLjk4LjI0OSwxMDQuMTU0LjE2Mi4yMzciLCJyb2xlIjoiUzFmVHZtQWlmSXQyaVlRUTBCckIiLCJ0cmFja19pZCI6IkliUHVrUVpLUExRNENZU2VINzU4IiwiZGV2aWNlIjoiaW9zIiwiY291bnRyeSI6IlVTIiwiaWF0IjoxNTYxNDU0OTUyfQ.cuXW8o6liFeoCHvQLSIVSZCyxs_yjEU6hQKO9TiguAM/IbPukQZKPLQ4CYSeH758_trd_preview.mp4/index.m3u8?ut=st=1561454833~exp=1561458433~acl=/hls/IbPukQZKPLQ4CYSeH758_trd_preview.mp4/*~hmac=5a779364d449b9e06de80a3f947ee717e62d5ab6fc15c681daa795205761484b",
				),
				filename:  "test.m3u8",
				folder:    "./",
				proxyAddr: "http://localhost:8888/cache?r=1&file=",
				ps:        ps,
			},
			hlsResponse(),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Running with %s args - %+v", tt.name, tt.args)
			got, err := DownloadHLSURL(tt.args.url, tt.args.filename, tt.args.folder, tt.args.proxyAddr, tt.args.ps)
			if (err != nil) != tt.wantErr {
				t.Errorf("DownloadHLSURL() error = %v, wantErr %v", err, tt.wantErr)
				os.Remove(filepath.Join(tt.args.folder, tt.args.filename))
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DownloadHLSURL() = %v, want %v", got, tt.want)
			}
			os.Remove(filepath.Join(tt.args.folder, tt.args.filename))
		})
	}
}

func TestDownloadHLSPlaylist(t *testing.T) {
	ps := pubsub.New(1)
	type args struct {
		url       string
		storage   string
		proxyAddr string
		ps        *pubsub.PubSub
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"download",
			args{
				url:       "https://audio.udux.com/hls/0fa9a977f15c41508efe788b085751a5/eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NjE0NTg1NTIsImlwIjoiNDUuMjIyLjk4LjI0OSwxMDQuMTU0LjE2Mi4yMzciLCJyb2xlIjoiUzFmVHZtQWlmSXQyaVlRUTBCckIiLCJ0cmFja19pZCI6IkliUHVrUVpLUExRNENZU2VINzU4IiwiZGV2aWNlIjoiaW9zIiwiY291bnRyeSI6IlVTIiwiaWF0IjoxNTYxNDU0OTUyfQ.cuXW8o6liFeoCHvQLSIVSZCyxs_yjEU6hQKO9TiguAM/IbPukQZKPLQ4CYSeH758_trd_preview.mp4/index.m3u8?ut=st=1561454833~exp=1561458433~acl=/hls/IbPukQZKPLQ4CYSeH758_trd_preview.mp4/*~hmac=5a779364d449b9e06de80a3f947ee717e62d5ab6fc15c681daa795205761484b",
				storage:   "./",
				proxyAddr: "http://localhost:8888/cache?r=1&file=",
				ps:        ps,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DownloadHLSPlaylist(tt.args.url, tt.args.storage, tt.args.proxyAddr, tt.args.ps); (err != nil) != tt.wantErr {
				t.Errorf("DownloadHLSPlaylist() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
