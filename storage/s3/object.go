package s3

// for interacting with S3-compatible object storage

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/cwinters8/webster/headers"
	"github.com/cwinters8/webster/utils"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Client struct {
	min *minio.Client
}

func NewClient(endpoint, region, accessKey, secretAccessKey string) (*Client, error) {
	c, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretAccessKey, ""),
		Secure: true,
		Region: region,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}
	return &Client{c}, nil
}

// NewBucket creates a new bucket and optionally enables versioning
//
// If the bucket already exists, a message is printed and an error is not returned.
// If versioning is true, it will still be enabled.
func (c *Client) NewBucket(ctx context.Context, name string, versioning bool) error {
	exists, err := c.min.BucketExists(ctx, name)
	if err != nil {
		return fmt.Errorf("failed to check if bucket exists: %w", err)
	}
	if !exists {
		if err := c.min.MakeBucket(ctx, name, minio.MakeBucketOptions{}); err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
	} else {
		fmt.Printf("bucket %s already exists\n", name)
	}
	if versioning {
		if err := c.min.EnableVersioning(ctx, name); err != nil {
			return fmt.Errorf("failed to enable bucket versioning: %w", err)
		}
	}
	return nil
}

func (c *Client) BucketExists(ctx context.Context, name string) (bool, error) {
	exists, err := c.min.BucketExists(ctx, name)
	if err != nil {
		return false, fmt.Errorf("failed to check if bucket exists: %w", err)
	}
	return exists, nil
}

func (c *Client) RemoveBucket(ctx context.Context, name string) error {
	if err := c.min.RemoveBucket(ctx, name); err != nil {
		return fmt.Errorf("failed to remove bucket: %w", err)
	}
	return nil
}

// PutFile writes a file at path to bucket/name
func (c *Client) PutFile(ctx context.Context, bucket, name, path string, options *PutOptions) (versionID string, size int64, err error) {
	opts := minio.PutObjectOptions{}
	if options != nil {
		opts = *options.toMin()
	}
	info, err := c.min.FPutObject(ctx, bucket, name, path, opts)
	if err != nil {
		return "", 0, fmt.Errorf("failed to put file: %w", err)
	}
	return info.VersionID, info.Size, nil
}

// GetFile attempts to retrieve the file specified by name from bucket and writes the file to target
//
// if target is a directory, the base filename from name will be appended to the target path
//
// An empty string can be passed for versionID to retrieve the latest version
func (c *Client) GetFile(ctx context.Context, bucket, name, target, versionID string) error {
	info, err := os.Stat(target)
	if err != nil {
		return fmt.Errorf("failed to stat target: %w", err)
	}
	if info.IsDir() {
		fileName := path.Base(name)
		target = path.Join(target, fileName)
	}
	if err := c.min.FGetObject(ctx, bucket, name, target, minio.GetObjectOptions{
		VersionID: versionID,
	}); err != nil {
		return fmt.Errorf("failed to get file: %w", err)
	}
	return nil
}

var ErrDoesNotExist = utils.Error{Msg: "does not exist"}

func (c *Client) FileExists(ctx context.Context, bucket, name string) (bool, error) {
	if _, err := c.min.StatObject(ctx, bucket, name, minio.GetObjectOptions{}); err != nil {
		if errors.Is(ErrDoesNotExist, err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to get file info: %w", err)
	}
	return true, nil
}

func (c *Client) RemoveFile(ctx context.Context, bucket, name string) error {
	for obj := range c.min.ListObjects(ctx, bucket, minio.ListObjectsOptions{
		Prefix:       name,
		WithVersions: true,
	}) {
		if err := c.min.RemoveObject(ctx, bucket, name, minio.RemoveObjectOptions{
			VersionID: obj.VersionID,
		}); err != nil {
			return fmt.Errorf("failed to remove object: %w", err)
		}
	}
	return nil
}

type PutOptions struct {
	ContentType        headers.ContentType
	ContentEncoding    headers.ContentEncoding
	ContentDisposition headers.ContentDisposition
}

func (opts *PutOptions) toMin() *minio.PutObjectOptions {
	return &minio.PutObjectOptions{
		ContentType:        string(opts.ContentType),
		ContentEncoding:    string(opts.ContentEncoding),
		ContentDisposition: string(opts.ContentDisposition),
	}
}
