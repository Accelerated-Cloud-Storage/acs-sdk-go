// Package client provides a Go client for interacting with the Accelerated Cloud Storage service.
package client

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	pb "github.com/AcceleratedCloudStorage/acs-sdk-go/internal/generated"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/yaml.v2"
)

// CreateBucket sends a request to create a new bucket.
// It requires a bucket name and region specification.
// Returns an error if the bucket creation fails.
func (client *ACSClient) CreateBucket(ctx context.Context, bucket, region string) error {
	req := &pb.CreateBucketRequest{
		Bucket: bucket,
		Region: region,
	}

	_, err := client.client.CreateBucket(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create bucket: %v", err)
	}

	return nil
}

// DeleteBucket requests deletion of the specified bucket.
// Returns an error if the bucket deletion fails or the bucket doesn't exist.
func (client *ACSClient) DeleteBucket(ctx context.Context, bucket string) error {
	req := &pb.DeleteBucketRequest{
		Bucket: bucket,
	}

	_, err := client.client.DeleteBucket(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to delete bucket: %v", err)
	}

	return nil
}

// ListBuckets retrieves all buckets from the server.
// Returns a list of bucket objects and an error if the operation fails.
func (client *ACSClient) ListBuckets(ctx context.Context) ([]*pb.Bucket, error) {
	req := &pb.ListBucketsRequest{}

	resp, err := client.client.ListBuckets(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list buckets: %v", err)
	}

	return resp.Buckets, nil
}

// PutObject uploads data to the specified bucket and key.
// It automatically handles compression for large objects when beneficial.
// Returns an error if the upload fails.
func (client *ACSClient) PutObject(ctx context.Context, bucket, key string, data []byte) error {
	// Only compress if data is larger than threshold and compression would be beneficial
	const compressionThreshold = 100 * 1024 * 1024 // 1MB threshold
	isCompressed := false
	if len(data) >= compressionThreshold {
		var buf bytes.Buffer
		gw, err := gzip.NewWriterLevel(&buf, gzip.BestSpeed) // Use fastest compression
		if err != nil {
			return fmt.Errorf("failed to create gzip writer: %v", err)
		}
		if _, err := gw.Write(data); err != nil {
			return fmt.Errorf("failed to compress data: %v", err)
		}
		if err := gw.Close(); err != nil {
			return fmt.Errorf("failed to close gzip writer: %v", err)
		}

		// Only use compression if it actually reduces size
		if buf.Len() < len(data) {
			isCompressed = true
			data = buf.Bytes()
		}
	}

	stream, err := client.client.PutObject(ctx)
	if err != nil {
		return fmt.Errorf("failed to start PutObject stream: %v", err)
	}

	// Send parameters
	err = stream.Send(&pb.PutObjectRequest{
		Data: &pb.PutObjectRequest_Parameters{
			Parameters: &pb.PutObjectInput{
				Bucket:       bucket,
				Key:          key,
				IsCompressed: &isCompressed,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to send parameters: %v", err)
	}

	// Send data in chunks
	const chunkSize = 64 * 1024 // 64KB chunks
	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}

		err := stream.Send(&pb.PutObjectRequest{
			Data: &pb.PutObjectRequest_Chunk{
				Chunk: data[i:end],
			},
		})
		if err != nil {
			return fmt.Errorf("failed to send chunk: %v", err)
		}
	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		return fmt.Errorf("failed to close stream: %v", err)
	}

	return nil
}

// GetObject downloads the specified object from the server.
// Returns the object's data and an error if the download fails.
func (client *ACSClient) GetObject(ctx context.Context, bucket, key string) ([]byte, error) {
	req := &pb.GetObjectRequest{
		Bucket: bucket,
		Key:    key,
	}

	stream, err := client.client.GetObject(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to start GetObject stream: %v", err)
	}

	var data []byte
	firstMessage := true

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error receiving chunk: %v", err)
		}

		// Skip metadata message for now
		if firstMessage {
			firstMessage = false
			continue
		}

		// Append chunk data
		data = append(data, resp.GetChunk()...)
	}

	return data, nil
}

// DeleteObject removes a single object from a bucket.
// Returns an error if the deletion fails.
func (client *ACSClient) DeleteObject(ctx context.Context, bucket, key string) error {
	req := &pb.DeleteObjectRequest{
		Bucket: bucket,
		Key:    key,
	}

	_, err := client.client.DeleteObject(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to delete object: %v", err)
	}

	return nil
}

// HeadObject retrieves metadata for a specific object.
// Returns the object's metadata and an error if the operation fails.
func (client *ACSClient) HeadObject(ctx context.Context, bucket, key string) (*HeadObjectOutput, error) {
	req := &pb.HeadObjectRequest{
		Bucket: bucket,
		Key:    key,
	}

	resp, err := client.client.HeadObject(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to head object: %v", err)
	}

	return &HeadObjectOutput{
		ContentType:          resp.Metadata.ContentType,
		ContentLength:        resp.Metadata.Size,
		LastModified:         resp.Metadata.LastModified.AsTime(),
		ETag:                 resp.Metadata.Etag,
		ContentEncoding:      resp.Metadata.ContentEncoding,
		ContentLanguage:      resp.Metadata.ContentLanguage,
		VersionId:            resp.Metadata.VersionId,
		ServerSideEncryption: resp.Metadata.ServerSideEncryption,
		UserMetadata:         resp.Metadata.UserMetadata,
	}, nil
}

// DeleteObjects requests bulk deletion of objects in a bucket.
// Returns an error if any object deletion fails.
func (client *ACSClient) DeleteObjects(bucket string, keys []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Convert keys to object identifiers
	objects := make([]*pb.ObjectIdentifier, len(keys))
	for i, key := range keys {
		objects[i] = &pb.ObjectIdentifier{Key: key}
	}

	req := &pb.DeleteObjectsRequest{
		Bucket:  bucket,
		Objects: objects,
	}

	resp, err := client.client.DeleteObjects(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to delete objects: %v", err)
	}

	// Check for partial failures
	if len(resp.DeletedObjects) != len(keys) {
		return fmt.Errorf("some objects failed to delete")
	}

	return nil
}

// ListObjects retrieves object keys from the server based on given options.
// Returns a list of object keys and an error if the operation fails.
func (client *ACSClient) ListObjects(ctx context.Context, bucket string, opts *ListObjectsOptions) ([]string, error) {
	req := &pb.ListObjectsRequest{
		Bucket: bucket,
	}
	if opts != nil {
		if opts.Prefix != "" {
			req.Prefix = &opts.Prefix
		}
		if opts.StartAfter != "" {
			req.StartAfter = &opts.StartAfter
		}
		if opts.MaxKeys > 0 {
			req.MaxKeys = &opts.MaxKeys
		}
	}

	stream, err := client.client.ListObjects(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list objects: %v", err)
	}

	var keys []string
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if obj := resp.GetObject(); obj != nil {
			keys = append(keys, obj.Key)
		}
	}

	return keys, nil
}

// HeadBucket retrieves metadata for a specific bucket.
// Returns the bucket's metadata and an error if the operation fails.
func (client *ACSClient) HeadBucket(ctx context.Context, bucket string) (*HeadBucketOutput, error) {
	req := &pb.HeadBucketRequest{
		Bucket: bucket,
	}

	resp, err := client.client.HeadBucket(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to head bucket: %v", err)
	}

	return &HeadBucketOutput{
		Region: resp.BucketRegion,
	}, nil
}

// RotateKey checks if key rotation is needed and performs it if necessary.
// The force parameter can be used to require rotation regardless of timing.
// Returns an error if the rotation fails.
func (client *ACSClient) RotateKey(ctx context.Context, force bool) error {
	// Get the home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %v", err)
	}

	credsFile := filepath.Join(homeDir, ".acs", "credentials.yaml")

	// Read the entire credentials file
	data, err := os.ReadFile(credsFile)
	if err != nil {
		return fmt.Errorf("failed to read credentials file: %v", err)
	}

	var profiles profileCredentials
	if err := yaml.Unmarshal(data, &profiles); err != nil {
		return fmt.Errorf("failed to unmarshal credentials: %v", err)
	}

	// Get current profile name from environment variable
	profile := os.Getenv("ACS_PROFILE")
	if profile == "" {
		profile = "default"
	}

	creds, ok := profiles[profile]
	if !ok {
		return fmt.Errorf("profile '%s' not found in credentials file", profile)
	}

	resp, err := client.client.RotateKey(ctx, &pb.RotateKeyRequest{
		AccessKeyId: creds.AccessKeyID,
		Force:       &force,
	})
	if err != nil {
		return fmt.Errorf("key rotation failed: %v", err)
	}
	if !resp.Rotated {
		return nil
	}

	// Update only the current profile's credentials
	updatedCreds := creds
	updatedCreds.SecretAccessKey = resp.NewSecretAccessKey
	profiles[profile] = updatedCreds

	// Marshal updated profiles
	data, err = yaml.Marshal(profiles)
	if err != nil {
		return fmt.Errorf("failed to marshal credentials: %v", err)
	}

	// Write back to file
	if err := os.WriteFile(credsFile, data, 0600); err != nil {
		return fmt.Errorf("failed to update credentials file: %v", err)
	}

	return nil
}

// ShareBucket informs the service about a bucket that has been shared with it.
// Returns an error if the sharing operation fails or permissions are insufficient.
func (client *ACSClient) ShareBucket(ctx context.Context, bucket string) error {
	req := &pb.ShareBucketRequest{
		BucketName: bucket,
	}

	_, err := client.client.ShareBucket(ctx, req)
	if err != nil {
		// Check for specific error types and provide clear messages
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				return fmt.Errorf("bucket does not exist: %v", err)
			case codes.PermissionDenied:
				return fmt.Errorf("service lacks permission to access bucket: %v", err)
			case codes.InvalidArgument:
				return fmt.Errorf("invalid bucket name: %v", err)
			}
		}
		return fmt.Errorf("failed to share bucket: %v", err)
	}

	return nil
}

// CopyObject copies an object from a source bucket/key to a destination bucket/key.
// Returns an error if the copy operation fails.
func (client *ACSClient) CopyObject(ctx context.Context, bucket, copySource, key string) error {
	req := &pb.CopyObjectRequest{
		Bucket:     bucket,
		CopySource: copySource,
		Key:        key,
	}

	_, err := client.client.CopyObject(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to copy object: %v", err)
	}

	return nil
}
