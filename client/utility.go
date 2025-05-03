// Copyright 2025 Accelerated Cloud Storage Corporation. All Rights Reserved.
// Package client provides a Go client for interacting with the Accelerated Cloud Storage service.
package client

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pierrec/lz4/v4"
	"gopkg.in/yaml.v2"
)

// WithRange specifies a range for the GetObject operation
// The range should be in the format "bytes=start-end" (e.g., "bytes=0-9" for first 10 bytes)
func WithRange(rangeSpec string) GetObjectOption {
	return func(opts *GetObjectOptions) {
		opts.rangeSpec = rangeSpec
	}
}

// estimateCompressionRatio estimates the LZ4 compression ratio by sampling the data
func estimateCompressionRatio(data []byte) (float64, error) {
	totalSize := len(data)

	// Calculate sample size as 1% of total size, bounded between MIN and MAX
	targetSampleSize := int(float64(totalSize) * sampleRatio)
	if targetSampleSize < minSampleSize {
		targetSampleSize = minSampleSize
	}
	if targetSampleSize > maxSampleSize {
		targetSampleSize = maxSampleSize
	}

	// Take three samples: beginning, middle, and end
	perSampleSize := targetSampleSize / 3
	middle := totalSize / 2
	samples := [][]byte{
		data[:perSampleSize],
		data[middle-perSampleSize/2 : middle+perSampleSize/2],
		data[len(data)-perSampleSize:],
	}

	// Test compression ratio on samples
	var totalSampleSize int
	var totalCompressedSize int

	for _, sample := range samples {
		var buf bytes.Buffer
		w := lz4.NewWriter(&buf)
		w.Apply(lz4.CompressionLevelOption(0))
		if _, err := w.Write(sample); err != nil {
			return 0, fmt.Errorf("compression sample failed: %v", err)
		}
		if err := w.Close(); err != nil {
			return 0, fmt.Errorf("compression close failed: %v", err)
		}

		totalSampleSize += len(sample)
		totalCompressedSize += buf.Len()
	}

	return float64(totalCompressedSize) / float64(totalSampleSize), nil
}

// loadACSCredentials loads the service's credentials from the ~/.acs/credentials.yaml file.
// It creates the file with default values if it doesn't exist.
func loadACSCredentials() (*credentialsContents, error) {
	// Find home directory of user
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %v", err)
	}

	// Create the .acs directory if it doesn't exist
	acsDir := filepath.Join(homeDir, ".acs")
	if _, err := os.Stat(acsDir); os.IsNotExist(err) {
		fmt.Println("Accelerated Cloud Storage credential directory does not exist, creating it now . . .")
		if err := os.Mkdir(acsDir, 0700); err != nil {
			return nil, fmt.Errorf("failed to create .acs directory: %v", err)
		}
	}

	// Check if the credentials file exists
	credsFile := filepath.Join(acsDir, "credentials.yaml")
	if _, err := os.Stat(credsFile); os.IsNotExist(err) {
		fmt.Println("Accelerated Cloud Storage credentials file does not exist, creating it now . . .")
		defaultCreds := profileCredentials{
			"default": {
				AccessKeyID:     "your_access_key_id",
				SecretAccessKey: "your_secret_access_key",
			},
		}
		data, err := yaml.Marshal(defaultCreds)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal default credentials: %v", err)
		}
		if err := os.WriteFile(credsFile, data, 0600); err != nil {
			return nil, fmt.Errorf("failed to write default credentials file: %v", err)
		}
		creds := defaultCreds["default"]
		return &creds, nil
	}

	// Read the credentials file
	data, err := os.ReadFile(credsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read credentials file: %v", err)
	}

	var profiles profileCredentials
	if err := yaml.Unmarshal(data, &profiles); err != nil {
		return nil, fmt.Errorf("failed to unmarshal credentials: %v", err)
	}

	// Get profile from environment variable, default to "default" if not set
	profile := os.Getenv("ACS_PROFILE")
	if profile == "" {
		fmt.Println("ACS_PROFILE environment variable not set, using 'default' profile.")
		profile = "default"
	}

	creds, ok := profiles[profile]
	if !ok {
		return nil, fmt.Errorf("profile '%s' not found in credentials file", profile)
	}

	return &creds, nil
}
