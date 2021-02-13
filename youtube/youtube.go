package youtube

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os/exec"
	"regexp"

	"github.com/Jeffail/gabs"
)

var (
	youtubeURL = regexp.MustCompile(`(?:https?:\/\/)?(?:www\.)?(?:youtu\.be\/|youtube\.com\/(?:embed\/|v\/|watch\?v=|watch\?.+&v=))((\w|-){11})(?:\S+)?`)
)

func Download(URL string) (string, error) {
	ytdl, err := exec.Command("youtube-dl", "-f", "bestaudio", "--get-url", URL).Output()
	if err != nil {
		return "", err
	}
	return string(ytdl), nil
}

func Search(searchArg string) (string, string, error) {
	if youtubeURL.MatchString(searchArg) {
		YtVideo, err := Download(youtubeURL.FindAllStringSubmatch(searchArg, 2)[0][1])
		if err != nil {
			return "", "", err
		}
		return "<" + searchArg + ">", YtVideo, nil
	}
	res, err := http.Get(fmt.Sprintf("https://www.googleapis.com/youtube/v3/search?part=snippet,id&maxResults=1&type=video&q=%s&key=AIzaSyBYczXOYKc6kZWKhn3V9m-vKR-s7CS-Ync", url.QueryEscape(searchArg)))
	if err != nil {
		return "", "", err
	}
	defer res.Body.Close()
	bytes, _ := ioutil.ReadAll(res.Body) // probably risky but I doubt this will ever error
	JSON, _ := gabs.ParseJSON(bytes)     // probably risky but I doubt this will ever error
	children, err := JSON.S("items").Children()
	if err != nil || len(children) == 0 {
		return "", "", fmt.Errorf("No Video")
	}
	videoID, _ := children[0].Path("id.videoId").Data().(string)
	videoLink, err := Download(videoID)
	if err != nil {
		return "", "", err
	}
	return "<https://www.youtube.com/watch?v=" + videoID + ">", videoLink, nil
}
