package s3_test

import (
	"context"
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
	bucket := "webster-test"
	if err := client.NewBucket(ctx, bucket, false); err != nil {
		t.Fatalf("failed to create bucket: %v", err)
	}
	defer func(t *testing.T) {
		if err := client.RemoveBucket(ctx, bucket); err != nil {
			t.Fatal(err)
		}
		t.Logf("removed bucket %s", bucket)
	}(t)
	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Errorf("wanted bucket %s to exist", bucket)
	}
}
