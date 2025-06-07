package download

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/grafov/m3u8"
	"github.com/schollz/progressbar/v3"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

type fileMetadata struct {
	path   string
	number uint
}

type urlMetadata struct {
	url    string
	number uint
}

func Download(
	threads int,
	videoName string,
	segments []*m3u8.MediaSegment,
	segmentsAmount uint,
	rawDownloadurl string,
) {
	err := os.Mkdir(videoName, 0700)
	if err != nil {
		panic(err)
	}

	workerContext := context.Background()
	segmentsDataChan := make(chan urlMetadata, threads)
	fileMetadataChan := make(chan fileMetadata, segmentsAmount)

	go passSegments(segments, segmentsDataChan, segmentsAmount)

	if threads > int(segmentsAmount) {
		threads = int(segmentsAmount)
	}
	for range threads {
		go worker(workerContext, videoName, segmentsDataChan, fileMetadataChan, rawDownloadurl)
	}

	bar := progressbar.Default(int64(segmentsAmount), "Downloading video...")
	var counter int
	paths := make([]string, segmentsAmount)
	for {
		fileMedatada := <-fileMetadataChan
		paths[fileMedatada.number] = fileMedatada.path
		bar.Add(1)
		counter++
		if counter == int(segmentsAmount) {
			break
		}
	}

	compileVideo(paths, videoName)
}

func worker(
	ctx context.Context,
	dirName string,
	segmentUrlsChan <-chan urlMetadata,
	fileMetadataChan chan<- fileMetadata,
	rawDownloadurl string,
) {

	// TODO: нужно будет сделать так, чтобы те сегменты, которые не скачались возвращались в стек и скачивались повторно (макс 3 раза)

	for {
		select {
		case <-ctx.Done():
			return
		case segmentData, ok := <-segmentUrlsChan:
			if !ok {
				return
			}

			splitUrl := strings.Split(segmentData.url, "/")
			filename := splitUrl[len(splitUrl)-1]

			filePath := path.Join(dirName, filename)

			file, err := os.Create(filePath)

			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			downloadUrl := fmt.Sprintf("%s/%s", rawDownloadurl, segmentData.url)

			res, err := http.Get(downloadUrl)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			data, err := io.ReadAll(res.Body)

			_, err = file.Write(data)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			file.Close()

			fileMetadataChan <- fileMetadata{path: filePath, number: segmentData.number}
		}
	}
}

func passSegments(
	segments []*m3u8.MediaSegment,
	segmentsDataChan chan<- urlMetadata,
	segmentsAmount uint,
) {
	for i := range segmentsAmount {
		segmentsDataChan <- urlMetadata{url: segments[i].URI, number: i}
	}
}

func compileVideo(filePaths []string, videoName string) {
	tsFilename := fmt.Sprintf("%s.ts", videoName)
	mp4Filename := fmt.Sprintf("%s.mp4", videoName)
	completeTsFile, err := os.Create(tsFilename) // TODO: сделать так, чтобы файл создавался с разными расширениями
	if err != nil {
		panic(err)
	}

	for _, filePath := range filePaths {
		file, err := os.Open(filePath)
		if err != nil {
			panic(err)
		}

		_, err = io.Copy(completeTsFile, file)
		if err != nil {
			panic(err)
		}
	}

	completeTsFile.Close()

	videoFile, err := os.Create(mp4Filename)
	if err != nil {
		panic(err)
	}
	defer videoFile.Close()

	ffmpegKwArgs := ffmpeg_go.KwArgs{
		"vcodec": "copy",
		"acodec": "copy",
	}
	ffmpeg_go.Input(tsFilename).Output(mp4Filename, ffmpegKwArgs).WithOutput(videoFile).Silent(true).OverWriteOutput().Run()

	err = os.RemoveAll(videoName)
	if err != nil {
		panic(err)
	}
	err = os.Remove(tsFilename)
	if err != nil {
		panic(err)
	}

	os.Exit(0)
}
