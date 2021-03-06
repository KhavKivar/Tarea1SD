// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package paquete

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// LogisticaClienteClient is the client API for LogisticaCliente service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LogisticaClienteClient interface {
	// Servicios Logistica-Cliente
	EnviarPedido(ctx context.Context, in *Orden, opts ...grpc.CallOption) (*OrdenRecibida, error)
	SolicitarSeguimiento(ctx context.Context, in *Seguimiento, opts ...grpc.CallOption) (*Estado, error)
	// Servicios Logistica-Camion
	SolicitudPaquetes(ctx context.Context, in *TipoCamion, opts ...grpc.CallOption) (*Paquete, error)
	ResultadoEntrega(ctx context.Context, in *PaqueteRecibido, opts ...grpc.CallOption) (*OrdenRecibida, error)
	ActualizarEstado(ctx context.Context, in *EstadoPaquete, opts ...grpc.CallOption) (*OrdenRecibida, error)
}

type logisticaClienteClient struct {
	cc grpc.ClientConnInterface
}

func NewLogisticaClienteClient(cc grpc.ClientConnInterface) LogisticaClienteClient {
	return &logisticaClienteClient{cc}
}

func (c *logisticaClienteClient) EnviarPedido(ctx context.Context, in *Orden, opts ...grpc.CallOption) (*OrdenRecibida, error) {
	out := new(OrdenRecibida)
	err := c.cc.Invoke(ctx, "/paquete.logistica_cliente/EnviarPedido", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logisticaClienteClient) SolicitarSeguimiento(ctx context.Context, in *Seguimiento, opts ...grpc.CallOption) (*Estado, error) {
	out := new(Estado)
	err := c.cc.Invoke(ctx, "/paquete.logistica_cliente/SolicitarSeguimiento", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logisticaClienteClient) SolicitudPaquetes(ctx context.Context, in *TipoCamion, opts ...grpc.CallOption) (*Paquete, error) {
	out := new(Paquete)
	err := c.cc.Invoke(ctx, "/paquete.logistica_cliente/SolicitudPaquetes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logisticaClienteClient) ResultadoEntrega(ctx context.Context, in *PaqueteRecibido, opts ...grpc.CallOption) (*OrdenRecibida, error) {
	out := new(OrdenRecibida)
	err := c.cc.Invoke(ctx, "/paquete.logistica_cliente/ResultadoEntrega", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logisticaClienteClient) ActualizarEstado(ctx context.Context, in *EstadoPaquete, opts ...grpc.CallOption) (*OrdenRecibida, error) {
	out := new(OrdenRecibida)
	err := c.cc.Invoke(ctx, "/paquete.logistica_cliente/ActualizarEstado", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LogisticaClienteServer is the server API for LogisticaCliente service.
// All implementations must embed UnimplementedLogisticaClienteServer
// for forward compatibility
type LogisticaClienteServer interface {
	// Servicios Logistica-Cliente
	EnviarPedido(context.Context, *Orden) (*OrdenRecibida, error)
	SolicitarSeguimiento(context.Context, *Seguimiento) (*Estado, error)
	// Servicios Logistica-Camion
	SolicitudPaquetes(context.Context, *TipoCamion) (*Paquete, error)
	ResultadoEntrega(context.Context, *PaqueteRecibido) (*OrdenRecibida, error)
	ActualizarEstado(context.Context, *EstadoPaquete) (*OrdenRecibida, error)
	mustEmbedUnimplementedLogisticaClienteServer()
}

// UnimplementedLogisticaClienteServer must be embedded to have forward compatible implementations.
type UnimplementedLogisticaClienteServer struct {
}

func (UnimplementedLogisticaClienteServer) EnviarPedido(context.Context, *Orden) (*OrdenRecibida, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EnviarPedido not implemented")
}
func (UnimplementedLogisticaClienteServer) SolicitarSeguimiento(context.Context, *Seguimiento) (*Estado, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SolicitarSeguimiento not implemented")
}
func (UnimplementedLogisticaClienteServer) SolicitudPaquetes(context.Context, *TipoCamion) (*Paquete, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SolicitudPaquetes not implemented")
}
func (UnimplementedLogisticaClienteServer) ResultadoEntrega(context.Context, *PaqueteRecibido) (*OrdenRecibida, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResultadoEntrega not implemented")
}
func (UnimplementedLogisticaClienteServer) ActualizarEstado(context.Context, *EstadoPaquete) (*OrdenRecibida, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ActualizarEstado not implemented")
}
func (UnimplementedLogisticaClienteServer) mustEmbedUnimplementedLogisticaClienteServer() {}

// UnsafeLogisticaClienteServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LogisticaClienteServer will
// result in compilation errors.
type UnsafeLogisticaClienteServer interface {
	mustEmbedUnimplementedLogisticaClienteServer()
}

func RegisterLogisticaClienteServer(s *grpc.Server, srv LogisticaClienteServer) {
	s.RegisterService(&_LogisticaCliente_serviceDesc, srv)
}

func _LogisticaCliente_EnviarPedido_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Orden)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogisticaClienteServer).EnviarPedido(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/paquete.logistica_cliente/EnviarPedido",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogisticaClienteServer).EnviarPedido(ctx, req.(*Orden))
	}
	return interceptor(ctx, in, info, handler)
}

func _LogisticaCliente_SolicitarSeguimiento_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Seguimiento)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogisticaClienteServer).SolicitarSeguimiento(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/paquete.logistica_cliente/SolicitarSeguimiento",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogisticaClienteServer).SolicitarSeguimiento(ctx, req.(*Seguimiento))
	}
	return interceptor(ctx, in, info, handler)
}

func _LogisticaCliente_SolicitudPaquetes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TipoCamion)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogisticaClienteServer).SolicitudPaquetes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/paquete.logistica_cliente/SolicitudPaquetes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogisticaClienteServer).SolicitudPaquetes(ctx, req.(*TipoCamion))
	}
	return interceptor(ctx, in, info, handler)
}

func _LogisticaCliente_ResultadoEntrega_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PaqueteRecibido)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogisticaClienteServer).ResultadoEntrega(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/paquete.logistica_cliente/ResultadoEntrega",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogisticaClienteServer).ResultadoEntrega(ctx, req.(*PaqueteRecibido))
	}
	return interceptor(ctx, in, info, handler)
}

func _LogisticaCliente_ActualizarEstado_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EstadoPaquete)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogisticaClienteServer).ActualizarEstado(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/paquete.logistica_cliente/ActualizarEstado",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogisticaClienteServer).ActualizarEstado(ctx, req.(*EstadoPaquete))
	}
	return interceptor(ctx, in, info, handler)
}

var _LogisticaCliente_serviceDesc = grpc.ServiceDesc{
	ServiceName: "paquete.logistica_cliente",
	HandlerType: (*LogisticaClienteServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "EnviarPedido",
			Handler:    _LogisticaCliente_EnviarPedido_Handler,
		},
		{
			MethodName: "SolicitarSeguimiento",
			Handler:    _LogisticaCliente_SolicitarSeguimiento_Handler,
		},
		{
			MethodName: "SolicitudPaquetes",
			Handler:    _LogisticaCliente_SolicitudPaquetes_Handler,
		},
		{
			MethodName: "ResultadoEntrega",
			Handler:    _LogisticaCliente_ResultadoEntrega_Handler,
		},
		{
			MethodName: "ActualizarEstado",
			Handler:    _LogisticaCliente_ActualizarEstado_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "paquete/paquete.proto",
}
