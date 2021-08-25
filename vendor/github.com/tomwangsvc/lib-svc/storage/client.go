package storage

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"cloud.google.com/go/storage"
	lib_errors "github.com/tomwangsvc/lib-svc/errors"
	lib_log "github.com/tomwangsvc/lib-svc/log"
)

type Client interface {
	Bucket(name string) *storage.BucketHandle
	Close() error
	ReadFromBucket(ctx context.Context, bucket, key string) (object []byte, contentType string, err error)
}

func NewClient(ctx context.Context) (Client, error) {
	lib_log.Info(ctx, "Initializing")

	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		return nil, lib_errors.Wrap(err, "Failed initializing storage client")
	}

	lib_log.Info(ctx, "Initialized")
	return client{
		storageClient: storageClient,
	}, nil
}

type client struct {
	storageClient *storage.Client
}

func (c client) Bucket(name string) *storage.BucketHandle {
	return c.storageClient.Bucket(name)
}

func (c client) Close() error {
	if err := c.storageClient.Close(); err != nil {
		return lib_errors.Wrap(err, "Failed closing storage client")
	}
	return nil
}

func (c client) ReadFromBucket(ctx context.Context, bucket, key string) (object []byte, contentType string, err error) {
	lib_log.Info(ctx, "Reading", lib_log.FmtString("bucket", bucket), lib_log.FmtString("key", key))

	object, contentType, err = c.readFromBucket(ctx, bucket, key)
	if err != nil {
		err = lib_errors.Wrapf(err, "Failed reading from object from bucket=%s, key=%s", bucket, key)
	}

	lib_log.Info(ctx, "Read", lib_log.FmtString("bucket", bucket), lib_log.FmtString("key", key), lib_log.FmtInt("len(object)", len(object)), lib_log.FmtString("contentType", contentType))
	return
}

func (c client) readFromBucket(ctx context.Context, bucket, key string) (object []byte, contentType string, err error) {
	b := c.storageClient.Bucket(bucket)
	bucketObject := b.Object(key)

	reader, err := bucketObject.NewReader(ctx)
	if err != nil {
		if err.Error() == storage.ErrObjectNotExist.Error() {
			err = lib_errors.NewCustomWithCause(http.StatusNotFound, fmt.Sprintf("Object does not exist, bucket=%s, key=%s", bucket, key), err)
			return
		}
		err = lib_errors.Wrapf(err, "Failed creating new reader for bucket=%s, key=%s", bucket, key)
		return
	}
	defer reader.Close()
	contentType = reader.ContentType()

	buffer := new(bytes.Buffer)
	if _, err = buffer.ReadFrom(reader); err != nil {
		err = lib_errors.Wrapf(err, "Failed reading from reader for bucket=%s, key=%s", bucket, key)
		return
	}
	object = buffer.Bytes()

	return
}
