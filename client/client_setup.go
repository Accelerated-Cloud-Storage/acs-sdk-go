// Copyright 2025 Accelerated Cloud Storage Corporation. All Rights Reserved.
// Package client provides a Go client for interacting with the Accelerated Cloud Storage service.
package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"

	pb "github.com/AcceleratedCloudStorage/acs-sdk-go/generated"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

// ACSClient wraps the gRPC connection and client for ObjectStorageCache.
// It provides high-level operations for interacting with the ACS service.
type ACSClient struct {
	client pb.ObjectStorageCacheClient
	conn   *grpc.ClientConn
	retry  RetryConfig
}

var _ = embed.FS{}
//go:embed internal/ca-chain.pem
var embeddedCACert []byte
func loadClientTLSCredentials() (credentials.TransportCredentials, error) {
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(embeddedCACert) {
		log.Fatalf("Failed to append CA certificates")
	}

	tlsConfig := &tls.Config{
		RootCAs: certPool,
	}

	return credentials.NewTLS(tlsConfig), nil
}

// NewClient initializes a new gRPC client with authentication.
// It establishes a secure connection to the ACS service, loads credentials,
// and performs initial authentication. It also checks for key rotation needs.
// Returns an error if connection, authentication, or credential loading fails.
func NewClient() (*ACSClient, error) {
	tlsCredentials, err := loadClientTLSCredentials()
	if err != nil {
		log.Fatalf("Failed to load TLS credentials: %v", err)
	}
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(tlsCredentials),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(1024*1024*1024), // 1GB
			grpc.MaxCallSendMsgSize(1024*1024*1024), // 1GB
		),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second, // 10 seconds between pings
			Timeout:             5 * time.Second,  // 5 seconds timeout for pings
			PermitWithoutStream: true,
		}),
	}

	// Create connection
	conn, err := grpc.NewClient(serverAddress, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %v", err)
	}

	// Create client with default retry config
	client := &ACSClient{
		client: pb.NewObjectStorageCacheClient(conn),
		conn:   conn,
		retry:  DefaultRetryConfig,
	}

	// Load credentials from disk
	serviceCreds, err := loadCredentials()
	if err != nil {
		return nil, fmt.Errorf("failed to load credentials: %v", err)
	}

	// Perform authentication
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = client.client.Authenticate(ctx, &pb.AuthRequest{
		AccessKeyId:     serviceCreds.AccessKeyID,
		SecretAccessKey: serviceCreds.SecretAccessKey,
	})
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("authentication failed: %v", err)
	}
	// After successful authentication, check if key rotation is needed
	if err := client.RotateKey(ctx, false); err != nil {
		// Log the error but don't fail the connection
		fmt.Printf("Warning: Key rotation check failed: %v\n", err)
	}

	return client, nil
}

// Close terminates the client connection.
// It should be called when the client is no longer needed to free resources.
func (client *ACSClient) Close() error {
	if client.conn != nil {
		return client.conn.Close()
	}
	return nil
}

// loadCredentials loads the credentials from the ~/.acs/credentials.yaml file.
// It creates the credentials file with default values if it doesn't exist.
// The function respects the ACS_PROFILE environment variable to select the appropriate profile.
func loadCredentials() (*credentialsContents, error) {
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
		fmt.Println("No ACS_PROFILE environment variable set, using 'default' profile.")
		profile = "default"
	}

	creds, ok := profiles[profile]
	if !ok {
		return nil, fmt.Errorf("profile '%s' not found in credentials file", profile)
	}

	return &creds, nil
}
