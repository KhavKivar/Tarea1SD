package main

import (
	"context"
	"log"
	"net"
	//"fmt"
	"strconv"

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
	 timestamp string
	 id string
	 producto string
	 valor int32
	 tienda string
	 destino string
	 prioritario int32
	 estado string
	 seguimiento string
}

var ordenes []orden

var retail []orden
var pyme []orden

var num_seguimiento int


func (s *server) EnviarPedido(ctx context.Context, in *pb.Orden) (*pb.OrdenRecibida, error) {
	log.Printf("Pedido Recibido Con id %v desde  %v hacia  %v", in.GetId(), in.GetTienda(), in.GetDestino())
	
	var orden_nueva orden

	orden_nueva.timestamp = ""
	orden_nueva.id = in.GetId()
	orden_nueva.producto = in.GetProducto()
	orden_nueva.valor = in.GetValor()
	orden_nueva.tienda = in.GetTienda()
	orden_nueva.destino = in.GetDestino()
	orden_nueva.prioritario = in.GetPrioritario()
	orden_nueva.estado = "En bodega"
	
	if in.GetTienda() == "pyme" {
		orden_nueva.seguimiento = strconv.Itoa(num_seguimiento)
		num_seguimiento = num_seguimiento + 1
	}else {
		orden_nueva.seguimiento = "0"
	}

	ordenes = append(ordenes,orden_nueva)
	log.Printf("El numero de seguimiento de la orden es: %v", orden_nueva.seguimiento)

	return &pb.OrdenRecibida{Message: "Orden recibida " + in.GetId()}, nil
}

func main() {
	
	num_seguimiento = 11111

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
