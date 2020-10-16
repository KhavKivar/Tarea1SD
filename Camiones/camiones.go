package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	pb "google.golang.org/TAREA1SD/Logistica/paquete"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50052"
)

var tiempo int
var tiempoEntrega int

type paquete struct {
	id           string
	tipo         string
	valor        int32
	origen       string
	destino      string
	intentos     int32
	fechaEntrega string
	tipoCamion   string
	estado       string
}

var allPedidos []paquete
var listRetail []paquete
var listRetail2 []paquete
var listNormal []paquete

func getPaquete(c pb.LogisticaClienteClient, tipoCamion string, idCamion string) (bool, paquete) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SolicitudPaquetes(ctx, &pb.TipoCamion{Tipo: tipoCamion})
	if err != nil {
		log.Fatalf("Error al obtener el paquete: %v", err)
	}
	var newPaquete paquete
	if r.GetId() != "null" {
		log.Printf("Paquete recibido: id:%v, Origen: %v, Destino:%v, Camion encargado: %v", r.GetId(), r.GetOrigen(), r.GetDestino(), idCamion)
		// Paquete recibido
		newPaquete.id = r.GetId()
		newPaquete.tipo = r.GetTipo()
		newPaquete.valor = r.GetValor()
		newPaquete.origen = r.GetOrigen()
		newPaquete.destino = r.GetDestino()
		newPaquete.intentos = r.GetIntentos()
		newPaquete.estado = "En proceso de entrega"
		newPaquete.fechaEntrega = ""

		if idCamion == "Retails 1" {
			newPaquete.tipoCamion = idCamion
			listRetail = append(listRetail, newPaquete)
		}
		if idCamion == "Retails 2" {
			newPaquete.tipoCamion = idCamion
			listRetail2 = append(listRetail2, newPaquete)
		}
		if idCamion == "normal" {
			newPaquete.tipoCamion = idCamion
			listNormal = append(listNormal, newPaquete)
		}
		allPedidos = append(allPedidos, newPaquete)
		return true, newPaquete
	}
	return false, newPaquete
}

func random80() bool {
	n := rand.Intn(10)
	if n < 8 {
		return true
	}
	return false
}

func clienteRecibe(maxIntentos int) int {
	intentos := 0
	time.Sleep(time.Duration(tiempoEntrega) * time.Millisecond)
	if random80() {
		return intentos
	}
	intentos++
	if intentos == maxIntentos {
		return intentos
	}
	time.Sleep(time.Duration(tiempoEntrega) * time.Millisecond)
	if random80() {
		return intentos
	}
	intentos++
	if intentos == maxIntentos {
		return intentos
	}
	time.Sleep(time.Duration(tiempoEntrega) * time.Millisecond)
	if random80() {
		return intentos
	}
	intentos++
	return intentos
}

func updateValue(id string, est string, fecha string, intentos int32) {
	i := 0
	for i < len(allPedidos) {
		if allPedidos[i].id == id {
			var obj = allPedidos[i]
			obj.estado = est
			obj.fechaEntrega = fecha
			obj.intentos = intentos
			allPedidos[i] = obj
			return
		}
		i++
	}
}

func sendEstado(p1 paquete) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, _ := cb.ActualizarEstado(ctx, &pb.EstadoPaquete{Id: p1.id, Estado: p1.estado})
	log.Printf("M: %v", r.GetMessage())
}

func entregarPyme(p1 paquete) {
	//Calculamos el max de intentos en base al valor del producto
	maxInt := math.Floor(float64(p1.valor) / 10)
	if maxInt == 0 {
		maxInt++
	}
	p1.estado = "En camino"
	sendEstado(p1)
	var intentos = clienteRecibe(int(maxInt))
	t := time.Now()
	p1.fechaEntrega = t.Format("2006-01-02 15:04:05")
	p1.intentos = int32(intentos)
	log.Printf("Paquete Entregado id: %v por camion: %v  Fecha Entrega:%v \n", p1.id, p1.tipoCamion, p1.fechaEntrega)
	if intentos == int(maxInt) {
		//El pedido no pudo ser entregado
		p1.estado = "No Recibido"
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r, _ := cb.ResultadoEntrega(ctx, &pb.PaqueteRecibido{Id: p1.id, Intentos: p1.intentos, Estado: p1.estado, Tipo: p1.tipo})
		log.Printf("M: %v", r.GetMessage())

	} else {
		//El pedido fue entregado con exito
		p1.estado = "Recibido"
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r, _ := cb.ResultadoEntrega(ctx, &pb.PaqueteRecibido{Id: p1.id, Intentos: p1.intentos, Estado: p1.estado, Tipo: p1.tipo})
		log.Printf("M: %v", r.GetMessage())
	}
	updateValue(p1.id, p1.estado, p1.fechaEntrega, p1.intentos)
}

func entregarRetails(p1 paquete) {
	maxInt := 3
	p1.estado = "En camino"
	sendEstado(p1)

	var intentos = clienteRecibe(maxInt)
	t := time.Now()
	p1.fechaEntrega = t.Format("2006-01-02 15:04:05")
	p1.intentos = int32(intentos)

	log.Printf("Paquete Entregado id: %v por camion: %v\n", p1.id, p1.tipoCamion)
	if intentos == 3 {
		//El pedido no pudo ser entregado
		p1.estado = "No Recibido"
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r, _ := cb.ResultadoEntrega(ctx, &pb.PaqueteRecibido{Id: p1.id, Intentos: p1.intentos, Estado: p1.estado, Tipo: p1.tipo})
		log.Printf("M: %v", r.GetMessage())
	} else {
		//El pedido fue entregado con exito
		p1.estado = "Recibido"
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r, _ := cb.ResultadoEntrega(ctx, &pb.PaqueteRecibido{Id: p1.id, Intentos: p1.intentos, Estado: p1.estado, Tipo: p1.tipo})
		log.Printf("M: %v", r.GetMessage())

	}
	updateValue(p1.id, p1.estado, p1.fechaEntrega, p1.intentos)
}

func entregaDispatcher(p1 paquete) {
	if p1.origen == "pyme" {
		entregarPyme(p1)
	} else {
		entregarRetails(p1)
	}
}

func entregarPedido(p1 paquete, p2 paquete) {
	if p2.id != "null" {
		//Entregamos el paquete con mas valor primero
		if p1.valor >= p2.valor {
			entregaDispatcher(p1)
			entregaDispatcher(p2)
		} else {
			entregaDispatcher(p2)
			entregaDispatcher(p1)
		}
	} else {
		entregaDispatcher(p1)
	}
}

func dispatcher(c pb.LogisticaClienteClient, tipoCamion string, idCamion string, wg *sync.WaitGroup) bool {

	//Recibio un paquete
	var result, primerPaquete = getPaquete(c, tipoCamion, idCamion)
	if result == true {
		//Intentamos obtener el segundo paquete
		var result2, segundoPaquete = getPaquete(c, tipoCamion, idCamion)
		if !result2 {
			log.Printf("Fallo en obtener el segundo paquete\n")
			var exito bool = false
			var rr bool
			var tercerPaquete paquete
			//Ticker con tasa de refresco 100ms
			Ti := time.NewTicker(100 * time.Millisecond)
			mychannel := make(chan bool)
			go func() {
				// Using for loop
				for {
					// Select statement
					select {
					// Case statement
					case <-mychannel:
						rr, tercerPaquete = getPaquete(c, tipoCamion, idCamion)
						if rr {
							exito = true
							Ti.Stop()
							mychannel <- true
							break
						}
					}
				}
			}()
			time.Sleep(time.Duration(tiempo) * time.Millisecond)
			// Calling Stop() method
			Ti.Stop()
			// Setting the value of channel
			mychannel <- true

			//Se obtuvo el segundo paquete
			if exito {
				entregarPedido(primerPaquete, tercerPaquete)
			} else {
				tercerPaquete.id = "null"
				entregarPedido(primerPaquete, tercerPaquete)
			}
		} else {
			//Procesamos el pedido
			entregarPedido(primerPaquete, segundoPaquete)
		}
	}
	wg.Done()
	return true
}

var cb pb.LogisticaClienteClient

func threeCamion(c pb.LogisticaClienteClient, wg *sync.WaitGroup) {
	go dispatcher(c, "retails", "Retails 1", wg)
	go dispatcher(c, "retails", "Retails 2", wg)
	go dispatcher(c, "normal", "normal", wg)
}

func main() {
	runtime.GOMAXPROCS(3)
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Bienvenido a la simulacion de PrestigioExpress--Camiones \n")
	fmt.Print("Tasa de refresco 500ms \n")
	fmt.Print("Ingrese el tiempo en milisegundos a esperar por el segundo pedido\n")
	text, _ := reader.ReadString('\n')
	tiempo, _ = strconv.Atoi(strings.TrimSuffix(text, "\n"))
	fmt.Print("Ingrese el retardo en milisegundos en entregar un pedido\n")
	reader = bufio.NewReader(os.Stdin)
	text, _ = reader.ReadString('\n')
	tiempoEntrega, _ = strconv.Atoi(strings.TrimSuffix(text, "\n"))

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	c := pb.NewLogisticaClienteClient(conn)
	cb = c

	r1 := time.NewTicker(500 * time.Millisecond)
	var w sync.WaitGroup

	for {
		select {
		case <-r1.C:
			w.Add(3)
			threeCamion(c, &w)
			w.Wait()
		}
	}
}
