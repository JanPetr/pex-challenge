package app

import (
	"fmt"
	"image"
	"sync"

	"github.com/EdlinOrg/prominentcolor"
	"github.com/janpetr/pex-challenge/pkg/logger"
	"github.com/janpetr/pex-challenge/pkg/metric"
)

type AnalyzedImage struct {
	url    string
	colors []string
}

func AnalyzeImages(images <-chan DownloadedImage) <-chan AnalyzedImage {
	parsedImages := make(chan AnalyzedImage)

	go func() {
		var wg sync.WaitGroup

		for i := range images {
			wg.Add(1)
			go func(i DownloadedImage) {
				defer wg.Done()

				colors, err := Analyze(i.image)
				if err != nil {
					// Report the error, but do not stop the pipeline.
					logger.Errorf(err, "failed image analysis, image URL: %q", i.url)
					metric.AddInt64(FailedAnalysisCnt, 1)

					return
				}

				p := AnalyzedImage{
					url:    i.url,
					colors: colors,
				}

				parsedImages <- p
				metric.AddInt64(AnalyzedCnt, 1)
			}(i)
		}
		wg.Wait()

		close(parsedImages)
	}()

	return parsedImages
}

func Analyze(img image.Image) ([]string, error) {
	cols, err := prominentcolor.KmeansWithArgs(prominentcolor.ArgumentNoCropping, img)
	if err != nil {
		return []string{}, err
	}

	var colors []string
	for _, c := range cols {
		colors = append(colors, fmt.Sprintf("#%s", c.AsString()))
	}

	return colors, nil
}
