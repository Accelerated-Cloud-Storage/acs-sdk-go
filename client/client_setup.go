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
	"time"

	pb "github.com/AcceleratedCloudStorage/acs-sdk-go/generated"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

// ACSClient wraps the gRPC connection and client for ObjectStorageCache.
// It provides high-level operations for interacting with the ACS service.
type ACSClient struct {
	client  pb.ObjectStorageCacheClient
	conn    *grpc.ClientConn
	retry   RetryConfig
	session *Session
}

// Ensure compliation
var _ = embed.FS{}

//go:embed internal/ca-chain.pem
var embeddedCACert []byte

// loadClientTLSCredentials loads the CA certificates from the embedded file and returns
// the TransportCredentials with the loaded certificates.
func loadClientTLSCredentials() (credentials.TransportCredentials, error) {
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(embeddedCACert) {
		log.Fatalf("Failed to append CA certificates")
	}

	tlsConfig := &tls.Config{
		RootCAs:    certPool,
		MinVersion: tls.VersionTLS12,
	}

	return credentials.NewTLS(tlsConfig), nil
}

// NewClient initializes a new gRPC client with authentication.
// It establishes a secure connection to the ACS service, loads credentials,
// and performs initial authentication.
func NewClient(session *Session) (*ACSClient, error) {
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
		client:  pb.NewObjectStorageCacheClient(conn),
		conn:    conn,
		retry:   DefaultRetryConfig,
		session: session, // Store the session
	}

	// Load credentials from disk
	serviceCreds, err := loadACSCredentials()
	if err != nil {
		return nil, fmt.Errorf("failed to load credentials: %v", err)
	}

	// Prepare authentication request
	authReq := &pb.AuthRequest{
		AccessKeyId:     serviceCreds.AccessKeyID,
		SecretAccessKey: serviceCreds.SecretAccessKey,
	}

	// Add region if provided in session
	if session != nil && session.Region != "" {
		authReq.Region = &session.Region
	} else {
		DefaultRegion := "us-east-1"
		authReq.Region = &DefaultRegion
	}

	// Perform authentication
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = client.client.Authenticate(ctx, authReq)
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

