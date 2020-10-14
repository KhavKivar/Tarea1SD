package main

import (
	"context"
	"log"
	"net"
	"strings"

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
	timestamp   string
	id          string
	producto    string
	valor       int32
	tienda      string
	destino     string
	prioritario int32
	estado      string
	seguimiento string
}

type paquete struct{
	id string
	seguimiento string
	tipo string
	valor int32
	intentos int32
	estado string
}

var ordenes []orden

var retail []paquete
var normal []paquete
var prioritario []paquete

var num_seguimiento int

func (s *server) EnviarPedido(ctx context.Context, in *pb.Orden) (*pb.OrdenRecibida, error) {
	log.Printf("Pedido Recibido Con id %v desde  %v hacia  %v", in.GetId(), in.GetTienda(), in.GetDestino())

	var orden_nueva orden
	var pack paquete

	orden_nueva.timestamp = ""
	orden_nueva.id = in.GetId()
	orden_nueva.producto = in.GetProducto()
	orden_nueva.valor = in.GetValor()
	orden_nueva.tienda = in.GetTienda()
	orden_nueva.destino = in.GetDestino()
	orden_nueva.prioritario = in.GetPrioritario()
	orden_nueva.estado = "En bodega"
	orden_nueva.seguimiento = "0"

	pack.id = in.GetId()
	pack.seguimiento = "0"
	pack.tipo = ""
	pack.valor = in.GetValor()
	pack.intentos = 0
	pack.estado = "En bodega"


	if in.GetTienda() == "pyme" {
		orden_nueva.seguimiento = strconv.Itoa(num_seguimiento)
		pack.seguimiento = strconv.Itoa(num_seguimiento)
		num_seguimiento = num_seguimiento + 1
		
		if orden_nueva.prioritario == 1 {
			pack.tipo = "prioritario"
			prioritario = append(prioritario, pack)
		} else {
			pack.tipo = "normal"
			normal = append(normal, pack)
		}
	}else {
		pack.tipo = "retail"
		retail = append(retail, pack)
	}

	ordenes = append(ordenes, orden_nueva)
	log.Printf("El numero de seguimiento de la orden es: %v", orden_nueva.seguimiento)

	return &pb.OrdenRecibida{Message: "Orden recibida " + in.GetId()}, nil
}

func (s *server) SolicitarSeguimiento(ctx context.Context, in *pb.Seguimiento) (*pb.Estado, error) {
	i := 0
	log.Printf("Consulta recibida por el numero de seguimiento: %v", in.GetSeguimiento())
	for i < len(ordenes) {
		var x = strings.TrimSuffix(ordenes[i].seguimiento, "\n")
		var y = strings.TrimSuffix(in.GetSeguimiento(), "\n")
		if x == y && in.GetSeguimiento() != "0" {
			return &pb.Estado{Estado: "El estado de la orden es " + ordenes[i].estado}, nil
		}
		i++
	}
	return &pb.Estado{Estado: "La orden no existe"}, nil
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
