# ACS SDK for Go
The Go SDK for Accelerated Cloud Storage's Object Storage offering. 

`acs-sdk-go` is the ACS SDK for the Go programming language.

The SDK requires a minimum version of `Go 1.23.4`.

Check out the [Release Notes] for information about the latest bug fixes, updates, and features added to the SDK.

Jump To:
* [Getting Started](#getting-started)
* [Getting Help](#getting-help)

### Go version support policy

The SDK follows the upstream [release policy](https://go.dev/doc/devel/release#policy)
with an additional six months of support for the most recently deprecated
language version.

**ACS reserves the right to drop support for unsupported Go versions earlier to
address critical security issues.**

## Getting started
[![Website](https://img.shields.io/badge/Website-Console-blue)](https://acceleratedcloudstorage.io) [![API Reference](https://img.shields.io/badge/API-Reference-blue.svg)](https://pkg.go.dev/github.com/AcceleratedCloudStorage/acs-sdk-go) [![Demo](https://img.shields.io/badge/Demo-Videos-blue.svg)](https://www.youtube.com/@AcceleratedCloudStorageSales) 

To get started working with the SDK setup your project for Go modules, and retrieve the SDK dependencies with `go get`. This example shows how you can use the SDK to make an API request using the SDK's client.

#### Get credentials

Get your credentials and setup payments from the console on the [website](https://acceleratedcloudstorage.io).

Next, set up credentials (in e.g. ``~/.acs/credentials``):
```
default:
    access_key_id = YOUR_KEY
    secret_access_key = YOUR_SECRET
```
Note: You can include multiple profiles and set them using the ACS_PROFILE environment variable. See the examples/config folder for a sample file. 

#### Initialize Project
```sh
$ mkdir ~/helloacs
$ cd ~/helloacs
$ go mod init helloacs
```
#### Add SDK Dependencies
```sh
$ go get github.com/AcceleratedCloudStorage/acs-sdk-go/client
```

#### Write Code
You can either use the client for an interface similar to the AWS SDK or a FUSE mount for a file system interface. Check out the example folder for more details.

## Share bucket

You can also bring your existing buckets into the service by setting a bucket policy and then sharing the bucket with the service.

### Step 1: Setting a bucket policy

Here is the AWS reference guide for [bucket policies](https://docs.aws.amazon.com/AmazonS3/latest/userguide/add-bucket-policy.html). You can set the following bucket policy through the AWS Console or SDK to enable ACS to access it.

```
{
"Version": "2012-10-17",
   "Statement": [
    {
     "Sid": "AllowUserFullAccess", 
     "Effect": "Allow",
     "Principal": {
      "AWS": "arn:aws:iam::160885293701:root"
     },
     "Action": [
      "s3:*"
     ],
     "Resource": [
      "arn:aws:s3:::BUCKETNAME",
      "arn:aws:s3:::BUCKETNAME/*"
     ]
    }
   ]
}
```

### Step 2: Notify ACS of this newly shared bucket

```
// Create a new client
acsClient, err := client.NewClient(&client.Session{Region: "us-east-1"})
defer acsClient.Close()
// Share a bucket
err = client.ShareBucket(context, BUCKETNAME)
```

## Getting Help

Please use these community resources for getting help. 

### Feedback

If you encounter a bug with the ACS SDK for Go we would like to hear about it.
Search the [existing issues][Issues] and see if others are also experiencing the same issue before opening a new issue. Please include the version of ACS SDK for Go, Go language, and OS you’re using. Please also include reproduction case when appropriate. Keeping the list of open issues lean will help us respond in a timely manner.

### Discussion  

We have a discussion forum where you can read about announcements, product ideas, partcipate in Q&A. Here is a link to the [discussion].

### Contact us 

Email us at sales@acceleratedcloudstorage.com if you have any further questions or concerns.  

[Dep]: https://github.com/golang/dep
[Issues]: https://github.com/AcceleratedCloudStorage/acs-sdk-go/issues
[Discussion]: https://github.com/AcceleratedCloudStorage/acs-sdk-go/discussions
[Release Notes]: https://github.com/AcceleratedCloudStorage/acs-sdk-go/blob/main/CHANGELOG.md