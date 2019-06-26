package downloader

import (
	"net/url"
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

func Test_PrefixedHlsFilename(t *testing.T) {
	type args struct {
		prefix string
		url    *url.URL
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"test",
			args{
				prefix: "http://127.0.0.1:7071/cache?r=1&file=",
				url: mustParseURL(
					"https://audio.udux.com/hls/8b3acdbe-7b5c-4a7e-ae55-aaa61ccd3cf8/eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NjE1OTE2MDcsImlwIjoiMTI5LjU2LjM3LjEzNywxMDQuMTU1LjE1NC4yNDgiLCJ1c2VyX2lkIjoiYUJQc3NvcDFkZ1o3TVh3bkd1T092SlFzcFZ0MiIsInN1YnNjcmlwdGlvbiI6IlJjWnluMlNkN0VHaURmNklnYmVHIiwicGxhbiI6Ilk0YThyMnYydUNGR2FLT1BYMHY0Iiwicm9sZSI6Ik5vVkpRWHh4azdIb0FzOEtlTDM1IiwidHJhY2tfaWQiOiJ0TWk0b3d6VE1QaVJNQVNWMWxkSiIsImRldmljZSI6ImlvcyIsImNvdW50cnkiOiJVUyIsImlhdCI6MTU2MTU4ODAwN30.fcSRQyh3z4vSj0yCIrlO9Tp32YMvfkanbdhvt3i5zOY/tMi4owzTMPiRMASV1ldJ_trd.mp4/segment-1-a1.ts",
				),
			},
			"b5724af1fb4efc7723e4293faf14ed702cea91da95a10f7abdf7b6b12c832511",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PrefixedHlsFilename(tt.args.prefix, tt.args.url); got != tt.want {
				t.Errorf("PrefixedHlsFilename() = %v, want %v", got, tt.want)
			}
		})
	}
}
