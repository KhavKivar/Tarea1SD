package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

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

type paquete struct {
	id          string
	seguimiento string
	tipo        string
	valor       int32
	intentos    int32
	estado      string
}

var ordenes []orden

var retail []paquete
var normal []paquete
var prioritario []paquete

var numSeguimiento int

func (s *server) EnviarPedido(ctx context.Context, in *pb.Orden) (*pb.OrdenRecibida, error) {
	log.Printf("Pedido Recibido Con id %v desde  %v hacia  %v", in.GetId(), in.GetTienda(), in.GetDestino())

	var ordenNueva orden
	var pack paquete
	t := time.Now()
	ordenNueva.timestamp = t.Format("2006-01-02 15:04:05")
	ordenNueva.id = in.GetId()
	ordenNueva.producto = in.GetProducto()
	ordenNueva.valor = in.GetValor()
	ordenNueva.tienda = in.GetTienda()
	ordenNueva.destino = in.GetDestino()
	ordenNueva.prioritario = in.GetPrioritario()
	ordenNueva.estado = "En bodega"
	ordenNueva.seguimiento = "0"

	pack.id = in.GetId()
	pack.seguimiento = "0"
	pack.tipo = ""
	pack.valor = in.GetValor()
	pack.intentos = 0
	pack.estado = "En bodega"

	if in.GetTienda() == "pyme" {
		ordenNueva.seguimiento = strconv.Itoa(numSeguimiento)
		pack.seguimiento = strconv.Itoa(numSeguimiento)
		numSeguimiento = numSeguimiento + 1

		if ordenNueva.prioritario == 1 {
			pack.tipo = "prioritario"
			prioritario = append(prioritario, pack)
		} else {
			pack.tipo = "normal"
			normal = append(normal, pack)
		}
	} else {
		pack.tipo = "retail"
		retail = append(retail, pack)
	}

	//Se añade el pedido al archivo pedidos.csv
	f, err := os.OpenFile("pedidos.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(ordenNueva.timestamp + "," + ordenNueva.id + "," + ordenNueva.producto + "," + fmt.Sprint(ordenNueva.valor) + "," + ordenNueva.tienda + "," + ordenNueva.destino + "," + fmt.Sprint(ordenNueva.prioritario) + "," + ordenNueva.seguimiento + "\n"); err != nil {
		log.Println(err)
	}

	ordenes = append(ordenes, ordenNueva)
	log.Printf("El numero de seguimiento de la orden es: %v", ordenNueva.seguimiento)
	return &pb.OrdenRecibida{Message: "Orden recibida " + in.GetId() + " " + ", Tu numero de seguimiento es:" + ordenNueva.seguimiento}, nil
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

	numSeguimiento = 11111

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
