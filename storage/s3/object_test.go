package s3_test

import (
	"context"
	"os"
	"path"
	"testing"

	"github.com/cwinters8/webster/storage/s3"
	"github.com/cwinters8/webster/utils"
)

func TestBucketOps(t *testing.T) {
	endpoint, region, accessKey, secretKey := "S3_ENDPOINT", "S3_REGION", "S3_ACCESS_KEY", "S3_SECRET_KEY"
	env, err := utils.GetEnv(endpoint, accessKey, secretKey)
	if err != nil {
		t.Fatalf("failed to get env vars: %v", err)
	}
	client, err := s3.NewClient(env[endpoint], env[region], env[accessKey], env[secretKey])
	if err != nil {
		t.Fatalf("failed to create S3 client: %v", err)
	}
	ctx := context.Background()

	// new bucket
	bucket := "webster-test"
	if err := client.NewBucket(ctx, bucket, false); err != nil {
		t.Fatalf("failed to create bucket: %v", err)
	}
	defer func(t *testing.T) {
		if err := client.RemoveBucket(ctx, bucket); err != nil {
			t.Fatal(err)
		}
	}(t)
	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Errorf("wanted bucket %s to exist", bucket)
	}

	t.Run("file", func(t *testing.T) {
		file := "hello.txt"
		content := "Hello World"
		if err := os.MkdirAll("tmp", 0755); err != nil {
			t.Fatalf("failed to create parent tmp dir: %v", err)
		}
		tmp, err := os.MkdirTemp("tmp", "")
		if err != nil {
			t.Fatalf("failed to create tmp dir: %v", err)
		}
		p := path.Join(tmp, file)
		if err := os.WriteFile(p, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write file: %v", err)
		}

		_, size, err := client.PutFile(ctx, bucket, file, p, nil)
		if err != nil {
			t.Fatalf("failed to put file in bucket: %v", err)
		}
		defer func(t *testing.T) {
			if err := client.RemoveFile(ctx, bucket, file); err != nil {
				t.Fatalf("failed to clean up file: %v", err)
			}
		}(t)
		if size == 0 {
			t.Errorf("wanted size of put file to be greater than 0")
		}

		// retrieve it and verify contents
		getTmp, err := os.MkdirTemp("tmp", "")
		if err != nil {
			t.Fatalf("failed to create tmp dir: %v", err)
		}
		if err := client.GetFile(ctx, bucket, file, getTmp, ""); err != nil {
			t.Fatal(err)
		}
		out, err := os.ReadFile(path.Join(getTmp, file))
		if err != nil {
			t.Fatalf("failed to read file: %v", err)
		}
		got := string(out)
		if got != content {
			t.Errorf("wanted file contents `%s`; got `%s`", content, got)
		}
	})
}
