package app

import (
	"image"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAnalyzeImages(t *testing.T) {
	testImagePaths := []string{
		"testdata/images/1nlgrn49x7ry.jpg",
		"testdata/images/4m5yk8gjrtzy.jpg",
		"testdata/images/d8021b5i2moy.jpg",
		"testdata/images/FApqk3D.jpg",
	}

	downloadedImages := make(chan DownloadedImage, len(testImagePaths))
	for _, ip := range testImagePaths {
		f, err := os.Open(ip)
		require.NoError(t, err)

		img, _, err := image.Decode(f)
		require.NoError(t, err)

		downloadedImages <- DownloadedImage{
			url:   ip,
			image: img,
		}
	}
	close(downloadedImages)

	ais := AnalyzeImages(downloadedImages)

	var analyzedImages []AnalyzedImage
	for ai := range ais {
		analyzedImages = append(analyzedImages, ai)
	}

	require.Equal(t, len(testImagePaths), len(analyzedImages))
	for _, aImage := range analyzedImages {
		require.Contains(t, testImagePaths, aImage.url)
		require.Len(t, aImage.colors, 3)
	}
}
