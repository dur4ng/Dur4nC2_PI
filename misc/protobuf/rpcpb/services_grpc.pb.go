// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.22.1
// source: rpcpb/services.proto

package rpcpb

import (
	"Dur4nC2/misc/protobuf/clientpb"
	"Dur4nC2/misc/protobuf/commonpb"
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// TeamServerRPCClient is the client API for TeamServerRPC service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TeamServerRPCClient interface {
	// *** Beacons ***
	ListBeacons(ctx context.Context, in *commonpb.Empty, opts ...grpc.CallOption) (*clientpb.Beacons, error)
	GetBeacon(ctx context.Context, in *clientpb.Beacon, opts ...grpc.CallOption) (*clientpb.Beacon, error)
	DeleteBeacon(ctx context.Context, in *clientpb.Beacon, opts ...grpc.CallOption) (*commonpb.Empty, error)
	GetBeaconTasks(ctx context.Context, in *clientpb.Beacon, opts ...grpc.CallOption) (*clientpb.BeaconTasks, error)
	GetBeaconTaskContent(ctx context.Context, in *clientpb.BeaconTask, opts ...grpc.CallOption) (*clientpb.BeaconTask, error)
	CancelBeaconTask(ctx context.Context, in *clientpb.BeaconTask, opts ...grpc.CallOption) (*clientpb.BeaconTask, error)
	// *** Operators ***
	ListOperators(ctx context.Context, in *commonpb.Empty, opts ...grpc.CallOption) (*clientpb.Operators, error)
	// *** Hosts ***
	ListHosts(ctx context.Context, in *commonpb.Empty, opts ...grpc.CallOption) (*clientpb.Hosts, error)
	GetHost(ctx context.Context, in *clientpb.Host, opts ...grpc.CallOption) (*clientpb.Host, error)
	DeleteHost(ctx context.Context, in *clientpb.Host, opts ...grpc.CallOption) (*commonpb.Empty, error)
	AddHostIOC(ctx context.Context, in *clientpb.IOC, opts ...grpc.CallOption) (*commonpb.Empty, error)
	DeleteHostIOC(ctx context.Context, in *clientpb.IOC, opts ...grpc.CallOption) (*commonpb.Empty, error)
	DeleteHostLoot(ctx context.Context, in *clientpb.Loot, opts ...grpc.CallOption) (*commonpb.Empty, error)
	// *** Loots ***
	AddLoot(ctx context.Context, in *clientpb.Loot, opts ...grpc.CallOption) (*clientpb.Loot, error)
	DeleteLoot(ctx context.Context, in *clientpb.Loot, opts ...grpc.CallOption) (*commonpb.Empty, error)
	UpdateLoot(ctx context.Context, in *clientpb.Loot, opts ...grpc.CallOption) (*clientpb.Loot, error)
	GetLootContent(ctx context.Context, in *clientpb.Loot, opts ...grpc.CallOption) (*clientpb.Loot, error)
	ListLoot(ctx context.Context, in *commonpb.Empty, opts ...grpc.CallOption) (*clientpb.Loots, error)
	GetLoot(ctx context.Context, in *clientpb.Loot, opts ...grpc.CallOption) (*clientpb.Loots, error)
	// *** Listeners ***
	StartHTTPSListener(ctx context.Context, in *clientpb.HTTPListenerReq, opts ...grpc.CallOption) (*clientpb.HTTPListener, error)
	StartHTTPListener(ctx context.Context, in *clientpb.HTTPListenerReq, opts ...grpc.CallOption) (*clientpb.HTTPListener, error)
	// *** Stager Listener ***
	StartTCPStagerListener(ctx context.Context, in *clientpb.StagerListenerReq, opts ...grpc.CallOption) (*clientpb.StagerListener, error)
	StartHTTPStagerListener(ctx context.Context, in *clientpb.StagerListenerReq, opts ...grpc.CallOption) (*clientpb.StagerListener, error)
}

type teamServerRPCClient struct {
	cc grpc.ClientConnInterface
}

func NewTeamServerRPCClient(cc grpc.ClientConnInterface) TeamServerRPCClient {
	return &teamServerRPCClient{cc}
}

func (c *teamServerRPCClient) ListBeacons(ctx context.Context, in *commonpb.Empty, opts ...grpc.CallOption) (*clientpb.Beacons, error) {
	out := new(clientpb.Beacons)
	err := c.cc.Invoke(ctx, "/rpcpb.TeamServerRPC/ListBeacons", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *teamServerRPCClient) GetBeacon(ctx context.Context, in *clientpb.Beacon, opts ...grpc.CallOption) (*clientpb.Beacon, error) {
	out := new(clientpb.Beacon)
	err := c.cc.Invoke(ctx, "/rpcpb.TeamServerRPC/GetBeacon", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *teamServerRPCClient) DeleteBeacon(ctx context.Context, in *clientpb.Beacon, opts ...grpc.CallOption) (*commonpb.Empty, error) {
	out := new(commonpb.Empty)
	err := c.cc.Invoke(ctx, "/rpcpb.TeamServerRPC/DeleteBeacon", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *teamServerRPCClient) GetBeaconTasks(ctx context.Context, in *clientpb.Beacon, opts ...grpc.CallOption) (*clientpb.BeaconTasks, error) {
	out := new(clientpb.BeaconTasks)
	err := c.cc.Invoke(ctx, "/rpcpb.TeamServerRPC/GetBeaconTasks", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *teamServerRPCClient) GetBeaconTaskContent(ctx context.Context, in *clientpb.BeaconTask, opts ...grpc.CallOption) (*clientpb.BeaconTask, error) {
	out := new(clientpb.BeaconTask)
	err := c.cc.Invoke(ctx, "/rpcpb.TeamServerRPC/GetBeaconTaskContent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *teamServerRPCClient) CancelBeaconTask(ctx context.Context, in *clientpb.BeaconTask, opts ...grpc.CallOption) (*clientpb.BeaconTask, error) {
	out := new(clientpb.BeaconTask)
	err := c.cc.Invoke(ctx, "/rpcpb.TeamServerRPC/CancelBeaconTask", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *teamServerRPCClient) ListOperators(ctx context.Context, in *commonpb.Empty, opts ...grpc.CallOption) (*clientpb.Operators, error) {
	out := new(clientpb.Operators)
	err := c.cc.Invoke(ctx, "/rpcpb.TeamServerRPC/ListOperators", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *teamServerRPCClient) ListHosts(ctx context.Context, in *commonpb.Empty, opts ...grpc.CallOption) (*clientpb.Hosts, error) {
	out := new(clientpb.Hosts)
	err := c.cc.Invoke(ctx, "/rpcpb.TeamServerRPC/ListHosts", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *teamServerRPCClient) GetHost(ctx context.Context, in *clientpb.Host, opts ...grpc.CallOption) (*clientpb.Host, error) {
	out := new(clientpb.Host)
	err := c.cc.Invoke(ctx, "/rpcpb.TeamServerRPC/GetHost", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *teamServerRPCClient) DeleteHost(ctx context.Context, in *clientpb.Host, opts ...grpc.CallOption) (*commonpb.Empty, error) {
	out := new(commonpb.Empty)
	err := c.cc.Invoke(ctx, "/rpcpb.TeamServerRPC/DeleteHost", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *teamServerRPCClient) AddHostIOC(ctx context.Context, in *clientpb.IOC, opts ...grpc.CallOption) (*commonpb.Empty, error) {
	out := new(commonpb.Empty)
	err := c.cc.Invoke(ctx, "/rpcpb.TeamServerRPC/AddHostIOC", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *teamServerRPCClient) DeleteHostIOC(ctx context.Context, in *clientpb.IOC, opts ...grpc.CallOption) (*commonpb.Empty, error) {
	out := new(commonpb.Empty)
	err := c.cc.Invoke(ctx, "/rpcpb.TeamServerRPC/DeleteHostIOC", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *teamServerRPCClient) DeleteHostLoot(ctx context.Context, in *clientpb.Loot, opts ...grpc.CallOption) (*commonpb.Empty, error) {
	out := new(commonpb.Empty)
	err := c.cc.Invoke(ctx, "/rpcpb.TeamServerRPC/DeleteHostLoot", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *teamServerRPCClient) AddLoot(ctx context.Context, in *clientpb.Loot, opts ...grpc.CallOption) (*clientpb.Loot, error) {
	out := new(clientpb.Loot)
	err := c.cc.Invoke(ctx, "/rpcpb.TeamServerRPC/AddLoot", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *teamServerRPCClient) DeleteLoot(ctx context.Context, in *clientpb.Loot, opts ...grpc.CallOption) (*commonpb.Empty, error) {
	out := new(commonpb.Empty)
	err := c.cc.Invoke(ctx, "/rpcpb.TeamServerRPC/DeleteLoot", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *teamServerRPCClient) UpdateLoot(ctx context.Context, in *clientpb.Loot, opts ...grpc.CallOption) (*clientpb.Loot, error) {
	out := new(clientpb.Loot)
	err := c.cc.Invoke(ctx, "/rpcpb.TeamServerRPC/UpdateLoot", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *teamServerRPCClient) GetLootContent(ctx context.Context, in *clientpb.Loot, opts ...grpc.CallOption) (*clientpb.Loot, error) {
	out := new(clientpb.Loot)
	err := c.cc.Invoke(ctx, "/rpcpb.TeamServerRPC/GetLootContent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *teamServerRPCClient) ListLoot(ctx context.Context, in *commonpb.Empty, opts ...grpc.CallOption) (*clientpb.Loots, error) {
	out := new(clientpb.Loots)
	err := c.cc.Invoke(ctx, "/rpcpb.TeamServerRPC/ListLoot", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *teamServerRPCClient) GetLoot(ctx context.Context, in *clientpb.Loot, opts ...grpc.CallOption) (*clientpb.Loots, error) {
	out := new(clientpb.Loots)
	err := c.cc.Invoke(ctx, "/rpcpb.TeamServerRPC/GetLoot", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *teamServerRPCClient) StartHTTPSListener(ctx context.Context, in *clientpb.HTTPListenerReq, opts ...grpc.CallOption) (*clientpb.HTTPListener, error) {
	out := new(clientpb.HTTPListener)
	err := c.cc.Invoke(ctx, "/rpcpb.TeamServerRPC/StartHTTPSListener", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *teamServerRPCClient) StartHTTPListener(ctx context.Context, in *clientpb.HTTPListenerReq, opts ...grpc.CallOption) (*clientpb.HTTPListener, error) {
	out := new(clientpb.HTTPListener)
	err := c.cc.Invoke(ctx, "/rpcpb.TeamServerRPC/StartHTTPListener", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *teamServerRPCClient) StartTCPStagerListener(ctx context.Context, in *clientpb.StagerListenerReq, opts ...grpc.CallOption) (*clientpb.StagerListener, error) {
	out := new(clientpb.StagerListener)
	err := c.cc.Invoke(ctx, "/rpcpb.TeamServerRPC/StartTCPStagerListener", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *teamServerRPCClient) StartHTTPStagerListener(ctx context.Context, in *clientpb.StagerListenerReq, opts ...grpc.CallOption) (*clientpb.StagerListener, error) {
	out := new(clientpb.StagerListener)
	err := c.cc.Invoke(ctx, "/rpcpb.TeamServerRPC/StartHTTPStagerListener", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TeamServerRPCServer is the server API for TeamServerRPC service.
// All implementations must embed UnimplementedTeamServerRPCServer
// for forward compatibility
type TeamServerRPCServer interface {
	// *** Beacons ***
	ListBeacons(context.Context, *commonpb.Empty) (*clientpb.Beacons, error)
	GetBeacon(context.Context, *clientpb.Beacon) (*clientpb.Beacon, error)
	DeleteBeacon(context.Context, *clientpb.Beacon) (*commonpb.Empty, error)
	GetBeaconTasks(context.Context, *clientpb.Beacon) (*clientpb.BeaconTasks, error)
	GetBeaconTaskContent(context.Context, *clientpb.BeaconTask) (*clientpb.BeaconTask, error)
	CancelBeaconTask(context.Context, *clientpb.BeaconTask) (*clientpb.BeaconTask, error)
	// *** Operators ***
	ListOperators(context.Context, *commonpb.Empty) (*clientpb.Operators, error)
	// *** Hosts ***
	ListHosts(context.Context, *commonpb.Empty) (*clientpb.Hosts, error)
	GetHost(context.Context, *clientpb.Host) (*clientpb.Host, error)
	DeleteHost(context.Context, *clientpb.Host) (*commonpb.Empty, error)
	AddHostIOC(context.Context, *clientpb.IOC) (*commonpb.Empty, error)
	DeleteHostIOC(context.Context, *clientpb.IOC) (*commonpb.Empty, error)
	DeleteHostLoot(context.Context, *clientpb.Loot) (*commonpb.Empty, error)
	// *** Loots ***
	AddLoot(context.Context, *clientpb.Loot) (*clientpb.Loot, error)
	DeleteLoot(context.Context, *clientpb.Loot) (*commonpb.Empty, error)
	UpdateLoot(context.Context, *clientpb.Loot) (*clientpb.Loot, error)
	GetLootContent(context.Context, *clientpb.Loot) (*clientpb.Loot, error)
	ListLoot(context.Context, *commonpb.Empty) (*clientpb.Loots, error)
	GetLoot(context.Context, *clientpb.Loot) (*clientpb.Loots, error)
	// *** Listeners ***
	StartHTTPSListener(context.Context, *clientpb.HTTPListenerReq) (*clientpb.HTTPListener, error)
	StartHTTPListener(context.Context, *clientpb.HTTPListenerReq) (*clientpb.HTTPListener, error)
	// *** Stager Listener ***
	StartTCPStagerListener(context.Context, *clientpb.StagerListenerReq) (*clientpb.StagerListener, error)
	StartHTTPStagerListener(context.Context, *clientpb.StagerListenerReq) (*clientpb.StagerListener, error)
	mustEmbedUnimplementedTeamServerRPCServer()
}

// UnimplementedTeamServerRPCServer must be embedded to have forward compatible implementations.
type UnimplementedTeamServerRPCServer struct {
}

func (UnimplementedTeamServerRPCServer) ListBeacons(context.Context, *commonpb.Empty) (*clientpb.Beacons, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListBeacons not implemented")
}
func (UnimplementedTeamServerRPCServer) GetBeacon(context.Context, *clientpb.Beacon) (*clientpb.Beacon, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBeacon not implemented")
}
func (UnimplementedTeamServerRPCServer) DeleteBeacon(context.Context, *clientpb.Beacon) (*commonpb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteBeacon not implemented")
}
func (UnimplementedTeamServerRPCServer) GetBeaconTasks(context.Context, *clientpb.Beacon) (*clientpb.BeaconTasks, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBeaconTasks not implemented")
}
func (UnimplementedTeamServerRPCServer) GetBeaconTaskContent(context.Context, *clientpb.BeaconTask) (*clientpb.BeaconTask, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBeaconTaskContent not implemented")
}
func (UnimplementedTeamServerRPCServer) CancelBeaconTask(context.Context, *clientpb.BeaconTask) (*clientpb.BeaconTask, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CancelBeaconTask not implemented")
}
func (UnimplementedTeamServerRPCServer) ListOperators(context.Context, *commonpb.Empty) (*clientpb.Operators, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListOperators not implemented")
}
func (UnimplementedTeamServerRPCServer) ListHosts(context.Context, *commonpb.Empty) (*clientpb.Hosts, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListHosts not implemented")
}
func (UnimplementedTeamServerRPCServer) GetHost(context.Context, *clientpb.Host) (*clientpb.Host, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetHost not implemented")
}
func (UnimplementedTeamServerRPCServer) DeleteHost(context.Context, *clientpb.Host) (*commonpb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteHost not implemented")
}
func (UnimplementedTeamServerRPCServer) AddHostIOC(context.Context, *clientpb.IOC) (*commonpb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddHostIOC not implemented")
}
func (UnimplementedTeamServerRPCServer) DeleteHostIOC(context.Context, *clientpb.IOC) (*commonpb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteHostIOC not implemented")
}
func (UnimplementedTeamServerRPCServer) DeleteHostLoot(context.Context, *clientpb.Loot) (*commonpb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteHostLoot not implemented")
}
func (UnimplementedTeamServerRPCServer) AddLoot(context.Context, *clientpb.Loot) (*clientpb.Loot, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddLoot not implemented")
}
func (UnimplementedTeamServerRPCServer) DeleteLoot(context.Context, *clientpb.Loot) (*commonpb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteLoot not implemented")
}
func (UnimplementedTeamServerRPCServer) UpdateLoot(context.Context, *clientpb.Loot) (*clientpb.Loot, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateLoot not implemented")
}
func (UnimplementedTeamServerRPCServer) GetLootContent(context.Context, *clientpb.Loot) (*clientpb.Loot, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLootContent not implemented")
}
func (UnimplementedTeamServerRPCServer) ListLoot(context.Context, *commonpb.Empty) (*clientpb.Loots, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListLoot not implemented")
}
func (UnimplementedTeamServerRPCServer) GetLoot(context.Context, *clientpb.Loot) (*clientpb.Loots, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLoot not implemented")
}
func (UnimplementedTeamServerRPCServer) StartHTTPSListener(context.Context, *clientpb.HTTPListenerReq) (*clientpb.HTTPListener, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartHTTPSListener not implemented")
}
func (UnimplementedTeamServerRPCServer) StartHTTPListener(context.Context, *clientpb.HTTPListenerReq) (*clientpb.HTTPListener, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartHTTPListener not implemented")
}
func (UnimplementedTeamServerRPCServer) StartTCPStagerListener(context.Context, *clientpb.StagerListenerReq) (*clientpb.StagerListener, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartTCPStagerListener not implemented")
}
func (UnimplementedTeamServerRPCServer) StartHTTPStagerListener(context.Context, *clientpb.StagerListenerReq) (*clientpb.StagerListener, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartHTTPStagerListener not implemented")
}
func (UnimplementedTeamServerRPCServer) mustEmbedUnimplementedTeamServerRPCServer() {}

// UnsafeTeamServerRPCServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TeamServerRPCServer will
// result in compilation errors.
type UnsafeTeamServerRPCServer interface {
	mustEmbedUnimplementedTeamServerRPCServer()
}

func RegisterTeamServerRPCServer(s grpc.ServiceRegistrar, srv TeamServerRPCServer) {
	s.RegisterService(&TeamServerRPC_ServiceDesc, srv)
}

func _TeamServerRPC_ListBeacons_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(commonpb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TeamServerRPCServer).ListBeacons(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpcpb.TeamServerRPC/ListBeacons",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TeamServerRPCServer).ListBeacons(ctx, req.(*commonpb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _TeamServerRPC_GetBeacon_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(clientpb.Beacon)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TeamServerRPCServer).GetBeacon(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpcpb.TeamServerRPC/GetBeacon",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TeamServerRPCServer).GetBeacon(ctx, req.(*clientpb.Beacon))
	}
	return interceptor(ctx, in, info, handler)
}

func _TeamServerRPC_DeleteBeacon_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(clientpb.Beacon)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TeamServerRPCServer).DeleteBeacon(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpcpb.TeamServerRPC/DeleteBeacon",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TeamServerRPCServer).DeleteBeacon(ctx, req.(*clientpb.Beacon))
	}
	return interceptor(ctx, in, info, handler)
}

func _TeamServerRPC_GetBeaconTasks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(clientpb.Beacon)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TeamServerRPCServer).GetBeaconTasks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpcpb.TeamServerRPC/GetBeaconTasks",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TeamServerRPCServer).GetBeaconTasks(ctx, req.(*clientpb.Beacon))
	}
	return interceptor(ctx, in, info, handler)
}

func _TeamServerRPC_GetBeaconTaskContent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(clientpb.BeaconTask)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TeamServerRPCServer).GetBeaconTaskContent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpcpb.TeamServerRPC/GetBeaconTaskContent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TeamServerRPCServer).GetBeaconTaskContent(ctx, req.(*clientpb.BeaconTask))
	}
	return interceptor(ctx, in, info, handler)
}

func _TeamServerRPC_CancelBeaconTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(clientpb.BeaconTask)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TeamServerRPCServer).CancelBeaconTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpcpb.TeamServerRPC/CancelBeaconTask",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TeamServerRPCServer).CancelBeaconTask(ctx, req.(*clientpb.BeaconTask))
	}
	return interceptor(ctx, in, info, handler)
}

func _TeamServerRPC_ListOperators_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(commonpb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TeamServerRPCServer).ListOperators(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpcpb.TeamServerRPC/ListOperators",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TeamServerRPCServer).ListOperators(ctx, req.(*commonpb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _TeamServerRPC_ListHosts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(commonpb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TeamServerRPCServer).ListHosts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpcpb.TeamServerRPC/ListHosts",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TeamServerRPCServer).ListHosts(ctx, req.(*commonpb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _TeamServerRPC_GetHost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(clientpb.Host)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TeamServerRPCServer).GetHost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpcpb.TeamServerRPC/GetHost",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TeamServerRPCServer).GetHost(ctx, req.(*clientpb.Host))
	}
	return interceptor(ctx, in, info, handler)
}

func _TeamServerRPC_DeleteHost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(clientpb.Host)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TeamServerRPCServer).DeleteHost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpcpb.TeamServerRPC/DeleteHost",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TeamServerRPCServer).DeleteHost(ctx, req.(*clientpb.Host))
	}
	return interceptor(ctx, in, info, handler)
}

func _TeamServerRPC_AddHostIOC_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(clientpb.IOC)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TeamServerRPCServer).AddHostIOC(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpcpb.TeamServerRPC/AddHostIOC",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TeamServerRPCServer).AddHostIOC(ctx, req.(*clientpb.IOC))
	}
	return interceptor(ctx, in, info, handler)
}

func _TeamServerRPC_DeleteHostIOC_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(clientpb.IOC)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TeamServerRPCServer).DeleteHostIOC(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpcpb.TeamServerRPC/DeleteHostIOC",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TeamServerRPCServer).DeleteHostIOC(ctx, req.(*clientpb.IOC))
	}
	return interceptor(ctx, in, info, handler)
}

func _TeamServerRPC_DeleteHostLoot_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(clientpb.Loot)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TeamServerRPCServer).DeleteHostLoot(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpcpb.TeamServerRPC/DeleteHostLoot",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TeamServerRPCServer).DeleteHostLoot(ctx, req.(*clientpb.Loot))
	}
	return interceptor(ctx, in, info, handler)
}

func _TeamServerRPC_AddLoot_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(clientpb.Loot)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TeamServerRPCServer).AddLoot(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpcpb.TeamServerRPC/AddLoot",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TeamServerRPCServer).AddLoot(ctx, req.(*clientpb.Loot))
	}
	return interceptor(ctx, in, info, handler)
}

func _TeamServerRPC_DeleteLoot_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(clientpb.Loot)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TeamServerRPCServer).DeleteLoot(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpcpb.TeamServerRPC/DeleteLoot",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TeamServerRPCServer).DeleteLoot(ctx, req.(*clientpb.Loot))
	}
	return interceptor(ctx, in, info, handler)
}

func _TeamServerRPC_UpdateLoot_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(clientpb.Loot)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TeamServerRPCServer).UpdateLoot(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpcpb.TeamServerRPC/UpdateLoot",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TeamServerRPCServer).UpdateLoot(ctx, req.(*clientpb.Loot))
	}
	return interceptor(ctx, in, info, handler)
}

func _TeamServerRPC_GetLootContent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(clientpb.Loot)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TeamServerRPCServer).GetLootContent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpcpb.TeamServerRPC/GetLootContent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TeamServerRPCServer).GetLootContent(ctx, req.(*clientpb.Loot))
	}
	return interceptor(ctx, in, info, handler)
}

func _TeamServerRPC_ListLoot_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(commonpb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TeamServerRPCServer).ListLoot(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpcpb.TeamServerRPC/ListLoot",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TeamServerRPCServer).ListLoot(ctx, req.(*commonpb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _TeamServerRPC_GetLoot_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(clientpb.Loot)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TeamServerRPCServer).GetLoot(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpcpb.TeamServerRPC/GetLoot",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TeamServerRPCServer).GetLoot(ctx, req.(*clientpb.Loot))
	}
	return interceptor(ctx, in, info, handler)
}

func _TeamServerRPC_StartHTTPSListener_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(clientpb.HTTPListenerReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TeamServerRPCServer).StartHTTPSListener(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpcpb.TeamServerRPC/StartHTTPSListener",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TeamServerRPCServer).StartHTTPSListener(ctx, req.(*clientpb.HTTPListenerReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _TeamServerRPC_StartHTTPListener_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(clientpb.HTTPListenerReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TeamServerRPCServer).StartHTTPListener(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpcpb.TeamServerRPC/StartHTTPListener",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TeamServerRPCServer).StartHTTPListener(ctx, req.(*clientpb.HTTPListenerReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _TeamServerRPC_StartTCPStagerListener_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(clientpb.StagerListenerReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TeamServerRPCServer).StartTCPStagerListener(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpcpb.TeamServerRPC/StartTCPStagerListener",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TeamServerRPCServer).StartTCPStagerListener(ctx, req.(*clientpb.StagerListenerReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _TeamServerRPC_StartHTTPStagerListener_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(clientpb.StagerListenerReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TeamServerRPCServer).StartHTTPStagerListener(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpcpb.TeamServerRPC/StartHTTPStagerListener",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TeamServerRPCServer).StartHTTPStagerListener(ctx, req.(*clientpb.StagerListenerReq))
	}
	return interceptor(ctx, in, info, handler)
}

// TeamServerRPC_ServiceDesc is the grpc.ServiceDesc for TeamServerRPC service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TeamServerRPC_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "rpcpb.TeamServerRPC",
	HandlerType: (*TeamServerRPCServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListBeacons",
			Handler:    _TeamServerRPC_ListBeacons_Handler,
		},
		{
			MethodName: "GetBeacon",
			Handler:    _TeamServerRPC_GetBeacon_Handler,
		},
		{
			MethodName: "DeleteBeacon",
			Handler:    _TeamServerRPC_DeleteBeacon_Handler,
		},
		{
			MethodName: "GetBeaconTasks",
			Handler:    _TeamServerRPC_GetBeaconTasks_Handler,
		},
		{
			MethodName: "GetBeaconTaskContent",
			Handler:    _TeamServerRPC_GetBeaconTaskContent_Handler,
		},
		{
			MethodName: "CancelBeaconTask",
			Handler:    _TeamServerRPC_CancelBeaconTask_Handler,
		},
		{
			MethodName: "ListOperators",
			Handler:    _TeamServerRPC_ListOperators_Handler,
		},
		{
			MethodName: "ListHosts",
			Handler:    _TeamServerRPC_ListHosts_Handler,
		},
		{
			MethodName: "GetHost",
			Handler:    _TeamServerRPC_GetHost_Handler,
		},
		{
			MethodName: "DeleteHost",
			Handler:    _TeamServerRPC_DeleteHost_Handler,
		},
		{
			MethodName: "AddHostIOC",
			Handler:    _TeamServerRPC_AddHostIOC_Handler,
		},
		{
			MethodName: "DeleteHostIOC",
			Handler:    _TeamServerRPC_DeleteHostIOC_Handler,
		},
		{
			MethodName: "DeleteHostLoot",
			Handler:    _TeamServerRPC_DeleteHostLoot_Handler,
		},
		{
			MethodName: "AddLoot",
			Handler:    _TeamServerRPC_AddLoot_Handler,
		},
		{
			MethodName: "DeleteLoot",
			Handler:    _TeamServerRPC_DeleteLoot_Handler,
		},
		{
			MethodName: "UpdateLoot",
			Handler:    _TeamServerRPC_UpdateLoot_Handler,
		},
		{
			MethodName: "GetLootContent",
			Handler:    _TeamServerRPC_GetLootContent_Handler,
		},
		{
			MethodName: "ListLoot",
			Handler:    _TeamServerRPC_ListLoot_Handler,
		},
		{
			MethodName: "GetLoot",
			Handler:    _TeamServerRPC_GetLoot_Handler,
		},
		{
			MethodName: "StartHTTPSListener",
			Handler:    _TeamServerRPC_StartHTTPSListener_Handler,
		},
		{
			MethodName: "StartHTTPListener",
			Handler:    _TeamServerRPC_StartHTTPListener_Handler,
		},
		{
			MethodName: "StartTCPStagerListener",
			Handler:    _TeamServerRPC_StartTCPStagerListener_Handler,
		},
		{
			MethodName: "StartHTTPStagerListener",
			Handler:    _TeamServerRPC_StartHTTPStagerListener_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "rpcpb/services.proto",
}
