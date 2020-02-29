package app

import (
	"testing"

	"github.com/janpetr/pex-challenge/pkg/metric"
	"github.com/stretchr/testify/require"
)

// Integration test - requires internet connection to make HTTP requests
// Possible improvement - inject HTTP client and mock it in tests
func TestDownloadImages(t *testing.T) {
	rps := 5

	tests := []struct {
		name          string
		urls          []string
		wantSuccesses int64
		wantFails     int64
	}{
		{
			name:          "single URL",
			urls:          []string{"http://i.imgur.com/FApqk3D.jpg"},
			wantSuccesses: 1,
			wantFails:     0,
		},
		{
			name:          "multiple URLs",
			urls:          []string{"http://i.imgur.com/FApqk3D.jpg", "http://i.imgur.com/TKLs9lo.jpg", "https://i.redd.it/d8021b5i2moy.jpg"},
			wantSuccesses: 3,
			wantFails:     0,
		},
		{
			name:          "404 URL",
			urls:          []string{"http://i.foo.com/bar.jpg"},
			wantSuccesses: 0,
			wantFails:     1,
		},
		{
			name:          "404 URL + success URLs",
			urls:          []string{"http://i.foo.com/bar.jpg", "http://i.imgur.com/TKLs9lo.jpg", "https://i.redd.it/d8021b5i2moy.jpg"},
			wantSuccesses: 2,
			wantFails:     1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metric.SetInt64(DownloadedCnt, 0)
			metric.SetInt64(FailedDownloadCnt, 0)

			urls := make(chan string, len(tt.urls))
			for _, u := range tt.urls {
				urls <- u
			}
			close(urls)

			dis := DownloadImages(urls, rps)

			var downloadedImages []DownloadedImage
			for di := range dis {
				downloadedImages = append(downloadedImages, di)
			}

			require.Equal(t, tt.wantSuccesses, metric.GetInt64(DownloadedCnt))
			require.Equal(t, tt.wantFails, metric.GetInt64(FailedDownloadCnt))

			for _, dImage := range downloadedImages {
				require.Contains(t, tt.urls, dImage.url)
				require.NotNil(t, dImage.image)
			}
		})
	}
}
