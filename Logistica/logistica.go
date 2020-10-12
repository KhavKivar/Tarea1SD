package main

import (
	"context"
	"log"
	"net"

	pb "google.golang.org/TAREA1SD/Logistica/paquete"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type server struct {
	pb.UnimplementedLogisticaClienteServer
}

func (s *server) EnviarPedido(ctx context.Context, in *pb.Orden) (*pb.OrdenRecibida, error) {
	log.Printf("Pedido Recibido Con id %v desde  %v hacia  %v", in.GetId(), in.GetTienda(), in.GetDestino())

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
