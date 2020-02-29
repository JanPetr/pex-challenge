package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/janpetr/pex-challenge/app"
	"github.com/janpetr/pex-challenge/pkg/metric"
)

var inputFile = flag.String("inputFile", "", "path to input file with image URLs")
var outputFile = flag.String("outputFile", "output.csv", "path to output file where the analyzed images are dumped")
var rps = flag.Int("rps", 20, "max requests per second")
var processDuplicateURLs = flag.Bool("processDuplicateURLs", false, `When set to 'false' the program won't process images from duplicate URLs and process them only once. For larged URL sets with a lot of unique URLs you might need to set it to 'true' as it might have a big memory impact.`)

func main() {
	start := time.Now()

	// Prepare arguments
	flag.Parse()

	if len(*inputFile) <= 0 {
		fmt.Println("Input file is not specified, please specify it")
		return
	}

	// Run
	fmt.Println("The program is running ...")

	urls, errr := app.ReadURLs(*inputFile, *processDuplicateURLs)
	images := app.DownloadImages(urls, *rps)
	parsedImages := app.AnalyzeImages(images)
	done, erre := app.ExportCSV(*outputFile, parsedImages)

	select {
	case err := <-errr:
		fmt.Println("Reading input file failed:", err)
		return
	case err := <-erre:
		fmt.Println("Writing output file failed:", err)
		return
	case <-done:
		fmt.Println("\nThe program has finished:")
		fmt.Printf("- URLs read: %d\n", metric.GetInt64(app.ReadCnt))
		fmt.Printf("- URLs forwarded to processing: %d\n", metric.GetInt64(app.ForwardedCnt))
		fmt.Printf("- Images downloaded: %d\n", metric.GetInt64(app.DownloadedCnt))
		fmt.Printf("- Downloads failed: %d\n", metric.GetInt64(app.FailedDownloadCnt))
		fmt.Printf("- Images analyzed: %d\n", metric.GetInt64(app.AnalyzedCnt))
		fmt.Printf("- Analysis failed: %d\n\n", metric.GetInt64(app.FailedAnalysisCnt))

		fmt.Printf("Totally processed %d images in %s\n", metric.GetInt64(app.ExportedCnt), time.Since(start))
		return
	}
}
