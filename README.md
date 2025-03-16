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
[![Website](https://img.shields.io/badge/Website-Console-blue)](https://acceleratedcloudstorage.com) [![API Reference](https://img.shields.io/badge/API-Reference-blue.svg)](https://pkg.go.dev/github.com/AcceleratedCloudStorage/acs-sdk-go) [![Demo](https://img.shields.io/badge/Demo-Videos-blue.svg)](https://www.youtube.com/@AcceleratedCloudStorageSales) 

To get started working with the SDK setup your project for Go modules, and retrieve the SDK dependencies with `go get`. This example shows how you can use the SDK to make an API request using the SDK's client.

#### Setup credientials 
Get your your credentials from the console on the [website](https://acceleratedcloudstorage.io).

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
Check out the example folder. 

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