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
	// "github.com/osiloke/streaming/log"
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
				folder:    "./t/",
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
				storage:   "./t/",
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

func TestGetHLSURLSPath(t *testing.T) {
	type args struct {
		url              *url.URL
		folder           string
		segmentURLPrefix string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			"testdownloadURL",
			args{
				url: mustParseURL(
					"https://audio.udux.com/hls/0fa9a977f15c41508efe788b085751a5/eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NjE0NTg1NTIsImlwIjoiNDUuMjIyLjk4LjI0OSwxMDQuMTU0LjE2Mi4yMzciLCJyb2xlIjoiUzFmVHZtQWlmSXQyaVlRUTBCckIiLCJ0cmFja19pZCI6IkliUHVrUVpLUExRNENZU2VINzU4IiwiZGV2aWNlIjoiaW9zIiwiY291bnRyeSI6IlVTIiwiaWF0IjoxNTYxNDU0OTUyfQ.cuXW8o6liFeoCHvQLSIVSZCyxs_yjEU6hQKO9TiguAM/IbPukQZKPLQ4CYSeH758_trd_preview.mp4/index.m3u8?ut=st=1561454833~exp=1561458433~acl=/hls/IbPukQZKPLQ4CYSeH758_trd_preview.mp4/*~hmac=5a779364d449b9e06de80a3f947ee717e62d5ab6fc15c681daa795205761484b",
				),
				folder:           "./t/",
				segmentURLPrefix: "http://localhost:8888/cache?r=1&file=",
			},
			[]string{"t/73d65d39a95a80c5874d043e8d248fdee38a1959f611f1fd7e490024302e58c9",
				"t/5b1821643ca90a30fb5f31827521504252dc9366abd5b3b9a0221dd8af806e13",
				"t/2fba1f499213fb508b9d2cf342921b307cf5be0ec9b650f2b489790b4a6d9e3c",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetHLSURLSPath(tt.args.url, tt.args.folder, tt.args.segmentURLPrefix)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetHLSURLSPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetHLSURLSPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveHLSPlaylist(t *testing.T) {
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
			"remove",
			args{
				url:       "https://audio.udux.com/hls/0fa9a977f15c41508efe788b085751a5/eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NjE0NTg1NTIsImlwIjoiNDUuMjIyLjk4LjI0OSwxMDQuMTU0LjE2Mi4yMzciLCJyb2xlIjoiUzFmVHZtQWlmSXQyaVlRUTBCckIiLCJ0cmFja19pZCI6IkliUHVrUVpLUExRNENZU2VINzU4IiwiZGV2aWNlIjoiaW9zIiwiY291bnRyeSI6IlVTIiwiaWF0IjoxNTYxNDU0OTUyfQ.cuXW8o6liFeoCHvQLSIVSZCyxs_yjEU6hQKO9TiguAM/IbPukQZKPLQ4CYSeH758_trd_preview.mp4/index.m3u8?ut=st=1561454833~exp=1561458433~acl=/hls/IbPukQZKPLQ4CYSeH758_trd_preview.mp4/*~hmac=5a779364d449b9e06de80a3f947ee717e62d5ab6fc15c681daa795205761484b",
				storage:   "./t/",
				proxyAddr: "http://localhost:8888/cache?r=1&file=",
				ps:        ps,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RemoveHLSPlaylist(tt.args.url, tt.args.storage, tt.args.proxyAddr, tt.args.ps); (err != nil) != tt.wantErr {
				t.Errorf("RemoveHLSPlaylist() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
