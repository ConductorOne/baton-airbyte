// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: c1/ratelimit/v1/ratelimit.proto

package v1

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
	RateLimiterService_Do_FullMethodName     = "/c1.ratelimit.v1.RateLimiterService/Do"
	RateLimiterService_Report_FullMethodName = "/c1.ratelimit.v1.RateLimiterService/Report"
)

// RateLimiterServiceClient is the client API for RateLimiterService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RateLimiterServiceClient interface {
	Do(ctx context.Context, in *DoRequest, opts ...grpc.CallOption) (*DoResponse, error)
	Report(ctx context.Context, in *ReportRequest, opts ...grpc.CallOption) (*ReportResponse, error)
}

type rateLimiterServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRateLimiterServiceClient(cc grpc.ClientConnInterface) RateLimiterServiceClient {
	return &rateLimiterServiceClient{cc}
}

func (c *rateLimiterServiceClient) Do(ctx context.Context, in *DoRequest, opts ...grpc.CallOption) (*DoResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DoResponse)
	err := c.cc.Invoke(ctx, RateLimiterService_Do_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rateLimiterServiceClient) Report(ctx context.Context, in *ReportRequest, opts ...grpc.CallOption) (*ReportResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ReportResponse)
	err := c.cc.Invoke(ctx, RateLimiterService_Report_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RateLimiterServiceServer is the server API for RateLimiterService service.
// All implementations should embed UnimplementedRateLimiterServiceServer
// for forward compatibility.
type RateLimiterServiceServer interface {
	Do(context.Context, *DoRequest) (*DoResponse, error)
	Report(context.Context, *ReportRequest) (*ReportResponse, error)
}

// UnimplementedRateLimiterServiceServer should be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedRateLimiterServiceServer struct{}

func (UnimplementedRateLimiterServiceServer) Do(context.Context, *DoRequest) (*DoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Do not implemented")
}
func (UnimplementedRateLimiterServiceServer) Report(context.Context, *ReportRequest) (*ReportResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Report not implemented")
}
func (UnimplementedRateLimiterServiceServer) testEmbeddedByValue() {}

// UnsafeRateLimiterServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RateLimiterServiceServer will
// result in compilation errors.
type UnsafeRateLimiterServiceServer interface {
	mustEmbedUnimplementedRateLimiterServiceServer()
}

func RegisterRateLimiterServiceServer(s grpc.ServiceRegistrar, srv RateLimiterServiceServer) {
	// If the following call pancis, it indicates UnimplementedRateLimiterServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&RateLimiterService_ServiceDesc, srv)
}

func _RateLimiterService_Do_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RateLimiterServiceServer).Do(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RateLimiterService_Do_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RateLimiterServiceServer).Do(ctx, req.(*DoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RateLimiterService_Report_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReportRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RateLimiterServiceServer).Report(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RateLimiterService_Report_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RateLimiterServiceServer).Report(ctx, req.(*ReportRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RateLimiterService_ServiceDesc is the grpc.ServiceDesc for RateLimiterService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RateLimiterService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "c1.ratelimit.v1.RateLimiterService",
	HandlerType: (*RateLimiterServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Do",
			Handler:    _RateLimiterService_Do_Handler,
		},
		{
			MethodName: "Report",
			Handler:    _RateLimiterService_Report_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "c1/ratelimit/v1/ratelimit.proto",
}
