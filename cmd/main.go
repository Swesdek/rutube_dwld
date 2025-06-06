package main

import (
	"flag"
	"log"

	"github.com/Swesdek/rutube-dwld/internal/download"
	"github.com/Swesdek/rutube-dwld/internal/platforms/rutube"
)

func main() {
	threads := flag.Int("threads", 10, "Number of threads for downloading")
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		log.Fatal("Error: url was not supplied")
	}
	videoName, segments, segmentsAmount, rawDownloadUrl := rutube.GetVideoInfo(args[0])

	download.Download(*threads, videoName, segments, segmentsAmount, rawDownloadUrl)
}
