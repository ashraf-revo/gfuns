package gfuns

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

func Split(url string, ffprobe FFprobe) ([]string, string, error) {
	dir, _ := createDir(ffprobe.Format.Filename)
	message, err := ffmpeg("-y", "-v", "error", "-i", url, "-f", "segment", "-codec:", "copy", "-segment_time", "60", dir+"/"+getBaseName(GetName(ffprobe.Format.Filename))+"_%d"+filepath.Ext(GetName(url)))
	return walk(dir), message, err
}
func Probe(url string) (FFprobe, error) {
	cmd := exec.Command(os.TempDir()+"/ffprobe", "-v", "quiet", "-print_format", "json", "-show_error", "-show_format", "-show_streams", url)
	fmt.Println(cmd)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
		return FFprobe{}, err
	}
	var app FFprobe
	return app, json.Unmarshal(out.Bytes(), &app)
}
func ffmpeg(arg ...string) (string, error) {
	cmd := exec.Command(os.TempDir()+"/ffmpeg", arg...)
	fmt.Println(cmd)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println(out.String(), err)
		log.Fatal(err)
		return out.String(), err
	}
	fmt.Println(out.String(), err)
	return out.String(), nil
}
func Mp4(url string, gcsEvent GCSEvent) ([]string, string, error) {
	dir, _ := createDir(gcsEvent.Pattern + "/" + GetName(gcsEvent.Name))
	message, err := ffmpeg("-y", "-v", "error", "-i", url, "-f", "mp4", "-s", strconv.FormatInt(int64(gcsEvent.Resolution.Width), 10)+"x"+strconv.FormatInt(int64(gcsEvent.Resolution.Height), 10), dir+"/"+getBaseName(GetName(gcsEvent.Name))+".mp4")
	return walk(dir), message, err
}
func Png(url string, gcsEvent GCSEvent) ([]string, string, error) {
	dir, _ := createDir(gcsEvent.Pattern + "/" + GetName(gcsEvent.Name))
	message, err := ffmpeg("-y", "-v", "error", "-i", url, "-vframes", "1", "-vf", "select='gte(n\\,10)',scale=320:-1", "-ss", "00:00:10", dir+"/"+getBaseName(GetName(gcsEvent.Name))+".png")
	return walk(dir), message, err
}
func Jpeg(url string, gcsEvent GCSEvent) ([]string, string, error) {
	dir, _ := createDir(gcsEvent.Pattern + "/" + GetName(gcsEvent.Name))
	message, err := ffmpeg("-y", "-v", "error", "-i", url, "-vf", "select='gte(n\\,10)',scale=144:-1", dir+"/"+getBaseName(GetName(gcsEvent.Name))+"_%d.jpeg")
	return walk(dir), message, err
}
func Webp(url string, gcsEvent GCSEvent) ([]string, string, error) {
	dir, _ := createDir(gcsEvent.Pattern + "/" + GetName(gcsEvent.Name))
	message, err := ffmpeg("-y", "-v", "error", "-i", url, "-loop", "0", "-vf", "select='gte(n\\,10)',scale=320:-1", "-ss", "00:00:2", "-t", "00:00:03", dir+"/"+getBaseName(GetName(gcsEvent.Name))+".webp")
	return walk(dir), message, err
}
func Hls(url string, gcsEvent GCSEvent) ([]string, string, error) {
	dir, _ := createDir(gcsEvent.Pattern + "/" + GetName(gcsEvent.Name))
	urls := generateMultipleUrl("asrevo-video", List(url))
	write(dir+"/files.txt", urls)
	name := getBaseName(GetName(gcsEvent.Name))
	message, err := ffmpeg("-y", "-v", "error", "-safe", "0", "-protocol_whitelist", "\"file,http,https,tcp,tls,concat\"", "-f", "concat", "-i", dir+"/files.txt", "-hls_segment_type", "mpegts", "-f", "hls", "-codec:", "copy", "-start_number", "0", "-hls_time", "2", "-hls_list_size", "0", "-hls_enc", "1", "-hls_enc_key", gcsEvent.Meta.Key, "-hls_enc_key_url", dir+"/"+name+".key", "-hls_enc_iv", gcsEvent.Meta.Iv, "-master_pl_name", name+".m3u8", dir+"/"+name+"_"+".m3u8")
	return walk(dir), message, err
}
