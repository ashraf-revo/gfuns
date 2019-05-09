package gfuns

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"google.golang.org/api/iterator"
	"io"
	"log"
	"os"
	"time"
)

func GenerateUrl(bucket, key string) (string, error) {
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
func List(parent string) []string {
	var result []string
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	iter := client.Bucket("asrevo-video").Objects(ctx, &storage.Query{Prefix: parent})
	for {
		o, err := iter.Next()
		if err == iterator.Done {
			break
		}
		result = append(result, o.Name)
	}
	return result
}
