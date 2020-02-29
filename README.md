# Pex Coding Challenge

This repo implements a tool to read and detect 3 most commin colors from images from provided list of URLs and dumps results to CSV file.
Original assignment: https://gist.github.com/ehmo/e736c827ca73d84581d812b3a27bb132

## How to run it

1. Setup your golang environment (https://golang.org/doc/) and make sure to have your GOPATH correctly exported.
2. [Go modules](https://blog.golang.org/using-go-modules) are used for dependency management so ensure to install Go 1.13.6 and update your local environment with
    ```
    $ export GO111MODULE=on
    $ export GOPROXY="https://proxy.golang.org,direct"
    ```
3. Clone this repository in `$GOPATH/src/github.com/janpetr`
4. Switch to the cloned directory:
    ```
    $ cd $GOPATH/src/github.com/janpetr/pex-challenge
    ```
4. Run the program with:
    ```
    $ go run . --inputFile input.txt
    ```
5. Done! You should see output similar to this one:
    ```
    The program is running ...
    2020/02/29 20:02:23 failed image download, image URL: ["https://i.redd.it/nrafqoujmety.jpg"]: 404 code returned, status: 404 Not Found
    
    The program has finished:
    - URLs read: 1000
    - URLs forwarded to processing: 40
    - Images downloaded: 39
    - Downloads failed: 1
    - Images analyzed: 39
    - Analysis failed: 0
    
    Totally processed 39 images in 3.801464214s
    ```
   
## Configuration of the program

The program accepts these flags:

### `--inputFile`

Specifies a file from which the program will read URLs. Each URL needs to be on a single line in the file.

**This parameter is mandatory.**

### `--outputFile`

Specifies a path and a name of output CSV file. 

_Defaults to `output.csv`_

### `--rps`

In order not to overload 3rd party services, the program implements rate limiting of outgoing requests.
This parameter specifies the number of requests per second shot to services providing the images.

_Defaults to `20`._

## Structure of the program

### `main.go`

File where it all begins. Parses CLI arguments, starts the pipeline, outputs results to CLI.

### `./app`

Package of the program itself. Contains 4 main files:
- `reader.go`
- `downloader.go`
- `analyzer.go`
- `exporter.go`

with its associated test files. Each file contains it's self-descriptive logic.

### `./pkg`

Contains packages / libraries with a potential to be used in broader context and with multiple services. 
Now it has dummy implementations of `logger` and `metrics` packages.

## Tests

Run the tests from `app/` directory: `$ go test .`.

The `downloader_test.go` contains and integration test and requires network connection to actually download the images.
Possible improvement would be to inject mocked HTTP client to not to rely on the connection.

The `analyzer.go` has a benchmark test for the image analysis.

## Flow

First of all the program parses the CLI arguments and starts the pipeline.

Pipeline starts by reading URLs from input file. 
Original [input file](https://gist.github.com/ehmo/e736c827ca73d84581d812b3a27bb132#file-input-txt) contains only 40
unique URLs out of 1000 URLs present in the file. Not processing duplicate images speeds up the program significantly.
As the program should be run on the limited resources and compromise on speed, I didn't want to hold the already processed
URLs in memory - for a large URL sets (billions of URLs) it might overflow the memory.
So I decided to store the URLs in file. It adds additional IO time, but spares the memory.

The read URLs are pushed through the channel to next step - `downloader`.

Downloader asynchronously downloads images from provided URLs, decodes them to `image.Image` and pushes them through 
another channel to `analyzer`.

Analyzer uses https://github.com/EdlinOrg/prominentcolor library to detect most common colors. 
The detection is done asynchronously and when done, the results are pushed through channel to `exporter`.

Exporter consumes the channel of analyzed images and dumps them to CSV file. 

## Colors detection

For the colors detection the program relies on a 3rd party library: https://github.com/EdlinOrg/prominentcolor

I hope it's not an issue as it's mentioned the algorithm itself is not important.

## Enjoy!

If you should have any questions, hit me on hey@jan-petr.com 🚀
