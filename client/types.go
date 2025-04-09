// Copyright 2025 Accelerated Cloud Storage Corporation. All Rights Reserved.
// Package client provides a Go client for interacting with the Accelerated Cloud Storage service.
package client

import "time"

// serverAddress is the endpoint for the ACS service.
const (
	serverAddress = "acceleratedcloudstorageproduction.com:50050"
)

// Session represents a client session configuration.
// It contains optional parameters that can be provided when creating a new client.
type Session struct {
	// Region specifies the AWS region to use for this session
	Region string
}

// HeadBucketOutput represents the metadata returned by HeadBucket operation.
// It contains information about a bucket's configuration and status.
type HeadBucketOutput struct {
	// Region specifies the region where the bucket is located
	Region string
}

// HeadObjectOutput represents the metadata returned by HeadObject operation.
// It contains detailed information about an object's properties and metadata.
type HeadObjectOutput struct {
	// ContentType specifies the MIME type of the object
	ContentType string
	// ContentEncoding specifies the encoding of the object content
	ContentEncoding string
	// ContentLanguage specifies the language the object content is in
	ContentLanguage string
	// ContentLength specifies the size of the object in bytes
	ContentLength int64
	// LastModified is the last modification timestamp
	LastModified time.Time
	// ETag is the entity tag for the object
	ETag string
	// UserMetadata contains user-defined metadata key-value pairs
	UserMetadata map[string]string
	// ServerSideEncryption specifies the type of server-side encryption used
	ServerSideEncryption string
	// VersionId is the version identifier for the object
	VersionId string
}

// ListObjectsOptions holds optional parameters for object listing.
// These options allow for customizing the object listing operation.
type ListObjectsOptions struct {
	// Prefix filters objects by prefix
	Prefix string
	// StartAfter returns objects lexicographically after this value
	StartAfter string
	// MaxKeys specifies the maximum number of keys to return
	MaxKeys int32
}

// credentialsContents holds the access key ID and secret access key.
// This structure matches the format of the credentials file.
type credentialsContents struct {
	// AccessKeyID is the AWS-style access key identifier
	AccessKeyID string `yaml:"access_key_id"`
	// SecretAccessKey is the secret key for authentication
	SecretAccessKey string `yaml:"secret_access_key"`
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
