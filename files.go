package func1

import (
	"archive/zip"
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func Download(gcsEvent GCSEvent, path os.File) error {
	out, err := os.Create(path.Name())
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer func() { _ = out.Close() }()

	resp, err := http.Get(gcsEvent.File.Url)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
func Unzip(path string) ([]string, error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		return walk(path), nil
	}
	defer func() {
		if r != nil {
			_ = r.Close()
		} else {
			log.Fatal(err)
		}
	}()

	dir, _ := ioutil.TempDir(os.TempDir(), "temp")
	fmt.Println(dir)
	for _, file := range r.Reader.File {
		zippedFile, err := file.Open()
		if err != nil {
			log.Fatal(err)
		}
		defer func() { _ = zippedFile.Close() }()
		extractedFilePath := filepath.Join(
			dir,
			file.Name,
		)

		if file.FileInfo().IsDir() {
			_ = os.MkdirAll(extractedFilePath, file.Mode())
		} else {
			outputFile, err := os.OpenFile(
				extractedFilePath,
				os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
				file.Mode(),
			)
			if err != nil {
				log.Fatal(err)
			}
			defer func() { _ = outputFile.Close() }()

			_, err = io.Copy(outputFile, zippedFile)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	return walk(dir), nil
}
func walk(dir string) []string {
	var files []string
	_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			log.Fatal(err)
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files
}
func Each(files []string, getKey func(it string, to GCSEvent) string, ctx context.Context, to GCSEvent) ([]string, error) {
	client, err := storage.NewClient(ctx)
	var result []string
	if err != nil {
		log.Fatal(err)
		return result, err
	}
	var tempErr error = nil
	for _, file := range files {
		key := getKey(file, to)
		err := push(client, ctx, to.Bucket, file, key)
		if err != nil && tempErr == nil {
			log.Fatal(err)
			tempErr = err
		}
		result = append(result, key)
		_ = os.Remove(file)
	}
	return result, tempErr
}
