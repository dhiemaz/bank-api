// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: rpc_bank.proto

package pb

import (
	context "context"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// BankServiceClient is the client API for BankService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BankServiceClient interface {
	// Auth gRPC calls
	Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error)
	Logout(ctx context.Context, in *LogoutRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	// User gRPC calls
	CreateUser(ctx context.Context, in *UserRequest, opts ...grpc.CallOption) (*UserResponse, error)
	GetUser(ctx context.Context, in *Username, opts ...grpc.CallOption) (*UserResponse, error)
	UpdateUser(ctx context.Context, in *UserUpdateRequest, opts ...grpc.CallOption) (*UserResponse, error)
	DeleteUser(ctx context.Context, in *Username, opts ...grpc.CallOption) (*empty.Empty, error)
}

type bankServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewBankServiceClient(cc grpc.ClientConnInterface) BankServiceClient {
	return &bankServiceClient{cc}
}

func (c *bankServiceClient) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
	out := new(LoginResponse)
	err := c.cc.Invoke(ctx, "/pb.BankService/Login", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bankServiceClient) Logout(ctx context.Context, in *LogoutRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/pb.BankService/Logout", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bankServiceClient) CreateUser(ctx context.Context, in *UserRequest, opts ...grpc.CallOption) (*UserResponse, error) {
	out := new(UserResponse)
	err := c.cc.Invoke(ctx, "/pb.BankService/CreateUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bankServiceClient) GetUser(ctx context.Context, in *Username, opts ...grpc.CallOption) (*UserResponse, error) {
	out := new(UserResponse)
	err := c.cc.Invoke(ctx, "/pb.BankService/GetUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bankServiceClient) UpdateUser(ctx context.Context, in *UserUpdateRequest, opts ...grpc.CallOption) (*UserResponse, error) {
	out := new(UserResponse)
	err := c.cc.Invoke(ctx, "/pb.BankService/UpdateUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bankServiceClient) DeleteUser(ctx context.Context, in *Username, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/pb.BankService/DeleteUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BankServiceServer is the server API for BankService service.
// All implementations must embed UnimplementedBankServiceServer
// for forward compatibility
type BankServiceServer interface {
	// Auth gRPC calls
	Login(context.Context, *LoginRequest) (*LoginResponse, error)
	Logout(context.Context, *LogoutRequest) (*empty.Empty, error)
	// User gRPC calls
	CreateUser(context.Context, *UserRequest) (*UserResponse, error)
	GetUser(context.Context, *Username) (*UserResponse, error)
	UpdateUser(context.Context, *UserUpdateRequest) (*UserResponse, error)
	DeleteUser(context.Context, *Username) (*empty.Empty, error)
	mustEmbedUnimplementedBankServiceServer()
}

// UnimplementedBankServiceServer must be embedded to have forward compatible implementations.
type UnimplementedBankServiceServer struct {
}

func (UnimplementedBankServiceServer) Login(context.Context, *LoginRequest) (*LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedBankServiceServer) Logout(context.Context, *LogoutRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Logout not implemented")
}
func (UnimplementedBankServiceServer) CreateUser(context.Context, *UserRequest) (*UserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUser not implemented")
}
func (UnimplementedBankServiceServer) GetUser(context.Context, *Username) (*UserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUser not implemented")
}
func (UnimplementedBankServiceServer) UpdateUser(context.Context, *UserUpdateRequest) (*UserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUser not implemented")
}
func (UnimplementedBankServiceServer) DeleteUser(context.Context, *Username) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteUser not implemented")
}
func (UnimplementedBankServiceServer) mustEmbedUnimplementedBankServiceServer() {}

// UnsafeBankServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BankServiceServer will
// result in compilation errors.
type UnsafeBankServiceServer interface {
	mustEmbedUnimplementedBankServiceServer()
}

func RegisterBankServiceServer(s grpc.ServiceRegistrar, srv BankServiceServer) {
	s.RegisterService(&BankService_ServiceDesc, srv)
}

func _BankService_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BankServiceServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.BankService/Login",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BankServiceServer).Login(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BankService_Logout_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LogoutRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BankServiceServer).Logout(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.BankService/Logout",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BankServiceServer).Logout(ctx, req.(*LogoutRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BankService_CreateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BankServiceServer).CreateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.BankService/CreateUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BankServiceServer).CreateUser(ctx, req.(*UserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BankService_GetUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Username)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BankServiceServer).GetUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.BankService/GetUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BankServiceServer).GetUser(ctx, req.(*Username))
	}
	return interceptor(ctx, in, info, handler)
}

func _BankService_UpdateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserUpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BankServiceServer).UpdateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.BankService/UpdateUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BankServiceServer).UpdateUser(ctx, req.(*UserUpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BankService_DeleteUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Username)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BankServiceServer).DeleteUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.BankService/DeleteUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BankServiceServer).DeleteUser(ctx, req.(*Username))
	}
	return interceptor(ctx, in, info, handler)
}

// BankService_ServiceDesc is the grpc.ServiceDesc for BankService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BankService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.BankService",
	HandlerType: (*BankServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Login",
			Handler:    _BankService_Login_Handler,
		},
		{
			MethodName: "Logout",
			Handler:    _BankService_Logout_Handler,
		},
		{
			MethodName: "CreateUser",
			Handler:    _BankService_CreateUser_Handler,
		},
		{
			MethodName: "GetUser",
			Handler:    _BankService_GetUser_Handler,
		},
		{
			MethodName: "UpdateUser",
			Handler:    _BankService_UpdateUser_Handler,
		},
		{
			MethodName: "DeleteUser",
			Handler:    _BankService_DeleteUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "rpc_bank.proto",
}
