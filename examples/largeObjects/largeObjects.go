// Copyright 2025 Accelerated Cloud Storage Corporation. All Rights Reserved.
package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/AcceleratedCloudStorage/acs-sdk-go/client"
)

func main() {
	context := context.Background()

	// Create a new GRPC client
	start := time.Now()
	acsClient, err := client.NewClient(&client.Session{Region: "us-east-1"})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Client creation took: %v\n", time.Since(start))
	defer acsClient.Close()

	// Create a new bucket
	bucketName := fmt.Sprintf("my-bucket-%d", time.Now().UnixNano())
	start = time.Now()
	err = acsClient.CreateBucket(context, bucketName)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Bucket creation took: %v\n", time.Since(start))

	// Create a new object
	objectName := fmt.Sprintf("my-object-%d", time.Now().UnixNano())
	// Create 10GB of data
	fmt.Println("Generating 10GB of test data...")
	dataStart := time.Now()
	objectData := make([]byte, 10*1024*1024*1024) // 10 GB
	for i := range objectData {
		objectData[i] = byte(rand.Intn(256)) // Fill with random bytes
	}
	fmt.Printf("Data generation took: %v\n", time.Since(dataStart))

	// Put object
	start = time.Now()
	err = acsClient.PutObject(context, bucketName, objectName, objectData)
	if err != nil {
		panic(err)
	}
	fmt.Printf("PutObject operation took: %v\n", time.Since(start))

	// Get object
	start = time.Now()
	data, err := acsClient.GetObject(context, bucketName, objectName)
	if err != nil {
		panic(err)
	}
	fmt.Printf("GetObject operation took: %v\n", time.Since(start))
	fmt.Printf("Retrieved data size: %d bytes\n", len(data))

	// Get object with range - retrieve only the first 1MB
	fmt.Println("\nPerforming ranged GetObject operation...")
	start = time.Now()
	rangeData, err := acsClient.GetObject(context, bucketName, objectName, client.WithRange("bytes=0-1048575"))
	if err != nil {
		panic(err)
	}
	fmt.Printf("Ranged GetObject operation took: %v\n", time.Since(start))
	fmt.Printf("Retrieved range data size: %d bytes\n", len(rangeData))

	// Get object with range - retrieve 1MB from the middle of the file
	fmt.Println("\nPerforming another ranged GetObject operation (middle of file)...")
	start = time.Now()
	middleOffset := 5 * 1024 * 1024 * 1024 // 5GB offset
	rangeMiddle, err := acsClient.GetObject(context, bucketName, objectName,
		client.WithRange(fmt.Sprintf("bytes=%d-%d", middleOffset, middleOffset+1048575)))
	if err != nil {
		panic(err)
	}
	fmt.Printf("Middle ranged GetObject operation took: %v\n", time.Since(start))
	fmt.Printf("Retrieved middle range data size: %d bytes\n", len(rangeMiddle))

	// Delete object
	start = time.Now()
	err = acsClient.DeleteObject(context, bucketName, objectName)
	if err != nil {
		panic(err)
	}
	fmt.Printf("DeleteObject operation took: %v\n", time.Since(start))

	// Delete bucket
	start = time.Now()
	err = acsClient.DeleteBucket(context, bucketName)
	if err != nil {
		panic(err)
	}
	fmt.Printf("DeleteBucket operation took: %v\n", time.Since(start))
}
