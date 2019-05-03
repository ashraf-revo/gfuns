package gfuns

import (
	"context"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func init() {
	run_init()
}

//Download
func Func1(response http.ResponseWriter, request *http.Request) {
	load()
	var gcsEvent GCSEvent
	err := decode(&gcsEvent, request)
	if err != nil {
		http.Error(response, err.Error(), 415)
	}
	if gcsEvent.File.Url == "" {
		http.Error(response, "Url Can't Be Empty", 400)
		return

	}
	if gcsEvent.File.Id == "" {
		http.Error(response, "Id Can't Be Empty", 400)
		return

	}
	path, _ := ioutil.TempFile(os.TempDir(), "temp*"+filepath.Ext(getName(gcsEvent.File.Url)))
	err = download(gcsEvent, *path)
	defer func() { _ = os.Remove(path.Name()) }()
	if err != nil {
		log.Fatal(err)
		http.Error(response, "Error Download this file "+gcsEvent.File.Url, 400)
		return

	}
	files, err := unzip(path.Name())
	if err != nil {
		log.Fatal(err)
		http.Error(response, "Error unzip this file "+gcsEvent.File.Url, 400)
		return
	}

	result, err := each(files, func(it string, to GCSEvent) string {
		hex := bson.NewObjectId().Hex()
		return to.Pattern + "/" + hex + "/" + hex + "/" + hex + filepath.Ext(getName(it))
	}, context.Background(), GCSEvent{Bucket: "asrevo-video", Pattern: gcsEvent.File.Id})

	if err != nil {
		log.Fatal(err)
	}
	for _, name := range result {
		go next("func2", GCSEvent{Bucket: "asrevo-video", Name: name})
	}
	_, _ = fmt.Fprint(response)
}
