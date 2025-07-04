// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v6.31.1
// source: fileshare.proto

package gen

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// UploadServiceClient is the client API for UploadService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UploadServiceClient interface {
	PreUpload(ctx context.Context, in *UploadRequest, opts ...grpc.CallOption) (*UploadTask, error)
	Upload(ctx context.Context, opts ...grpc.CallOption) (UploadService_UploadClient, error)
}

type uploadServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewUploadServiceClient(cc grpc.ClientConnInterface) UploadServiceClient {
	return &uploadServiceClient{cc}
}

func (c *uploadServiceClient) PreUpload(ctx context.Context, in *UploadRequest, opts ...grpc.CallOption) (*UploadTask, error) {
	out := new(UploadTask)
	err := c.cc.Invoke(ctx, "/fileshare.UploadService/PreUpload", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uploadServiceClient) Upload(ctx context.Context, opts ...grpc.CallOption) (UploadService_UploadClient, error) {
	stream, err := c.cc.NewStream(ctx, &UploadService_ServiceDesc.Streams[0], "/fileshare.UploadService/Upload", opts...)
	if err != nil {
		return nil, err
	}
	x := &uploadServiceUploadClient{stream}
	return x, nil
}

type UploadService_UploadClient interface {
	Send(*FileChunk) error
	CloseAndRecv() (*UploadStatus, error)
	grpc.ClientStream
}

type uploadServiceUploadClient struct {
	grpc.ClientStream
}

func (x *uploadServiceUploadClient) Send(m *FileChunk) error {
	return x.ClientStream.SendMsg(m)
}

func (x *uploadServiceUploadClient) CloseAndRecv() (*UploadStatus, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(UploadStatus)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// UploadServiceServer is the server API for UploadService service.
// All implementations must embed UnimplementedUploadServiceServer
// for forward compatibility
type UploadServiceServer interface {
	PreUpload(context.Context, *UploadRequest) (*UploadTask, error)
	Upload(UploadService_UploadServer) error
	mustEmbedUnimplementedUploadServiceServer()
}

// UnimplementedUploadServiceServer must be embedded to have forward compatible implementations.
type UnimplementedUploadServiceServer struct {
}

func (UnimplementedUploadServiceServer) PreUpload(context.Context, *UploadRequest) (*UploadTask, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PreUpload not implemented")
}
func (UnimplementedUploadServiceServer) Upload(UploadService_UploadServer) error {
	return status.Errorf(codes.Unimplemented, "method Upload not implemented")
}
func (UnimplementedUploadServiceServer) mustEmbedUnimplementedUploadServiceServer() {}

// UnsafeUploadServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UploadServiceServer will
// result in compilation errors.
type UnsafeUploadServiceServer interface {
	mustEmbedUnimplementedUploadServiceServer()
}

func RegisterUploadServiceServer(s grpc.ServiceRegistrar, srv UploadServiceServer) {
	s.RegisterService(&UploadService_ServiceDesc, srv)
}

func _UploadService_PreUpload_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UploadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UploadServiceServer).PreUpload(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fileshare.UploadService/PreUpload",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UploadServiceServer).PreUpload(ctx, req.(*UploadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UploadService_Upload_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(UploadServiceServer).Upload(&uploadServiceUploadServer{stream})
}

type UploadService_UploadServer interface {
	SendAndClose(*UploadStatus) error
	Recv() (*FileChunk, error)
	grpc.ServerStream
}

type uploadServiceUploadServer struct {
	grpc.ServerStream
}

func (x *uploadServiceUploadServer) SendAndClose(m *UploadStatus) error {
	return x.ServerStream.SendMsg(m)
}

func (x *uploadServiceUploadServer) Recv() (*FileChunk, error) {
	m := new(FileChunk)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// UploadService_ServiceDesc is the grpc.ServiceDesc for UploadService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UploadService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "fileshare.UploadService",
	HandlerType: (*UploadServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PreUpload",
			Handler:    _UploadService_PreUpload_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Upload",
			Handler:       _UploadService_Upload_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "fileshare.proto",
}

// DownloadServiceClient is the client API for DownloadService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DownloadServiceClient interface {
	PreDownloadWithCode(ctx context.Context, in *ShareLink, opts ...grpc.CallOption) (*DownloadSummary, error)
	PreDownload(ctx context.Context, in *DownloadRequest, opts ...grpc.CallOption) (*DownloadSummary, error)
	Download(ctx context.Context, in *DownloadTask, opts ...grpc.CallOption) (DownloadService_DownloadClient, error)
}

type downloadServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDownloadServiceClient(cc grpc.ClientConnInterface) DownloadServiceClient {
	return &downloadServiceClient{cc}
}

func (c *downloadServiceClient) PreDownloadWithCode(ctx context.Context, in *ShareLink, opts ...grpc.CallOption) (*DownloadSummary, error) {
	out := new(DownloadSummary)
	err := c.cc.Invoke(ctx, "/fileshare.DownloadService/PreDownloadWithCode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *downloadServiceClient) PreDownload(ctx context.Context, in *DownloadRequest, opts ...grpc.CallOption) (*DownloadSummary, error) {
	out := new(DownloadSummary)
	err := c.cc.Invoke(ctx, "/fileshare.DownloadService/PreDownload", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *downloadServiceClient) Download(ctx context.Context, in *DownloadTask, opts ...grpc.CallOption) (DownloadService_DownloadClient, error) {
	stream, err := c.cc.NewStream(ctx, &DownloadService_ServiceDesc.Streams[0], "/fileshare.DownloadService/Download", opts...)
	if err != nil {
		return nil, err
	}
	x := &downloadServiceDownloadClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type DownloadService_DownloadClient interface {
	Recv() (*FileChunk, error)
	grpc.ClientStream
}

type downloadServiceDownloadClient struct {
	grpc.ClientStream
}

func (x *downloadServiceDownloadClient) Recv() (*FileChunk, error) {
	m := new(FileChunk)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// DownloadServiceServer is the server API for DownloadService service.
// All implementations must embed UnimplementedDownloadServiceServer
// for forward compatibility
type DownloadServiceServer interface {
	PreDownloadWithCode(context.Context, *ShareLink) (*DownloadSummary, error)
	PreDownload(context.Context, *DownloadRequest) (*DownloadSummary, error)
	Download(*DownloadTask, DownloadService_DownloadServer) error
	mustEmbedUnimplementedDownloadServiceServer()
}

// UnimplementedDownloadServiceServer must be embedded to have forward compatible implementations.
type UnimplementedDownloadServiceServer struct {
}

func (UnimplementedDownloadServiceServer) PreDownloadWithCode(context.Context, *ShareLink) (*DownloadSummary, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PreDownloadWithCode not implemented")
}
func (UnimplementedDownloadServiceServer) PreDownload(context.Context, *DownloadRequest) (*DownloadSummary, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PreDownload not implemented")
}
func (UnimplementedDownloadServiceServer) Download(*DownloadTask, DownloadService_DownloadServer) error {
	return status.Errorf(codes.Unimplemented, "method Download not implemented")
}
func (UnimplementedDownloadServiceServer) mustEmbedUnimplementedDownloadServiceServer() {}

// UnsafeDownloadServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DownloadServiceServer will
// result in compilation errors.
type UnsafeDownloadServiceServer interface {
	mustEmbedUnimplementedDownloadServiceServer()
}

func RegisterDownloadServiceServer(s grpc.ServiceRegistrar, srv DownloadServiceServer) {
	s.RegisterService(&DownloadService_ServiceDesc, srv)
}

func _DownloadService_PreDownloadWithCode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShareLink)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DownloadServiceServer).PreDownloadWithCode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fileshare.DownloadService/PreDownloadWithCode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DownloadServiceServer).PreDownloadWithCode(ctx, req.(*ShareLink))
	}
	return interceptor(ctx, in, info, handler)
}

func _DownloadService_PreDownload_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DownloadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DownloadServiceServer).PreDownload(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fileshare.DownloadService/PreDownload",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DownloadServiceServer).PreDownload(ctx, req.(*DownloadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DownloadService_Download_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(DownloadTask)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(DownloadServiceServer).Download(m, &downloadServiceDownloadServer{stream})
}

type DownloadService_DownloadServer interface {
	Send(*FileChunk) error
	grpc.ServerStream
}

type downloadServiceDownloadServer struct {
	grpc.ServerStream
}

func (x *downloadServiceDownloadServer) Send(m *FileChunk) error {
	return x.ServerStream.SendMsg(m)
}

// DownloadService_ServiceDesc is the grpc.ServiceDesc for DownloadService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DownloadService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "fileshare.DownloadService",
	HandlerType: (*DownloadServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PreDownloadWithCode",
			Handler:    _DownloadService_PreDownloadWithCode_Handler,
		},
		{
			MethodName: "PreDownload",
			Handler:    _DownloadService_PreDownload_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Download",
			Handler:       _DownloadService_Download_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "fileshare.proto",
}

// ShareLinkServiceClient is the client API for ShareLinkService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ShareLinkServiceClient interface {
	GenerateLink(ctx context.Context, in *ShareLinkRequest, opts ...grpc.CallOption) (*ShareLinkResponse, error)
}

type shareLinkServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewShareLinkServiceClient(cc grpc.ClientConnInterface) ShareLinkServiceClient {
	return &shareLinkServiceClient{cc}
}

func (c *shareLinkServiceClient) GenerateLink(ctx context.Context, in *ShareLinkRequest, opts ...grpc.CallOption) (*ShareLinkResponse, error) {
	out := new(ShareLinkResponse)
	err := c.cc.Invoke(ctx, "/fileshare.ShareLinkService/GenerateLink", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ShareLinkServiceServer is the server API for ShareLinkService service.
// All implementations must embed UnimplementedShareLinkServiceServer
// for forward compatibility
type ShareLinkServiceServer interface {
	GenerateLink(context.Context, *ShareLinkRequest) (*ShareLinkResponse, error)
	mustEmbedUnimplementedShareLinkServiceServer()
}

// UnimplementedShareLinkServiceServer must be embedded to have forward compatible implementations.
type UnimplementedShareLinkServiceServer struct {
}

func (UnimplementedShareLinkServiceServer) GenerateLink(context.Context, *ShareLinkRequest) (*ShareLinkResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GenerateLink not implemented")
}
func (UnimplementedShareLinkServiceServer) mustEmbedUnimplementedShareLinkServiceServer() {}

// UnsafeShareLinkServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ShareLinkServiceServer will
// result in compilation errors.
type UnsafeShareLinkServiceServer interface {
	mustEmbedUnimplementedShareLinkServiceServer()
}

func RegisterShareLinkServiceServer(s grpc.ServiceRegistrar, srv ShareLinkServiceServer) {
	s.RegisterService(&ShareLinkService_ServiceDesc, srv)
}

func _ShareLinkService_GenerateLink_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShareLinkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShareLinkServiceServer).GenerateLink(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fileshare.ShareLinkService/GenerateLink",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShareLinkServiceServer).GenerateLink(ctx, req.(*ShareLinkRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ShareLinkService_ServiceDesc is the grpc.ServiceDesc for ShareLinkService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ShareLinkService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "fileshare.ShareLinkService",
	HandlerType: (*ShareLinkServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GenerateLink",
			Handler:    _ShareLinkService_GenerateLink_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "fileshare.proto",
}
