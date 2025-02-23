// Copyright 2025 Accelerated Cloud Storage Corporation. All Rights Reserved.
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/AcceleratedCloudStorage/acs-sdk-go/client"
)

func main() {
	context := context.Background()

	// Create a new GRPC client
	acsClient, err := client.NewClient()
	if err != nil {
		panic(err)
	}
	defer acsClient.Close()

	// Create a new bucket
	bucketName := fmt.Sprintf("my-bucket-%d", time.Now().UnixNano())
	err = acsClient.CreateBucket(context, bucketName, "us-east-1")
	if err != nil {
		panic(err)
	}
	defer acsClient.DeleteBucket(context, bucketName)

	// Create a new object
	objectName := "my-object"
	objectData := []byte("Hello, World!")
	err = acsClient.PutObject(context, bucketName, objectName, objectData)
	if err != nil {
		panic(err)
	}

	// Get the object
	data, err := acsClient.GetObject(context, bucketName, objectName)
	if err != nil {
		panic(err)
	}
	println(string(data))

	// Delete the object
	err = acsClient.DeleteObject(context, bucketName, objectName)
	if err != nil {
		panic(err)
	}

	// Delete the bucket
	err = acsClient.DeleteBucket(context, bucketName)
	if err != nil {
		panic(err)
	}
}