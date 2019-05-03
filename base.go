package func1

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func Run_init() {
	err := godotenv.Load("env")
	if err != nil {
		log.Fatal("Error loading env file", err)
	}
	load()
}
func load() {
	if !Exists(os.TempDir() + "/ffmpeg") {
		ffmpeg, _ := os.Create(os.TempDir() + "/ffmpeg")
		ffprobe, _ := os.Create(os.TempDir() + "/ffprobe")
		_ = Download(GCSEvent{File: File{Url: "https://storage.googleapis.com/asrevo/ffmpeg"}}, *ffmpeg)
		_ = Download(GCSEvent{File: File{Url: "https://storage.googleapis.com/asrevo/ffprobe"}}, *ffprobe)
		_ = chmod(ffmpeg.Name())
		_ = chmod(ffprobe.Name())
	}
}
