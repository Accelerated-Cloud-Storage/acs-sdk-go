package client

import "time"

const (
	serverAddress = "acceleratedcloudstorages3cache.com:50050"
)

// HeadBucketOutput represents the metadata returned by HeadBucket operation
type HeadBucketOutput struct {
	Region string
}

// HeadObjectOutput represents the metadata returned by HeadObject operation
type HeadObjectOutput struct {
	ContentType     string
	ContentEncoding string
	ContentLanguage string
	ContentLength   int64
	LastModified    time.Time
	ETag            string
	UserMetadata    map[string]string
	ServerSideEncryption string
	VersionId       string
}

// ListObjectsOptions holds optional parameters for object listing.
type ListObjectsOptions struct {
	Prefix     string
	StartAfter string
	MaxKeys    int32
}

// Credentials holds the access key ID and secret access key.
type credentialsContents struct {
	AccessKeyID     string `yaml:"access_key_id"`
	SecretAccessKey string `yaml:"secret_access_key"`
}

// ProfileCredentials holds multiple named credential profiles
type profileCredentials map[string]credentialsContents