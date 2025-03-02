// Copyright 2025 Accelerated Cloud Storage Corporation. All Rights Reserved.
// Package client provides a Go client for interacting with the Accelerated Cloud Storage service.
package client

// WithRange specifies a range for the GetObject operation
// The range should be in the format "bytes=start-end" (e.g., "bytes=0-9" for first 10 bytes)
func WithRange(rangeSpec string) GetObjectOption {
	return func(opts *GetObjectOptions) {
		opts.rangeSpec = rangeSpec
	}
}
