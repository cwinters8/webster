package s3_test

import (
	"context"
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/cwinters8/webster/storage/s3"
	"github.com/cwinters8/webster/utils"
)

func TestBucketOps(t *testing.T) {
	region, accessKey, secretKey := "S3_REGION", "S3_ACCESS_KEY", "S3_SECRET_KEY"
	env, err := utils.GetEnv(accessKey, secretKey)
	if err != nil {
		t.Fatalf("failed to get env vars: %v", err)
	}
	client, err := s3.NewClient("s3.amazonaws.com", env[region], env[accessKey], env[secretKey])
	if err != nil {
		t.Fatalf("failed to create S3 client: %v", err)
	}
	ctx := context.Background()

	newFile := func(t *testing.T, bucket, name string) (versionID string, err error) {
		content := "Hello World"
		if err := os.MkdirAll("tmp", 0755); err != nil {
			return "", fmt.Errorf("failed to create parent tmp dir: %v", err)
		}
		tmp, err := os.MkdirTemp("tmp", "")
		if err != nil {
			return "", fmt.Errorf("failed to create tmp dir: %v", err)
		}
		p := path.Join(tmp, name)
		if err := os.WriteFile(p, []byte(content), 0644); err != nil {
			return "", fmt.Errorf("failed to write file: %v", err)
		}

		version, size, err := client.PutFile(ctx, bucket, name, p, nil)
		if err != nil {
			return "", fmt.Errorf("failed to put file in bucket: %v", err)
		}
		if size == 0 {
			t.Errorf("wanted size of put file to be greater than 0")
		}

		// retrieve it and verify contents
		getTmp, err := os.MkdirTemp("tmp", "")
		if err != nil {
			return "", fmt.Errorf("failed to create tmp dir: %v", err)
		}
		if err := client.GetFile(ctx, bucket, name, getTmp, ""); err != nil {
			return "", err
		}
		out, err := os.ReadFile(path.Join(getTmp, name))
		if err != nil {
			return "", fmt.Errorf("failed to read file: %v", err)
		}
		got := string(out)
		if got != content {
			t.Errorf("wanted file contents `%s`; got `%s`", content, got)
		}
		return version, nil
	}

	t.Run("versioning", func(t *testing.T) {
		t.Parallel()
		// new bucket
		bucket := "webster-test-versioned"
		if err := client.NewBucket(ctx, bucket, true); err != nil {
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

		file := "hello.txt"
		version, err := newFile(t, bucket, file)
		defer func(t *testing.T) {
			if err := client.RemoveFile(ctx, bucket, file); err != nil {
				t.Fatalf("failed to clean up file: %v", err)
			}
		}(t)
		if err != nil {
			t.Fatal(err)
		}
		t.Log("version ID:", version)
		if len(version) == 0 {
			t.Error("wanted non-empty version ID")
		}
	})

	t.Run("no_versioning", func(t *testing.T) {
		t.Parallel()
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

		file := "hello.txt"
		version, err := newFile(t, bucket, file)
		defer func(t *testing.T) {
			if err := client.RemoveFile(ctx, bucket, file); err != nil {
				t.Fatalf("failed to clean up file: %v", err)
			}
		}(t)
		if err != nil {
			t.Fatal(err)
		}
		if len(version) != 0 {
			t.Errorf("wanted empty version ID; got %s", version)
		}
	})

	// TODO: verify file removal only removes the files that should be removed
	t.Run("remove", func(t *testing.T) {
		t.Parallel()
		bucket := "webster-delete"
		if err := client.NewBucket(ctx, bucket, false); err != nil {
			t.Fatalf("failed to create bucket: %v", err)
		}
		defer func(t *testing.T) {
			if err := client.RemoveBucket(ctx, bucket); err != nil {
				t.Fatal(err)
			}
		}(t)
		hello, world := "hello.txt", "world.txt"
		_, helloErr := newFile(t, bucket, hello)
		_, worldErr := newFile(t, bucket, world)
		defer func(t *testing.T) {
			if err := client.RemoveFile(ctx, bucket, hello); err != nil {
				t.Fatalf("failed to clean up file %s: %v", hello, err)
			}
			if err := client.RemoveFile(ctx, bucket, world); err != nil {
				t.Fatalf("failed to clean up file %s: %v", world, err)
			}
		}(t)
		if helloErr != nil || worldErr != nil {
			t.Fatalf("one or more files failed to create. errors:\n%v\n%v", helloErr, worldErr)
		}

		// remove hello and make sure world stays
		if err := client.RemoveFile(ctx, bucket, hello); err != nil {
			t.Fatalf("failed to remove file %s: %v", hello, err)
		}
		found, err := client.FileExists(ctx, bucket, world)
		if err != nil {
			t.Fatalf("failed to check file existence: %v", err)
		}
		if !found {
			t.Errorf("wanted %s to exist in the bucket", world)
		}
	})

	t.Run("nonexistent_file", func(t *testing.T) {
		t.Parallel()
		bucket := "webster-fake-file"
		if err := client.NewBucket(ctx, bucket, false); err != nil {
			t.Fatalf("failed to create bucket: %v", err)
		}
		defer func(t *testing.T) {
			if err := client.RemoveBucket(ctx, bucket); err != nil {
				t.Fatal(err)
			}
		}(t)
		// ensure client.FileExists returns false for a file that doesn't exist
		file := "faker.txt"
		found, err := client.FileExists(ctx, bucket, file)
		if err != nil {
			t.Fatalf("failed to check file existence: %v", err)
		}
		if found {
			t.Errorf("wanted false; %s should not exist", file)
		}
	})
}
