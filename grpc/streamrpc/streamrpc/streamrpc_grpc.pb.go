// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.18.1
// source: streamrpc.proto

package streamrpc

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

// MyServiceClient is the client API for MyService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MyServiceClient interface {
	MyMethod(ctx context.Context, opts ...grpc.CallOption) (MyService_MyMethodClient, error)
}

type myServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMyServiceClient(cc grpc.ClientConnInterface) MyServiceClient {
	return &myServiceClient{cc}
}

func (c *myServiceClient) MyMethod(ctx context.Context, opts ...grpc.CallOption) (MyService_MyMethodClient, error) {
	stream, err := c.cc.NewStream(ctx, &MyService_ServiceDesc.Streams[0], "/streamrpc.MyService/MyMethod", opts...)
	if err != nil {
		return nil, err
	}
	x := &myServiceMyMethodClient{stream}
	return x, nil
}

type MyService_MyMethodClient interface {
	Send(*MyMessage) error
	Recv() (*MyMessage, error)
	grpc.ClientStream
}

type myServiceMyMethodClient struct {
	grpc.ClientStream
}

func (x *myServiceMyMethodClient) Send(m *MyMessage) error {
	return x.ClientStream.SendMsg(m)
}

func (x *myServiceMyMethodClient) Recv() (*MyMessage, error) {
	m := new(MyMessage)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// MyServiceServer is the server API for MyService service.
// All implementations must embed UnimplementedMyServiceServer
// for forward compatibility
type MyServiceServer interface {
	MyMethod(MyService_MyMethodServer) error
	mustEmbedUnimplementedMyServiceServer()
}

// UnimplementedMyServiceServer must be embedded to have forward compatible implementations.
type UnimplementedMyServiceServer struct {
}

func (UnimplementedMyServiceServer) MyMethod(MyService_MyMethodServer) error {
	return status.Errorf(codes.Unimplemented, "method MyMethod not implemented")
}
func (UnimplementedMyServiceServer) mustEmbedUnimplementedMyServiceServer() {}

// UnsafeMyServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MyServiceServer will
// result in compilation errors.
type UnsafeMyServiceServer interface {
	mustEmbedUnimplementedMyServiceServer()
}

func RegisterMyServiceServer(s grpc.ServiceRegistrar, srv MyServiceServer) {
	s.RegisterService(&MyService_ServiceDesc, srv)
}

func _MyService_MyMethod_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(MyServiceServer).MyMethod(&myServiceMyMethodServer{stream})
}

type MyService_MyMethodServer interface {
	Send(*MyMessage) error
	Recv() (*MyMessage, error)
	grpc.ServerStream
}

type myServiceMyMethodServer struct {
	grpc.ServerStream
}

func (x *myServiceMyMethodServer) Send(m *MyMessage) error {
	return x.ServerStream.SendMsg(m)
}

func (x *myServiceMyMethodServer) Recv() (*MyMessage, error) {
	m := new(MyMessage)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// MyService_ServiceDesc is the grpc.ServiceDesc for MyService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MyService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "streamrpc.MyService",
	HandlerType: (*MyServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "MyMethod",
			Handler:       _MyService_MyMethod_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "streamrpc.proto",
}

// StreamRpcFileServiceClient is the client API for StreamRpcFileService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StreamRpcFileServiceClient interface {
	SendFile(ctx context.Context, in *StreamFileRequest, opts ...grpc.CallOption) (StreamRpcFileService_SendFileClient, error)
}

type streamRpcFileServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewStreamRpcFileServiceClient(cc grpc.ClientConnInterface) StreamRpcFileServiceClient {
	return &streamRpcFileServiceClient{cc}
}

func (c *streamRpcFileServiceClient) SendFile(ctx context.Context, in *StreamFileRequest, opts ...grpc.CallOption) (StreamRpcFileService_SendFileClient, error) {
	stream, err := c.cc.NewStream(ctx, &StreamRpcFileService_ServiceDesc.Streams[0], "/streamrpc.StreamRpcFileService/SendFile", opts...)
	if err != nil {
		return nil, err
	}
	x := &streamRpcFileServiceSendFileClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type StreamRpcFileService_SendFileClient interface {
	Recv() (*StreamFileReply, error)
	grpc.ClientStream
}

type streamRpcFileServiceSendFileClient struct {
	grpc.ClientStream
}

func (x *streamRpcFileServiceSendFileClient) Recv() (*StreamFileReply, error) {
	m := new(StreamFileReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// StreamRpcFileServiceServer is the server API for StreamRpcFileService service.
// All implementations must embed UnimplementedStreamRpcFileServiceServer
// for forward compatibility
type StreamRpcFileServiceServer interface {
	SendFile(*StreamFileRequest, StreamRpcFileService_SendFileServer) error
	mustEmbedUnimplementedStreamRpcFileServiceServer()
}

// UnimplementedStreamRpcFileServiceServer must be embedded to have forward compatible implementations.
type UnimplementedStreamRpcFileServiceServer struct {
}

func (UnimplementedStreamRpcFileServiceServer) SendFile(*StreamFileRequest, StreamRpcFileService_SendFileServer) error {
	return status.Errorf(codes.Unimplemented, "method SendFile not implemented")
}
func (UnimplementedStreamRpcFileServiceServer) mustEmbedUnimplementedStreamRpcFileServiceServer() {}

// UnsafeStreamRpcFileServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StreamRpcFileServiceServer will
// result in compilation errors.
type UnsafeStreamRpcFileServiceServer interface {
	mustEmbedUnimplementedStreamRpcFileServiceServer()
}

func RegisterStreamRpcFileServiceServer(s grpc.ServiceRegistrar, srv StreamRpcFileServiceServer) {
	s.RegisterService(&StreamRpcFileService_ServiceDesc, srv)
}

func _StreamRpcFileService_SendFile_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(StreamFileRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(StreamRpcFileServiceServer).SendFile(m, &streamRpcFileServiceSendFileServer{stream})
}

type StreamRpcFileService_SendFileServer interface {
	Send(*StreamFileReply) error
	grpc.ServerStream
}

type streamRpcFileServiceSendFileServer struct {
	grpc.ServerStream
}

func (x *streamRpcFileServiceSendFileServer) Send(m *StreamFileReply) error {
	return x.ServerStream.SendMsg(m)
}

// StreamRpcFileService_ServiceDesc is the grpc.ServiceDesc for StreamRpcFileService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var StreamRpcFileService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "streamrpc.StreamRpcFileService",
	HandlerType: (*StreamRpcFileServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SendFile",
			Handler:       _StreamRpcFileService_SendFile_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "streamrpc.proto",
}

// StreamRpcServiceClient is the client API for StreamRpcService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StreamRpcServiceClient interface {
	SayHelloWorld(ctx context.Context, in *StreamRequest, opts ...grpc.CallOption) (StreamRpcService_SayHelloWorldClient, error)
}

type streamRpcServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewStreamRpcServiceClient(cc grpc.ClientConnInterface) StreamRpcServiceClient {
	return &streamRpcServiceClient{cc}
}

func (c *streamRpcServiceClient) SayHelloWorld(ctx context.Context, in *StreamRequest, opts ...grpc.CallOption) (StreamRpcService_SayHelloWorldClient, error) {
	stream, err := c.cc.NewStream(ctx, &StreamRpcService_ServiceDesc.Streams[0], "/streamrpc.StreamRpcService/SayHelloWorld", opts...)
	if err != nil {
		return nil, err
	}
	x := &streamRpcServiceSayHelloWorldClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type StreamRpcService_SayHelloWorldClient interface {
	Recv() (*StreamReply, error)
	grpc.ClientStream
}

type streamRpcServiceSayHelloWorldClient struct {
	grpc.ClientStream
}

func (x *streamRpcServiceSayHelloWorldClient) Recv() (*StreamReply, error) {
	m := new(StreamReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// StreamRpcServiceServer is the server API for StreamRpcService service.
// All implementations must embed UnimplementedStreamRpcServiceServer
// for forward compatibility
type StreamRpcServiceServer interface {
	SayHelloWorld(*StreamRequest, StreamRpcService_SayHelloWorldServer) error
	mustEmbedUnimplementedStreamRpcServiceServer()
}

// UnimplementedStreamRpcServiceServer must be embedded to have forward compatible implementations.
type UnimplementedStreamRpcServiceServer struct {
}

func (UnimplementedStreamRpcServiceServer) SayHelloWorld(*StreamRequest, StreamRpcService_SayHelloWorldServer) error {
	return status.Errorf(codes.Unimplemented, "method SayHelloWorld not implemented")
}
func (UnimplementedStreamRpcServiceServer) mustEmbedUnimplementedStreamRpcServiceServer() {}

// UnsafeStreamRpcServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StreamRpcServiceServer will
// result in compilation errors.
type UnsafeStreamRpcServiceServer interface {
	mustEmbedUnimplementedStreamRpcServiceServer()
}

func RegisterStreamRpcServiceServer(s grpc.ServiceRegistrar, srv StreamRpcServiceServer) {
	s.RegisterService(&StreamRpcService_ServiceDesc, srv)
}

func _StreamRpcService_SayHelloWorld_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(StreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(StreamRpcServiceServer).SayHelloWorld(m, &streamRpcServiceSayHelloWorldServer{stream})
}

type StreamRpcService_SayHelloWorldServer interface {
	Send(*StreamReply) error
	grpc.ServerStream
}

type streamRpcServiceSayHelloWorldServer struct {
	grpc.ServerStream
}

func (x *streamRpcServiceSayHelloWorldServer) Send(m *StreamReply) error {
	return x.ServerStream.SendMsg(m)
}

// StreamRpcService_ServiceDesc is the grpc.ServiceDesc for StreamRpcService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var StreamRpcService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "streamrpc.StreamRpcService",
	HandlerType: (*StreamRpcServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SayHelloWorld",
			Handler:       _StreamRpcService_SayHelloWorld_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "streamrpc.proto",
}