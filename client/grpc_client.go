package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"

	pb "github.com/AcceleratedCloudStorage/acs-sdk-go/internal/generated"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

// GRPCClient wraps the gRPC connection and client for ObjectStorageCache.
type ACSClient struct {
	client pb.ObjectStorageCacheClient
	conn   *grpc.ClientConn
}

// loadCredentials loads the credentials from the ~/.acs/credentials.yaml file.
func loadCredentials() (*credentialsContents, error) {
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

// NewClient initializes a new gRPC client with authentication
func NewClient() (*ACSClient, error) {
	// Configure TLS with the system cert pool
	rootCAs, err := x509.SystemCertPool()
    if err != nil {
        rootCAs = x509.NewCertPool()
    }
    config := &tls.Config{
        ServerName: "acceleratedcloudstorages3cache.com",
        RootCAs:    rootCAs,
        MinVersion: tls.VersionTLS12,
		InsecureSkipVerify: false, // Keep this false for production
		NextProtos: []string{"h2", "http/1.1"}, // Add ALPN protocols
    }
    
    creds := credentials.NewTLS(config)
    
    opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
        grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(1024*1024*1024), // 1GB
			grpc.MaxCallSendMsgSize(1024*1024*1024),  // 1GB
    	),
        grpc.WithKeepaliveParams(keepalive.ClientParameters{
            Time:                20 * time.Second,
            Timeout:             60 * time.Second,
            PermitWithoutStream: true,
        }),
    }
	opts = append(opts, 
    grpc.WithDefaultServiceConfig(`{
        "loadBalancingPolicy": "round_robin",
        "methodConfig": [{
            "name": [{}],
            "retryPolicy": {
                "maxAttempts": 3,
                "initialBackoff": "0.1s",
                "maxBackoff": "1s",
                "backoffMultiplier": 2.0
            }
        }]
    }`),
)

    // Create connection
    conn, err := grpc.NewClient(serverAddress, opts...)
    if err != nil {
        return nil, fmt.Errorf("failed to connect: %v", err)
    }

    // Create client
    client := &ACSClient{
        client: pb.NewObjectStorageCacheClient(conn),
        conn:   conn,
    }

    // Load credentials
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
func (client *ACSClient) Close() error {
	if client.conn != nil {
		return client.conn.Close()
	}
	return nil
}
