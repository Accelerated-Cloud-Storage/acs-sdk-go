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
	acsClient, err := client.NewClient(&client.Session{Region: "us-east-1"})
	if err != nil {
		panic(err)
	}
	defer acsClient.Close()

	// Create a new bucket
	bucketName := fmt.Sprintf("my-bucket-%d", time.Now().UnixNano())
	err = acsClient.CreateBucket(context, bucketName)
	if err != nil {
		panic(err)
	}
	defer acsClient.DeleteBucket(context, bucketName)

	fmt.Println("Bucket created:", bucketName)
	// Create a new object
	objectName := "my-object"
	objectData := []byte("Hello, World!")
	err = acsClient.PutObject(context, bucketName, objectName, objectData)
	if err != nil {
		panic(err)
	}

	fmt.Println("Object created:", objectName)
	// Get the object
	data, err := acsClient.GetObject(context, bucketName, objectName)
	if err != nil {
		panic(err)
	}
	println(string(data))
	fmt.Println("Object retrieved:", objectName)

	// Delete the object
	err = acsClient.DeleteObject(context, bucketName, objectName)
	if err != nil {
		panic(err)
	}
	fmt.Println("Object deleted:", objectName)
	// Delete the bucket
	err = acsClient.DeleteBucket(context, bucketName)
	if err != nil {
		panic(err)
	}
}