// Copyright 2025 Accelerated Cloud Storage Corporation. All Rights Reserved.
// Package client provides a Go client for interacting with the Accelerated Cloud Storage service.
package client

import "time"

// serverAddress is the endpoint for the ACS service.
const (
	serverAddress = "acceleratedcloudstorages3cache.com:50050"
)

// HeadBucketOutput represents the metadata returned by HeadBucket operation.
// It contains information about a bucket's configuration and status.
type HeadBucketOutput struct {
	Region string // The region where the bucket is located
}

// HeadObjectOutput represents the metadata returned by HeadObject operation.
// It contains detailed information about an object's properties and metadata.
type HeadObjectOutput struct {
	ContentType          string            // MIME type of the object
	ContentEncoding      string            // Encoding of the object content
	ContentLanguage      string            // Language the object content is in
	ContentLength        int64             // Size of the object in bytes
	LastModified         time.Time         // Last modification timestamp
	ETag                 string            // Entity tag for the object
	UserMetadata         map[string]string // User-defined metadata key-value pairs
	ServerSideEncryption string            // Type of server-side encryption used
	VersionId            string            // Version identifier for the object
}

// ListObjectsOptions holds optional parameters for object listing.
// These options allow for customizing the object listing operation.
type ListObjectsOptions struct {
	Prefix     string // Filter objects by prefix
	StartAfter string // Return objects lexicographically after this value
	MaxKeys    int32  // Maximum number of keys to return
}

// credentialsContents holds the access key ID and secret access key.
// This structure matches the format of the credentials file.
type credentialsContents struct {
	AccessKeyID     string `yaml:"access_key_id"`     // AWS-style access key identifier
	SecretAccessKey string `yaml:"secret_access_key"` // Secret key for authentication
}

// profileCredentials holds multiple named credential profiles.
// Each profile contains a set of credentials for accessing the service.
type profileCredentials map[string]credentialsContents


// GetObjectOptions holds the options for GetObject
type GetObjectOptions struct {
	rangeSpec string
}

// GetObjectOption is a function that configures GetObjectOptions
type GetObjectOption func(*GetObjectOptions)
