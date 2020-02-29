package app

import (
	"bufio"
	"fmt"
	"os"

	"github.com/janpetr/pex-challenge/pkg/metric"
)

// Possible improvement - randomize the file name in order to let multiple programs run at once
var tmpFile = "/tmp/processed-urls-by-pex-challenge"

func ReadURLs(fileName string) (<-chan string, <-chan error) {
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

			err = os.Remove(tmpFile)
			if err != nil {
				errors <- err
			}
		}()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			url := scanner.Text()
			metric.AddInt64(ReadCnt, 1)

			isAllowed, err := isAllowed(url)
			if err != nil {
				errors <- err
			}

			if isAllowed {
				urls <- url
				metric.AddInt64(ForwardedCnt, 1)
			}
		}
	}()

	return urls, errors
}

func isAllowed(url string) (bool, error) {
	tmpFile, err := os.OpenFile(tmpFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return false, err
	}

	defer tmpFile.Close()

	scanner := bufio.NewScanner(tmpFile)
	for scanner.Scan() {
		processedUrl := scanner.Text()
		if processedUrl == url {
			return false, nil
		}
	}

	_, err = tmpFile.WriteString(fmt.Sprintf("%s\n", url))
	if err != nil {
		return false, err
	}

	return true, nil
}
