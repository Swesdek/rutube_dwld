package rutube

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/Swesdek/rutube-dwld/internal/interactions"
	"github.com/grafov/m3u8"
)

type videoInfo struct {
	Title        string `json:"title"`
	videoBalaner `json:"video_balancer"`
}

type videoBalaner struct {
	M3u8Url string `json:"m3u8"`
}

func GetVideoInfo(url string, private bool) (string, []*m3u8.MediaSegment, uint, string) {
	urlParts := strings.Split(url, "/")

	var privateInc int
	if private {
		privateInc += 1
	}

	videoID := urlParts[4+privateInc]
	var pk string

	match, err := regexp.MatchString(`\?p=(.*)`, urlParts[5+privateInc])
	if err != nil {
		panic(err)
	}

	if match {
		pk = fmt.Sprintf("/%s", urlParts[5+privateInc])
	}

	c := &http.Client{}

	apiURL := fmt.Sprintf("https://rutube.ru/api/play/options/%s%s", videoID, pk)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")

	res, err := c.Do(req)
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

	var mediaManifestURL string

	if len(masterPlaylist.Variants) == 1 {
		mediaManifestURL = resolutions[masterPlaylist.Variants[0].Resolution]
	} else {
		mediaManifestURL = interactions.SuggestResolution(resolutions)
	}

	res, err = http.Get(mediaManifestURL)
	if err != nil {
		panic(err)
	}

	mediaManifestData, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	buffer = bytes.NewBuffer(mediaManifestData)

	playlist, _, err = m3u8.Decode(*buffer, false)
	if err != nil {
		panic(err)
	}

	mediaPlaylist := playlist.(*m3u8.MediaPlaylist)
	splitMediaManURL := strings.Split(mediaManifestURL, "/")
	rawSegmentsURL := strings.Join(splitMediaManURL[:8], "/")
	return newVideoInfo.Title, mediaPlaylist.Segments, mediaPlaylist.Count(), rawSegmentsURL
}
