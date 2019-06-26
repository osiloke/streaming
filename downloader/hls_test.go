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
				prefix: "http://localhost:7071/cache?r=1&file=",
				url: mustParseURL(
					"https://audio.udux.com/hls/8b3acdbe-7b5c-4a7e-ae55-aaa61ccd3cf8/eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NjE1ODEwMzcsImlwIjoiMTI5LjU2LjM3LjEzNywxMDQuMTU0LjE2Mi4yMzciLCJ1c2VyX2lkIjoiYUJQc3NvcDFkZ1o3TVh3bkd1T092SlFzcFZ0MiIsInN1YnNjcmlwdGlvbiI6IlJjWnluMlNkN0VHaURmNklnYmVHIiwicGxhbiI6Ilk0YThyMnYydUNGR2FLT1BYMHY0Iiwicm9sZSI6Ik5vVkpRWHh4azdIb0FzOEtlTDM1IiwidHJhY2tfaWQiOiJ0TWk0b3d6VE1QaVJNQVNWMWxkSiIsImRldmljZSI6ImlvcyIsImNvdW50cnkiOiJVUyIsImlhdCI6MTU2MTU3NzQzN30.LsZOLYaKD_d_vPRkH5bqKzNCI8ccg16HaXnBFb6L9p4/tMi4owzTMPiRMASV1ldJ_trd.mp4/index.m3u8?ut=st=1561577317~exp=1561580917~acl=/hls/tMi4owzTMPiRMASV1ldJ_trd.mp4/*~hmac=995956a6e5ba53a91ca3b3adcd107d6e14b0fad31660c4b26b1282946b1b9ff2",
				),
			},
			"78e4a74ae6bf2b87ae5a9ee1cd9ef92c2d87b744c73a265cab5a0a6858944e35",
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