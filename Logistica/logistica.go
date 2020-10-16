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

type orden struct {
	id          string
	estado      string
	idCamion    string
	seguimiento string
	intentos    int32
}

type paquete struct {
	id       string
	tipo     string
	valor    int32
	origen   string
	destino  string
	intentos int32
}

var allQueue []orden

var retail []paquete
var normal []paquete
var prioritario []paquete

var numSeguimiento int

func (s *server) EnviarPedido(ctx context.Context, in *pb.Orden) (*pb.OrdenRecibida, error) {
	log.Printf("Pedido Recibido con id %v desde  %v hacia  %v", in.GetId(), in.GetTienda(), in.GetDestino())
	var pack paquete
	var ord orden

	//Se crea la orden
	ord.id = in.GetId()
	ord.estado = "En bodega"
	ord.idCamion = ""
	ord.seguimiento = "0"
	ord.intentos = 0

	//Se crea el paquete, y se añade a la cola correspondiente
	pack.id = in.GetId()
	pack.valor = in.GetValor()
	pack.origen = in.GetTienda()
	pack.destino = in.GetDestino()
	pack.intentos = 0

	if in.GetTienda() == "pyme" {
		ord.seguimiento = strconv.Itoa(numSeguimiento)
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

	allQueue = append(allQueue, ord)

	//Se añade el pedido al archivo pedidos.csv
	f, err := os.OpenFile("pedidos.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	t := time.Now()
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(t.Format("2006-01-02 15:04:05") + "," + in.GetId() + "," + in.GetProducto() + "," + fmt.Sprint(in.GetValor()) + "," + in.GetTienda() + "," + in.GetDestino() + "," + fmt.Sprint(in.GetPrioritario()) + "," + ord.seguimiento + "," + "En bodega" + "\n"); err != nil {
		log.Println(err)
	}

	log.Printf("Se envio Numero de seguimiento")
	return &pb.OrdenRecibida{Message: "Orden recibida " + in.GetId() + ",Tu numero de seguimiento es: " + ord.seguimiento}, nil
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

func (s *server) ActualizarEstado(ctx context.Context, in *pb.EstadoPaquete) (*pb.OrdenRecibida, error) {
	i := 0
	for i < len(allQueue) {
		if allQueue[i].id == in.GetId() {
			var obj = allQueue[i]
			obj.estado = in.GetEstado()
			allQueue[i] = obj
			return &pb.OrdenRecibida{Message: "Campo actualizado"}, nil
		}
		i++
	}
	return &pb.OrdenRecibida{Message: "No se encontro la id"}, nil
}

func (s *server) ResultadoEntrega(ctx context.Context, in *pb.PaqueteRecibido) (*pb.OrdenRecibida, error) {
	//Actualizar Pedido
	i := 0
	for i < len(allQueue) {
		if allQueue[i].id == in.GetId() {
			var obj = allQueue[i]
			obj.estado = in.GetEstado()
			allQueue[i] = obj
		}
		i++
	}

	log.Printf("Pedido id: %v intentos: %v Estado: %v\n ", in.GetId(), in.GetIntentos(), in.GetEstado())
	return &pb.OrdenRecibida{Message: "Recibido"}, nil
}

func (s *server) SolicitudPaquetes(ctx context.Context, in *pb.TipoCamion) (*pb.Paquete, error) {
	var y string = strings.TrimSuffix(in.Tipo, "\n")
	if y == "retails" {
		//Ver si hay paquetes en retails o prioritario
		if len(retail) > 0 {
			var aux = retail[0]
			//Dequeue
			retail = retail[1:]
			log.Printf("Paquete Enviado id: %v, origen: %v, destino: %v a Camion tipo: %v ", aux.id, aux.origen, aux.destino, in.GetTipo())
			return &pb.Paquete{Id: aux.id, Tipo: aux.tipo, Valor: aux.valor, Origen: aux.origen, Destino: aux.destino, Intentos: aux.intentos}, nil
		}
		if len(prioritario) > 0 {
			var aux = prioritario[0]
			//Dequeue
			prioritario = prioritario[1:]
			log.Printf("Paquete Enviado id: %v, origen: %v, destino: %v a Camion tipo: %v", aux.id, aux.origen, aux.destino, in.GetTipo())
			return &pb.Paquete{Id: aux.id, Tipo: aux.tipo, Valor: aux.valor, Origen: aux.origen, Destino: aux.destino, Intentos: aux.intentos}, nil
		}
	}
	if y == "normal" {
		//Ver si hay paquetes en prioritario o normal
		if len(prioritario) > 0 {
			var aux = prioritario[0]
			//Dequeue
			prioritario = prioritario[1:]
			log.Printf("Paquete Enviado id: %v, origen: %v, destino: %v a Camion tipo: %v", aux.id, aux.origen, aux.destino, in.GetTipo())
			return &pb.Paquete{Id: aux.id, Tipo: aux.tipo, Valor: aux.valor, Origen: aux.origen, Destino: aux.destino, Intentos: aux.intentos}, nil
		}
		if len(normal) > 0 {
			var aux = normal[0]
			//Dequeue
			normal = normal[1:]
			log.Printf("Paquete Enviado id: %v, origen: %v, destino: %v a Camion tipo: %v", aux.id, aux.origen, aux.destino, in.GetTipo())
			return &pb.Paquete{Id: aux.id, Tipo: aux.tipo, Valor: aux.valor, Origen: aux.origen, Destino: aux.destino, Intentos: aux.intentos}, nil
		}
	}
	return &pb.Paquete{Id: "null"}, nil
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
	if _, err := f.WriteString("timestamp,id,producto,valor,tienda,destino,prioritario,seguimiento,estado\n"); err != nil {
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
