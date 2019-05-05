package gfuns

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"
)

func Notifiy(data string, exchange string) {
	jsonValue, _ := json.Marshal(Result{Properties: Properties{ContentType: "application/json", DeliveryMode: 2, Priority: 0, Timestamp: int(time.Now().Unix()), MessageID: uuid.New().String()}, Payload: string(data), PayloadEncoding: "string", RoutingKey: exchange + ".s1"})
	_, err := http.Post("https://"+os.Getenv("username")+":"+os.Getenv("password")+"@"+os.Getenv("host")+"/api/exchanges/"+os.Getenv("username")+"/"+exchange+"/publish", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Fatal(err)
	}
}
func Next(fun string, gcsEvent GCSEvent) {
	show(fun, gcsEvent)
	jsonValue, _ := json.Marshal(gcsEvent)
	_, _ = http.Post("https://us-central1-ivory-program-229516.cloudfunctions.net/"+fun, "application/json", bytes.NewBuffer(jsonValue))
}
func show(fun string, gcsEvent GCSEvent) {
	jsonValue, _ := json.Marshal(gcsEvent)
	fmt.Println(fun, string(jsonValue))
}
func chmod(path string) error {
	return exec.Command("chmod", "+x", path).Run()

}
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
func GetResolutionImpls(probe FFprobe) []Resolution {
	temp := 0
	resolution := Resolution{}

	for i := 0; i < probe.Format.NbStreams; i++ {
		if probe.Streams[i].CodecType == "video" {
			if probe.Streams[i].Width*probe.Streams[i].Height > temp {
				temp = probe.Streams[i].Width * probe.Streams[i].Height
				resolution = Resolution{Width: (probe.Streams[i].Width * 2) / 2, Height: (probe.Streams[i].Height * 2) / 2}
			}
		}
	}
	return getLess(resolution)
}
func GetBeastResolution(probe FFprobe) Resolution {
	temp := 0
	resolution := Resolution{}

	for i := 0; i < probe.Format.NbStreams; i++ {
		if probe.Streams[i].CodecType == "video" {
			if probe.Streams[i].Width*probe.Streams[i].Height > temp {
				temp = probe.Streams[i].Width * probe.Streams[i].Height
				resolution = Resolution{Width: (probe.Streams[i].Width * 2) / 2, Height: (probe.Streams[i].Height * 2) / 2}
			}
		}
	}
	return resolution
}

func getLess(resolution Resolution) []Resolution {
	var temp []Resolution
	all := []Resolution{{Width: 256, Height: 144}, {Width: 426, Height: 240}, {Width: 640, Height: 360}, {Width: 854, Height: 480}, {Width: 1280, Height: 720}, {Width: 1920, Height: 1080}, {Width: 2560, Height: 1440}, {Width: 3840, Height: 2160}, {Width: 7680, Height: 4320}}
	for _, v := range all {
		if v.Width*v.Height <= resolution.Width*resolution.Height {
			temp = append(temp, Resolution{Id: bson.NewObjectId().Hex(), Width: v.Width, Height: v.Height})
		}

	}
	return temp
}
func createDir(rawUrl string) (string, error) {
	split := strings.Split(getPath(rawUrl), "/")
	sum := ""
	for _, v := range split[2 : len(split)-1] {
		sum += "/" + v
	}
	final := os.TempDir() + sum
	_ = os.RemoveAll(final)
	err := os.MkdirAll(final, 0700)
	return final, err
}
func getPath(rawUrl string) string {
	parse, _ := url.Parse(rawUrl)
	return parse.Path
}
func getBaseName(name string) string {
	split := strings.Split(name, ".")
	return split[0]
}
func GetName(rawUrl string) string {
	parse, _ := url.Parse(rawUrl)
	split := strings.Split(parse.Path, "/")
	return split[len(split)-1]
}
func Decode(v interface{}, request *http.Request) error {
	if request.Method != http.MethodPost {
		return errors.New("Method Not Allowed")
	}
	if request.Header.Get("Content-Type") != "application/json" {
		return errors.New("Unsupported Media Type")
	}
	return json.NewDecoder(request.Body).Decode(v)
}
