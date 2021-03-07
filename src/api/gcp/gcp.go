package gcp

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/url"
	"path/filepath"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

const (
	bucketName = "memeshop"
)

var imgCloudStore *ImageCloudStore

//ImageCloudStore struct
type ImageCloudStore struct {
	StorageClient *storage.Client
}

//NewStorageClient creates a new storage client instance
func NewStorageClient() {
	fp, err := filepath.Abs("../api/gcp/keys.json")
	if err != nil {
		log.Fatal(err)
	}

	storageClient, err := storage.NewClient(context.Background(), option.WithCredentialsFile(fp))
	if err != nil {
		log.Fatal(err)
	}
	imgCloudStore = &ImageCloudStore{storageClient}
}

//Connect disconnects from GCP service
func Connect() *ImageCloudStore {
	return imgCloudStore
}

//Disconnect disconnects from GCP service
func Disconnect() {
	if err := imgCloudStore.StorageClient.Close(); err != nil {
		log.Fatal(err)
	}
}

//DeleteFile deletes an uploaded file
func (i *ImageCloudStore) DeleteFile(filename string) error {
	return i.StorageClient.Bucket(bucketName).Object(filename).Delete(context.Background())
}

//UploadFile uploads a file
func (i *ImageCloudStore) UploadFile(f *[]byte, filename string) (string, error) {
	sw := i.StorageClient.Bucket(bucketName).Object(filename).NewWriter(context.Background())
	if _, err := io.Copy(sw, bytes.NewReader(*f)); err != nil {
		return "", err
	}

	if err := sw.Close(); err != nil {
		return "", err
	}

	u, err := url.Parse("/" + bucketName + "/" + sw.Attrs().Name)
	if err != nil {
		return "", err
	}

	return "https://storage.googleapis.com" + u.String(), nil
}
