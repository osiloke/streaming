package downloader

import (
	"reflect"
	"testing"
)

func TestGetSegmentURLS(t *testing.T) {
	type args struct {
		hlsRaw           []byte
		segmentURLPrefix string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			"test",
			args{
				hlsRaw: []byte(`#EXTM3U
				#EXT-X-TARGETDURATION:9
				#EXT-X-ALLOW-CACHE:YES
				#EXT-X-PLAYLIST-TYPE:VOD
				#EXT-X-VERSION:3
				#EXT-X-MEDIA-SEQUENCE:1
				#EXTINF:9.009,
				http://localhost:8888/cache?r=1&file=https://audio.udux.com/hls/0fa9a977f15c41508efe788b085751a5/eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NjE0NTg1NTIsImlwIjoiNDUuMjIyLjk4LjI0OSwxMDQuMTU0LjE2Mi4yMzciLCJyb2xlIjoiUzFmVHZtQWlmSXQyaVlRUTBCckIiLCJ0cmFja19pZCI6IkliUHVrUVpLUExRNENZU2VINzU4IiwiZGV2aWNlIjoiaW9zIiwiY291bnRyeSI6IlVTIiwiaWF0IjoxNTYxNDU0OTUyfQ.cuXW8o6liFeoCHvQLSIVSZCyxs_yjEU6hQKO9TiguAM/IbPukQZKPLQ4CYSeH758_trd_preview.mp4/segment-1-a1.ts
				#EXTINF:1.150,
				http://localhost:8888/cache?r=1&file=https://audio.udux.com/hls/0fa9a977f15c41508efe788b085751a5/eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NjE0NTg1NTIsImlwIjoiNDUuMjIyLjk4LjI0OSwxMDQuMTU0LjE2Mi4yMzciLCJyb2xlIjoiUzFmVHZtQWlmSXQyaVlRUTBCckIiLCJ0cmFja19pZCI6IkliUHVrUVpLUExRNENZU2VINzU4IiwiZGV2aWNlIjoiaW9zIiwiY291bnRyeSI6IlVTIiwiaWF0IjoxNTYxNDU0OTUyfQ.cuXW8o6liFeoCHvQLSIVSZCyxs_yjEU6hQKO9TiguAM/IbPukQZKPLQ4CYSeH758_trd_preview.mp4/segment-2-a1.ts
				#EXT-X-ENDLIST
				`),
				segmentURLPrefix: "http://localhost:8888/cache?r=1&file=",
			},
			[]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetSegmentURLS(tt.args.hlsRaw, tt.args.segmentURLPrefix); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSegmentURLS() = %v, want %v", got, tt.want)
			}
		})
	}
}
