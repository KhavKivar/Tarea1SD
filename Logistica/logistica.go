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

type server struct {
	pb.UnimplementedLogisticaClienteServer
}

type paquete struct {
	id          string
	seguimiento string
	tipo        string
	valor       int32
	intentos    int32
	estado      string
}

var allQueue []paquete
var retail []paquete
var normal []paquete
var prioritario []paquete

var numSeguimiento int

func (s *server) EnviarPedido(ctx context.Context, in *pb.Orden) (*pb.OrdenRecibida, error) {
	log.Printf("Pedido Recibido con id %v desde  %v hacia  %v", in.GetId(), in.GetTienda(), in.GetDestino())

	//Se crea el paquete, y se añade a la cola correspondiente
	var pack paquete
	pack.id = in.GetId()
	pack.valor = in.GetValor()
	pack.tipo = "0"
	pack.seguimiento = "0"
	pack.intentos = 0
	pack.estado = "En bodega"

	if in.GetTienda() == "pyme" {
		pack.seguimiento = strconv.Itoa(numSeguimiento)
		numSeguimiento = numSeguimiento + 1
		if in.GetPrioritario() == 1 {
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
	allQueue = append(allQueue, pack)

	//Se añade el pedido al archivo pedidos.csv
	f, err := os.OpenFile("pedidos.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	t := time.Now()
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(t.Format("2006-01-02 15:04:05") + "," + in.GetId() + "," + in.GetProducto() + "," + fmt.Sprint(in.GetValor()) + "," + in.GetTienda() + "," + in.GetDestino() + "," + fmt.Sprint(in.GetPrioritario()) + "," + pack.seguimiento + "," + "En bodega" + "\n"); err != nil {
		log.Println(err)
	}

	log.Printf("Se envio Numero de seguimiento")
	return &pb.OrdenRecibida{Message: "Orden recibida " + in.GetId() + " " + ",Tu numero de seguimiento es:" + pack.seguimiento}, nil
}

func (s *server) SolicitarSeguimiento(ctx context.Context, in *pb.Seguimiento) (*pb.Estado, error) {
	i := 0
	log.Printf("Consulta recibida por el numero de seguimiento: %v", in.GetSeguimiento())
	for i < len(allQueue) {
		var x = strings.TrimSuffix(allQueue[i].seguimiento, "\n")
		var y = strings.TrimSuffix(in.GetSeguimiento(), "\n")
		if x == y && in.GetSeguimiento() != "0" {
			return &pb.Estado{Estado: "El estado de la orden es " + allQueue[i].estado}, nil
		}
		i++
	}
	return &pb.Estado{Estado: "La orden no existe"}, nil
}

func (s *server) SolicitudPaquetes(ctx context.Context, in *pb.TipoCamion) (*pb.Paquete, error) {
	log.Printf("Camion %v solicita pedido\n", in.Tipo)
	log.Printf("Enviando paquete")
	return &pb.Paquete{Id: "hola", Seguimiento: "hola", Tipo: "hola", Valor: 1, Intentos: 1, Estado: "hola"}, nil
}

const (
	port  = ":50051"
	port2 = ":50052"
)

func runServer(l net.Listener) {
	s := grpc.NewServer()
	pb.RegisterLogisticaClienteServer(s, &server{})
	if err := s.Serve(l); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}

func main() {
	numSeguimiento = 11111
	//Eliminar pedidos.csv
	var err = os.Remove("pedidos.csv")
	if err != nil {
		return
	}
	f, err := os.OpenFile("pedidos.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	if _, err := f.WriteString("timestamp,id,producto,valor,tienda,destino,prioritario,seguimiento\n"); err != nil {
		log.Println(err)
	}

	//List port 50051 y 50052
	lis, err := net.Listen("tcp", port)
	lisCamion, err2 := net.Listen("tcp", port2)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	if err2 != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	//Servers

	go runServer(lis)
	runServer(lisCamion)
}
