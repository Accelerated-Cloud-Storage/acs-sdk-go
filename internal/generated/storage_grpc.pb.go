// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: storage.proto

package generated

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	ObjectStorageCache_CreateBucket_FullMethodName  = "/proto.ObjectStorageCache/CreateBucket"
	ObjectStorageCache_DeleteBucket_FullMethodName  = "/proto.ObjectStorageCache/DeleteBucket"
	ObjectStorageCache_ListBuckets_FullMethodName   = "/proto.ObjectStorageCache/ListBuckets"
	ObjectStorageCache_HeadBucket_FullMethodName    = "/proto.ObjectStorageCache/HeadBucket"
	ObjectStorageCache_PutObject_FullMethodName     = "/proto.ObjectStorageCache/PutObject"
	ObjectStorageCache_GetObject_FullMethodName     = "/proto.ObjectStorageCache/GetObject"
	ObjectStorageCache_DeleteObject_FullMethodName  = "/proto.ObjectStorageCache/DeleteObject"
	ObjectStorageCache_DeleteObjects_FullMethodName = "/proto.ObjectStorageCache/DeleteObjects"
	ObjectStorageCache_CopyObject_FullMethodName    = "/proto.ObjectStorageCache/CopyObject"
	ObjectStorageCache_HeadObject_FullMethodName    = "/proto.ObjectStorageCache/HeadObject"
	ObjectStorageCache_ListObjects_FullMethodName   = "/proto.ObjectStorageCache/ListObjects"
	ObjectStorageCache_Authenticate_FullMethodName  = "/proto.ObjectStorageCache/Authenticate"
	ObjectStorageCache_RotateKey_FullMethodName     = "/proto.ObjectStorageCache/RotateKey"
	ObjectStorageCache_ShareBucket_FullMethodName   = "/proto.ObjectStorageCache/ShareBucket"
)

// ObjectStorageCacheClient is the client API for ObjectStorageCache service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ObjectStorageCacheClient interface {
	// Bucket operations
	CreateBucket(ctx context.Context, in *CreateBucketRequest, opts ...grpc.CallOption) (*CreateBucketResponse, error)
	DeleteBucket(ctx context.Context, in *DeleteBucketRequest, opts ...grpc.CallOption) (*DeleteBucketResponse, error)
	ListBuckets(ctx context.Context, in *ListBucketsRequest, opts ...grpc.CallOption) (*ListBucketsResponse, error)
	HeadBucket(ctx context.Context, in *HeadBucketRequest, opts ...grpc.CallOption) (*HeadBucketResponse, error)
	// Object operations
	PutObject(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[PutObjectRequest, PutObjectResponse], error)
	GetObject(ctx context.Context, in *GetObjectRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[GetObjectResponse], error)
	DeleteObject(ctx context.Context, in *DeleteObjectRequest, opts ...grpc.CallOption) (*DeleteObjectResponse, error)
	DeleteObjects(ctx context.Context, in *DeleteObjectsRequest, opts ...grpc.CallOption) (*DeleteObjectsResponse, error)
	CopyObject(ctx context.Context, in *CopyObjectRequest, opts ...grpc.CallOption) (*CopyObjectResponse, error)
	HeadObject(ctx context.Context, in *HeadObjectRequest, opts ...grpc.CallOption) (*HeadObjectResponse, error)
	ListObjects(ctx context.Context, in *ListObjectsRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[ListObjectsResponse], error)
	// Configuration operations
	Authenticate(ctx context.Context, in *AuthRequest, opts ...grpc.CallOption) (*AuthResponse, error)
	RotateKey(ctx context.Context, in *RotateKeyRequest, opts ...grpc.CallOption) (*RotateKeyResponse, error)
	ShareBucket(ctx context.Context, in *ShareBucketRequest, opts ...grpc.CallOption) (*ShareBucketResponse, error)
}

type objectStorageCacheClient struct {
	cc grpc.ClientConnInterface
}

func NewObjectStorageCacheClient(cc grpc.ClientConnInterface) ObjectStorageCacheClient {
	return &objectStorageCacheClient{cc}
}

func (c *objectStorageCacheClient) CreateBucket(ctx context.Context, in *CreateBucketRequest, opts ...grpc.CallOption) (*CreateBucketResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateBucketResponse)
	err := c.cc.Invoke(ctx, ObjectStorageCache_CreateBucket_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *objectStorageCacheClient) DeleteBucket(ctx context.Context, in *DeleteBucketRequest, opts ...grpc.CallOption) (*DeleteBucketResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteBucketResponse)
	err := c.cc.Invoke(ctx, ObjectStorageCache_DeleteBucket_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *objectStorageCacheClient) ListBuckets(ctx context.Context, in *ListBucketsRequest, opts ...grpc.CallOption) (*ListBucketsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListBucketsResponse)
	err := c.cc.Invoke(ctx, ObjectStorageCache_ListBuckets_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *objectStorageCacheClient) HeadBucket(ctx context.Context, in *HeadBucketRequest, opts ...grpc.CallOption) (*HeadBucketResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(HeadBucketResponse)
	err := c.cc.Invoke(ctx, ObjectStorageCache_HeadBucket_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *objectStorageCacheClient) PutObject(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[PutObjectRequest, PutObjectResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &ObjectStorageCache_ServiceDesc.Streams[0], ObjectStorageCache_PutObject_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[PutObjectRequest, PutObjectResponse]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type ObjectStorageCache_PutObjectClient = grpc.ClientStreamingClient[PutObjectRequest, PutObjectResponse]

func (c *objectStorageCacheClient) GetObject(ctx context.Context, in *GetObjectRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[GetObjectResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &ObjectStorageCache_ServiceDesc.Streams[1], ObjectStorageCache_GetObject_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[GetObjectRequest, GetObjectResponse]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type ObjectStorageCache_GetObjectClient = grpc.ServerStreamingClient[GetObjectResponse]

func (c *objectStorageCacheClient) DeleteObject(ctx context.Context, in *DeleteObjectRequest, opts ...grpc.CallOption) (*DeleteObjectResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteObjectResponse)
	err := c.cc.Invoke(ctx, ObjectStorageCache_DeleteObject_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *objectStorageCacheClient) DeleteObjects(ctx context.Context, in *DeleteObjectsRequest, opts ...grpc.CallOption) (*DeleteObjectsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteObjectsResponse)
	err := c.cc.Invoke(ctx, ObjectStorageCache_DeleteObjects_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *objectStorageCacheClient) CopyObject(ctx context.Context, in *CopyObjectRequest, opts ...grpc.CallOption) (*CopyObjectResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CopyObjectResponse)
	err := c.cc.Invoke(ctx, ObjectStorageCache_CopyObject_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *objectStorageCacheClient) HeadObject(ctx context.Context, in *HeadObjectRequest, opts ...grpc.CallOption) (*HeadObjectResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(HeadObjectResponse)
	err := c.cc.Invoke(ctx, ObjectStorageCache_HeadObject_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *objectStorageCacheClient) ListObjects(ctx context.Context, in *ListObjectsRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[ListObjectsResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &ObjectStorageCache_ServiceDesc.Streams[2], ObjectStorageCache_ListObjects_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[ListObjectsRequest, ListObjectsResponse]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type ObjectStorageCache_ListObjectsClient = grpc.ServerStreamingClient[ListObjectsResponse]

func (c *objectStorageCacheClient) Authenticate(ctx context.Context, in *AuthRequest, opts ...grpc.CallOption) (*AuthResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AuthResponse)
	err := c.cc.Invoke(ctx, ObjectStorageCache_Authenticate_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *objectStorageCacheClient) RotateKey(ctx context.Context, in *RotateKeyRequest, opts ...grpc.CallOption) (*RotateKeyResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RotateKeyResponse)
	err := c.cc.Invoke(ctx, ObjectStorageCache_RotateKey_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *objectStorageCacheClient) ShareBucket(ctx context.Context, in *ShareBucketRequest, opts ...grpc.CallOption) (*ShareBucketResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ShareBucketResponse)
	err := c.cc.Invoke(ctx, ObjectStorageCache_ShareBucket_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ObjectStorageCacheServer is the server API for ObjectStorageCache service.
// All implementations must embed UnimplementedObjectStorageCacheServer
// for forward compatibility.
type ObjectStorageCacheServer interface {
	// Bucket operations
	CreateBucket(context.Context, *CreateBucketRequest) (*CreateBucketResponse, error)
	DeleteBucket(context.Context, *DeleteBucketRequest) (*DeleteBucketResponse, error)
	ListBuckets(context.Context, *ListBucketsRequest) (*ListBucketsResponse, error)
	HeadBucket(context.Context, *HeadBucketRequest) (*HeadBucketResponse, error)
	// Object operations
	PutObject(grpc.ClientStreamingServer[PutObjectRequest, PutObjectResponse]) error
	GetObject(*GetObjectRequest, grpc.ServerStreamingServer[GetObjectResponse]) error
	DeleteObject(context.Context, *DeleteObjectRequest) (*DeleteObjectResponse, error)
	DeleteObjects(context.Context, *DeleteObjectsRequest) (*DeleteObjectsResponse, error)
	CopyObject(context.Context, *CopyObjectRequest) (*CopyObjectResponse, error)
	HeadObject(context.Context, *HeadObjectRequest) (*HeadObjectResponse, error)
	ListObjects(*ListObjectsRequest, grpc.ServerStreamingServer[ListObjectsResponse]) error
	// Configuration operations
	Authenticate(context.Context, *AuthRequest) (*AuthResponse, error)
	RotateKey(context.Context, *RotateKeyRequest) (*RotateKeyResponse, error)
	ShareBucket(context.Context, *ShareBucketRequest) (*ShareBucketResponse, error)
	mustEmbedUnimplementedObjectStorageCacheServer()
}

// UnimplementedObjectStorageCacheServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedObjectStorageCacheServer struct{}

func (UnimplementedObjectStorageCacheServer) CreateBucket(context.Context, *CreateBucketRequest) (*CreateBucketResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateBucket not implemented")
}
func (UnimplementedObjectStorageCacheServer) DeleteBucket(context.Context, *DeleteBucketRequest) (*DeleteBucketResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteBucket not implemented")
}
func (UnimplementedObjectStorageCacheServer) ListBuckets(context.Context, *ListBucketsRequest) (*ListBucketsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListBuckets not implemented")
}
func (UnimplementedObjectStorageCacheServer) HeadBucket(context.Context, *HeadBucketRequest) (*HeadBucketResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HeadBucket not implemented")
}
func (UnimplementedObjectStorageCacheServer) PutObject(grpc.ClientStreamingServer[PutObjectRequest, PutObjectResponse]) error {
	return status.Errorf(codes.Unimplemented, "method PutObject not implemented")
}
func (UnimplementedObjectStorageCacheServer) GetObject(*GetObjectRequest, grpc.ServerStreamingServer[GetObjectResponse]) error {
	return status.Errorf(codes.Unimplemented, "method GetObject not implemented")
}
func (UnimplementedObjectStorageCacheServer) DeleteObject(context.Context, *DeleteObjectRequest) (*DeleteObjectResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteObject not implemented")
}
func (UnimplementedObjectStorageCacheServer) DeleteObjects(context.Context, *DeleteObjectsRequest) (*DeleteObjectsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteObjects not implemented")
}
func (UnimplementedObjectStorageCacheServer) CopyObject(context.Context, *CopyObjectRequest) (*CopyObjectResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CopyObject not implemented")
}
func (UnimplementedObjectStorageCacheServer) HeadObject(context.Context, *HeadObjectRequest) (*HeadObjectResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HeadObject not implemented")
}
func (UnimplementedObjectStorageCacheServer) ListObjects(*ListObjectsRequest, grpc.ServerStreamingServer[ListObjectsResponse]) error {
	return status.Errorf(codes.Unimplemented, "method ListObjects not implemented")
}
func (UnimplementedObjectStorageCacheServer) Authenticate(context.Context, *AuthRequest) (*AuthResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Authenticate not implemented")
}
func (UnimplementedObjectStorageCacheServer) RotateKey(context.Context, *RotateKeyRequest) (*RotateKeyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RotateKey not implemented")
}
func (UnimplementedObjectStorageCacheServer) ShareBucket(context.Context, *ShareBucketRequest) (*ShareBucketResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ShareBucket not implemented")
}
func (UnimplementedObjectStorageCacheServer) mustEmbedUnimplementedObjectStorageCacheServer() {}
func (UnimplementedObjectStorageCacheServer) testEmbeddedByValue()                            {}

// UnsafeObjectStorageCacheServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ObjectStorageCacheServer will
// result in compilation errors.
type UnsafeObjectStorageCacheServer interface {
	mustEmbedUnimplementedObjectStorageCacheServer()
}

func RegisterObjectStorageCacheServer(s grpc.ServiceRegistrar, srv ObjectStorageCacheServer) {
	// If the following call pancis, it indicates UnimplementedObjectStorageCacheServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&ObjectStorageCache_ServiceDesc, srv)
}

func _ObjectStorageCache_CreateBucket_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateBucketRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ObjectStorageCacheServer).CreateBucket(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ObjectStorageCache_CreateBucket_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ObjectStorageCacheServer).CreateBucket(ctx, req.(*CreateBucketRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ObjectStorageCache_DeleteBucket_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteBucketRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ObjectStorageCacheServer).DeleteBucket(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ObjectStorageCache_DeleteBucket_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ObjectStorageCacheServer).DeleteBucket(ctx, req.(*DeleteBucketRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ObjectStorageCache_ListBuckets_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListBucketsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ObjectStorageCacheServer).ListBuckets(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ObjectStorageCache_ListBuckets_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ObjectStorageCacheServer).ListBuckets(ctx, req.(*ListBucketsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ObjectStorageCache_HeadBucket_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HeadBucketRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ObjectStorageCacheServer).HeadBucket(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ObjectStorageCache_HeadBucket_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ObjectStorageCacheServer).HeadBucket(ctx, req.(*HeadBucketRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ObjectStorageCache_PutObject_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ObjectStorageCacheServer).PutObject(&grpc.GenericServerStream[PutObjectRequest, PutObjectResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type ObjectStorageCache_PutObjectServer = grpc.ClientStreamingServer[PutObjectRequest, PutObjectResponse]

func _ObjectStorageCache_GetObject_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(GetObjectRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ObjectStorageCacheServer).GetObject(m, &grpc.GenericServerStream[GetObjectRequest, GetObjectResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type ObjectStorageCache_GetObjectServer = grpc.ServerStreamingServer[GetObjectResponse]

func _ObjectStorageCache_DeleteObject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteObjectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ObjectStorageCacheServer).DeleteObject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ObjectStorageCache_DeleteObject_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ObjectStorageCacheServer).DeleteObject(ctx, req.(*DeleteObjectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ObjectStorageCache_DeleteObjects_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteObjectsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ObjectStorageCacheServer).DeleteObjects(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ObjectStorageCache_DeleteObjects_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ObjectStorageCacheServer).DeleteObjects(ctx, req.(*DeleteObjectsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ObjectStorageCache_CopyObject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CopyObjectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ObjectStorageCacheServer).CopyObject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ObjectStorageCache_CopyObject_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ObjectStorageCacheServer).CopyObject(ctx, req.(*CopyObjectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ObjectStorageCache_HeadObject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HeadObjectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ObjectStorageCacheServer).HeadObject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ObjectStorageCache_HeadObject_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ObjectStorageCacheServer).HeadObject(ctx, req.(*HeadObjectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ObjectStorageCache_ListObjects_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ListObjectsRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ObjectStorageCacheServer).ListObjects(m, &grpc.GenericServerStream[ListObjectsRequest, ListObjectsResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type ObjectStorageCache_ListObjectsServer = grpc.ServerStreamingServer[ListObjectsResponse]

func _ObjectStorageCache_Authenticate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ObjectStorageCacheServer).Authenticate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ObjectStorageCache_Authenticate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ObjectStorageCacheServer).Authenticate(ctx, req.(*AuthRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ObjectStorageCache_RotateKey_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RotateKeyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ObjectStorageCacheServer).RotateKey(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ObjectStorageCache_RotateKey_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ObjectStorageCacheServer).RotateKey(ctx, req.(*RotateKeyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ObjectStorageCache_ShareBucket_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShareBucketRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ObjectStorageCacheServer).ShareBucket(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ObjectStorageCache_ShareBucket_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ObjectStorageCacheServer).ShareBucket(ctx, req.(*ShareBucketRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ObjectStorageCache_ServiceDesc is the grpc.ServiceDesc for ObjectStorageCache service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ObjectStorageCache_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.ObjectStorageCache",
	HandlerType: (*ObjectStorageCacheServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateBucket",
			Handler:    _ObjectStorageCache_CreateBucket_Handler,
		},
		{
			MethodName: "DeleteBucket",
			Handler:    _ObjectStorageCache_DeleteBucket_Handler,
		},
		{
			MethodName: "ListBuckets",
			Handler:    _ObjectStorageCache_ListBuckets_Handler,
		},
		{
			MethodName: "HeadBucket",
			Handler:    _ObjectStorageCache_HeadBucket_Handler,
		},
		{
			MethodName: "DeleteObject",
			Handler:    _ObjectStorageCache_DeleteObject_Handler,
		},
		{
			MethodName: "DeleteObjects",
			Handler:    _ObjectStorageCache_DeleteObjects_Handler,
		},
		{
			MethodName: "CopyObject",
			Handler:    _ObjectStorageCache_CopyObject_Handler,
		},
		{
			MethodName: "HeadObject",
			Handler:    _ObjectStorageCache_HeadObject_Handler,
		},
		{
			MethodName: "Authenticate",
			Handler:    _ObjectStorageCache_Authenticate_Handler,
		},
		{
			MethodName: "RotateKey",
			Handler:    _ObjectStorageCache_RotateKey_Handler,
		},
		{
			MethodName: "ShareBucket",
			Handler:    _ObjectStorageCache_ShareBucket_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "PutObject",
			Handler:       _ObjectStorageCache_PutObject_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "GetObject",
			Handler:       _ObjectStorageCache_GetObject_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "ListObjects",
			Handler:       _ObjectStorageCache_ListObjects_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "storage.proto",
}
