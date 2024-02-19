// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: logs.proto

package logs

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

// LogServieClient is the client API for LogServie service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LogServieClient interface {
	WriteLog(ctx context.Context, in *LogRequest, opts ...grpc.CallOption) (*LogResponse, error)
}

type logServieClient struct {
	cc grpc.ClientConnInterface
}

func NewLogServieClient(cc grpc.ClientConnInterface) LogServieClient {
	return &logServieClient{cc}
}

func (c *logServieClient) WriteLog(ctx context.Context, in *LogRequest, opts ...grpc.CallOption) (*LogResponse, error) {
	out := new(LogResponse)
	err := c.cc.Invoke(ctx, "/logs.LogServie/WriteLog", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LogServieServer is the server API for LogServie service.
// All implementations must embed UnimplementedLogServieServer
// for forward compatibility
type LogServieServer interface {
	WriteLog(context.Context, *LogRequest) (*LogResponse, error)
	mustEmbedUnimplementedLogServieServer()
}

// UnimplementedLogServieServer must be embedded to have forward compatible implementations.
type UnimplementedLogServieServer struct {
}

func (UnimplementedLogServieServer) WriteLog(context.Context, *LogRequest) (*LogResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method WriteLog not implemented")
}
func (UnimplementedLogServieServer) mustEmbedUnimplementedLogServieServer() {}

// UnsafeLogServieServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LogServieServer will
// result in compilation errors.
type UnsafeLogServieServer interface {
	mustEmbedUnimplementedLogServieServer()
}

func RegisterLogServieServer(s grpc.ServiceRegistrar, srv LogServieServer) {
	s.RegisterService(&LogServie_ServiceDesc, srv)
}

func _LogServie_WriteLog_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LogRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogServieServer).WriteLog(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/logs.LogServie/WriteLog",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogServieServer).WriteLog(ctx, req.(*LogRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// LogServie_ServiceDesc is the grpc.ServiceDesc for LogServie service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LogServie_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "logs.LogServie",
	HandlerType: (*LogServieServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "WriteLog",
			Handler:    _LogServie_WriteLog_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "logs.proto",
}