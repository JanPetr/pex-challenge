package app

import (
	"bufio"
	"os"
	"sync"

	"github.com/janpetr/pex-challenge/pkg/metric"
)

var readURLs sync.Map

func ReadURLs(fileName string, processDuplicateURLs bool) (<-chan string, <-chan error) {
	urls := make(chan string)
	errors := make(chan error, 1)

	file, err := os.Open(fileName)
	if err != nil {
		errors <- err
	}

	go func() {
		defer func() {
			close(urls)
			err := file.Close()
			if err != nil {
				errors <- err
			}
		}()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			url := scanner.Text()
			metric.AddInt64(ReadCnt, 1)

			// Check duplicate URLs
			if !processDuplicateURLs {
				if _, ok := readURLs.Load(url); ok {
					continue
				}

				readURLs.Store(url, struct{}{})
			}

			urls <- url
			metric.AddInt64(ForwardedCnt, 1)
		}
	}()

	return urls, errors
}
