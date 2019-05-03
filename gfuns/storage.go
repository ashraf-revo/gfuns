package gfuns

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

func generateUrl(bucket, key string) (string, error) {
	return storage.SignedURL(bucket, key, &storage.SignedURLOptions{
		GoogleAccessID: os.Getenv("GoogleAccessID"),
		PrivateKey:     []byte(os.Getenv("PrivateKey")),
		Expires:        time.Now().Add(time.Second * 60),
		Method:         "GET",
	})
}
func push(client *storage.Client, ctx context.Context, bucket, file, object string) error {
	fmt.Println(bucket, file, object)
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer func() { _ = f.Close() }()

	wc := client.Bucket(bucket).Object(object).NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
