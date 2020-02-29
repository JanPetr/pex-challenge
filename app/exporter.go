package app

import (
	"fmt"
	"os"
	"strings"

	"github.com/janpetr/pex-challenge/pkg/metric"
)

func ExportCSV(outputFile string, parsedImages <-chan AnalyzedImage) (<-chan struct{}, <-chan error) {
	done := make(chan struct{})
	errors := make(chan error, 1)

	go func() {
		f, err := os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
		if err != nil {
			errors <- err
		}

		defer func() {
			err = f.Close()
			if err != nil {
				errors <- err
			}

			close(done)
		}()

		for pi := range parsedImages {
			_, err := f.WriteString(fmt.Sprintf("%s,%s\n", pi.url, strings.Join(pi.colors, ",")))
			if err != nil {
				errors <- err
			}

			metric.AddInt64(ExportedCnt, 1)
		}
	}()

	return done, errors
}
