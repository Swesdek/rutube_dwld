package rutube

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/grafov/m3u8"
)

type videoInfo struct {
	Title        string `json:"title"`
	videoBalaner `json:"video_balancer"`
}

type videoBalaner struct {
	M3u8Url string `json:"m3u8"`
}

func GetVideoInfo(url string) (string, []*m3u8.MediaSegment, uint, string) {
	urlParts := strings.Split(url, "/")

	videoId := urlParts[4]

	apiUrl := fmt.Sprintf("https://rutube.ru/api/play/options/%s", videoId)
	res, err := http.Get(apiUrl)
	if err != nil {
		panic(err)
	}

	decoder := json.NewDecoder(res.Body)
	var newVideoInfo videoInfo
	err = decoder.Decode(&newVideoInfo)
	if err != nil {
		panic(err)
	}

	res, err = http.Get(newVideoInfo.M3u8Url)
	if err != nil {
		panic(err)
	}

	masterManifestData, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	buffer := bytes.NewBuffer(masterManifestData)

	playlist, _, err := m3u8.Decode(*buffer, false)
	if err != nil {
		panic(err)
	}

	masterPlaylist := playlist.(*m3u8.MasterPlaylist)

	resolutions := make(map[string]string)
	for _, variant := range masterPlaylist.Variants {
		resolutions[variant.Resolution] = variant.URI
	}

	// Здесь нужно будет добавить интеракцию с пользователем, с вопросом о том,
	// какое нужно будет разрешение видео
	//
	// Пока будет лучшее разрешение

	mediaManifestUrl := resolutions[masterPlaylist.Variants[len(masterPlaylist.Variants)-1].Resolution]

	res, err = http.Get(mediaManifestUrl)
	if err != nil {
		panic(err)
	}

	mediaManifestData, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	buffer = bytes.NewBuffer(mediaManifestData)

	playlist, _, err = m3u8.Decode(*buffer, false)

	mediaPlaylist := playlist.(*m3u8.MediaPlaylist)

	splitMediaManUrl := strings.Split(mediaManifestUrl, "/")

	rawSegmentsUrl := strings.Join(splitMediaManUrl[:8], "/")

	return newVideoInfo.Title, mediaPlaylist.Segments, mediaPlaylist.Count(), rawSegmentsUrl
}
