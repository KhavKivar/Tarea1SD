package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"encoding/json"

	"github.com/streadway/amqp"

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
	tipo        string
	valor       int32
}

type paquete struct {
	id       string
	tipo     string
	valor    int32
	origen   string
	destino  string
	intentos int32
}

// Struct que se envia a finanzas
type finanzas struct {
	ID          string `json:"id"`
	Seguimiento string `json:"seguimiento"`
	Tipo        string `json:"tipo"`
	Valor       int32  `json:"valor"`
	Intentos    int32  `json:"intentos"`
	Estado      string `json:"estado"`
}

// Funcion que imprime errores de RabbitMQ
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

type safeStruct struct {
	retail []paquete
	mux    sync.Mutex
}

func (c *safeStruct) getSafe(indice int) paquete {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.retail[indice]
}

func (c *safeStruct) updateSafe(p1 paquete, indice int) {
	c.mux.Lock()
	c.retail[indice] = p1
	c.mux.Unlock()
}

func (c *safeStruct) appendSafe(p1 paquete) {
	c.mux.Lock()
	c.retail = append(c.retail, p1)
	c.mux.Unlock()
}

func (c *safeStruct) lenSafe() int {
	c.mux.Lock()
	defer c.mux.Unlock()
	return len(c.retail)
}

func (c *safeStruct) dequeueSafe() paquete {
	c.mux.Lock()
	var p1 paquete
	p1.id = "null"
	if len(c.retail) > 0 {
		p1 = c.retail[0]
		if len(c.retail) == 1 {
			c.retail = make([]paquete, 0)
		} else {
			c.retail = c.retail[1:]
		}

	}
	defer c.mux.Unlock()
	return p1
}

//EnviarAFinanzas ...
func EnviarAFinanzas(pack finanzas) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Falla en conectar a RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Falla en abrir un canal")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Falla en establecer una queue")

	body, err := json.Marshal(pack)
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Falla en enviar el mensaje")
}

var allQueue []orden

var numSeguimiento int

var mutexRetail safeStruct
var mutexNormal safeStruct
var mutexPrioritario safeStruct

//Funcion que recibe los pedidos de los clientes
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
	ord.valor = in.GetValor()

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
			ord.tipo = "prioritario"
			mutexPrioritario.appendSafe(pack)
		} else {
			pack.tipo = "normal"
			ord.tipo = "normal"
			mutexNormal.appendSafe(pack)
		}
	} else {
		pack.tipo = "retail"
		ord.tipo = "retail"
		mutexRetail.appendSafe(pack)
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
	return &pb.OrdenRecibida{Message: ord.seguimiento}, nil
}

// Funcion que responde el estado de un paquete de acuerdo a la consulta de un numero de seguimiento por parte de un cliente
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

//Funcion que actualiza los estos de los pedidos
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

//Funcion que recibe la informacion si el paquete fue entregado con exito o no
func (s *server) ResultadoEntrega(ctx context.Context, in *pb.PaqueteRecibido) (*pb.OrdenRecibida, error) {

	i := 0
	var nuevo finanzas
	for i < len(allQueue) {
		if allQueue[i].id == in.GetId() {
			var obj = allQueue[i]
			obj.estado = in.GetEstado()
			obj.intentos = in.GetIntentos()
			allQueue[i] = obj

			nuevo.ID = allQueue[i].id
			nuevo.Seguimiento = allQueue[i].seguimiento
			nuevo.Tipo = allQueue[i].tipo
			nuevo.Valor = allQueue[i].valor
			nuevo.Intentos = allQueue[i].intentos
			nuevo.Estado = allQueue[i].estado
			EnviarAFinanzas(nuevo)
		}
		i++
	}

	log.Printf("Orden finaliza proceso de entrega id: %v intentos: %v Estado: %v ", in.GetId(), in.GetIntentos(), in.GetEstado())
	return &pb.OrdenRecibida{Message: "Recibido"}, nil
}

//Funcion que recibe un solicitud para entregar un paquete por parte de camiones y le entrega el paquete correspondiente
func (s *server) SolicitudPaquetes(ctx context.Context, in *pb.TipoCamion) (*pb.Paquete, error) {
	var y string = strings.TrimSuffix(in.Tipo, "\n")
	if y == "retails" {
		//Ver si hay paquetes en retails o prioritario
		var aux = mutexRetail.dequeueSafe()
		if aux.id != "null" {
			log.Printf("Orden enviada a camion id: %v, origen: %v, destino: %v a Camion tipo: %v ", aux.id, aux.origen, aux.destino, in.GetTipo())
			return &pb.Paquete{Id: aux.id, Tipo: aux.tipo, Valor: aux.valor, Origen: aux.origen, Destino: aux.destino, Intentos: aux.intentos}, nil
		}
		aux = mutexPrioritario.dequeueSafe()
		if aux.id != "null" {
			log.Printf("Orden enviada a camion id: %v, origen: %v, destino: %v a Camion tipo: %v", aux.id, aux.origen, aux.destino, in.GetTipo())
			return &pb.Paquete{Id: aux.id, Tipo: aux.tipo, Valor: aux.valor, Origen: aux.origen, Destino: aux.destino, Intentos: aux.intentos}, nil
		}

	}
	if y == "normal" {
		//Ver si hay paquetes en prioritario o normal
		var aux = mutexPrioritario.dequeueSafe()
		if aux.id != "null" {
			log.Printf("Orden enviada a camion id: %v, origen: %v, destino: %v a Camion tipo: %v", aux.id, aux.origen, aux.destino, in.GetTipo())
			return &pb.Paquete{Id: aux.id, Tipo: aux.tipo, Valor: aux.valor, Origen: aux.origen, Destino: aux.destino, Intentos: aux.intentos}, nil
		}
		aux = mutexNormal.dequeueSafe()
		if aux.id != "null" {
			log.Printf("Orden enviada a camion id: %v, origen: %v, destino: %v a Camion tipo: %v", aux.id, aux.origen, aux.destino, in.GetTipo())
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
	log.Println("Server corriendo en puerto 50051 y 50052")

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
