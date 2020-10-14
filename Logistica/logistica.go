package main

import (
	"context"
	"log"
	"net"
	//"fmt"

	pb "google.golang.org/TAREA1SD/Logistica/paquete"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type server struct {
	pb.UnimplementedLogisticaClienteServer
}

type orden struct {
	 id string
	 producto string
	 valor int32
	 tienda string
	 destino string
	 prioritario int32
}

var ordenes []orden

func (s *server) EnviarPedido(ctx context.Context, in *pb.Orden) (*pb.OrdenRecibida, error) {
	log.Printf("Pedido Recibido Con id %v desde  %v hacia  %v", in.GetId(), in.GetTienda(), in.GetDestino())
	
	var orden_nueva orden

	orden_nueva.id = in.GetId()
	orden_nueva.producto = in.GetProducto()
	orden_nueva.valor = in.GetValor()
	orden_nueva.tienda = in.GetTienda()
	orden_nueva.destino = in.GetDestino()
	orden_nueva.prioritario = in.GetPrioritario()

	ordenes = append(ordenes,orden_nueva)

	return &pb.OrdenRecibida{Message: "Orden recibida " + in.GetId()}, nil
}

func main() {
	
	lis, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterLogisticaClienteServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
