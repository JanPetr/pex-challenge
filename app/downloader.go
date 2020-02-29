package app

import (
	"fmt"
	"image"
	"time"

	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"sync"

	"github.com/janpetr/pex-challenge/pkg/logger"
	"github.com/janpetr/pex-challenge/pkg/metric"
)

type DownloadedImage struct {
	url   string
	image image.Image
}

func DownloadImages(urls <-chan string, requestsPerSecond int) <-chan DownloadedImage {
	images := make(chan DownloadedImage)

	go func() {
		var wg sync.WaitGroup

		rate := time.Second / time.Duration(requestsPerSecond)
		throttle := time.Tick(rate)

		for u := range urls {
			wg.Add(1)
			<-throttle // rate limit requests
			go func(url string) {
				defer func() {
					wg.Done()
				}()

				i, err := getImage(url)
				if err != nil {
					// Report the error, but do not stop the pipeline.
					logger.Errorf(err, "failed image download, image URL: %q", url)
					metric.AddInt64(FailedDownloadCnt, 1)
					return
				}

				images <- DownloadedImage{
					url:   url,
					image: i,
				}

				metric.AddInt64(DownloadedCnt, 1)
			}(u)
		}

		wg.Wait()
		close(images)
	}()

	return images
}

func getImage(url string) (image.Image, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http get: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("%d code returned, status: %s", res.StatusCode, res.Status)
	}

	i, _, err := image.Decode(res.Body)
	if err != nil {
		return nil, fmt.Errorf("decode image: %w", err)
	}

	return i, nil
}
