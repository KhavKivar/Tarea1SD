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
	if r.GetId() != "null" && r.GetId() != "" {

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
			f, err := os.OpenFile("C1_RETAIL.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Println(err)
			}
			defer f.Close()
			if _, err := f.WriteString(newPaquete.id + "," + newPaquete.tipo + "," + fmt.Sprint(newPaquete.valor) + "," + newPaquete.origen + "," + newPaquete.destino + "," + fmt.Sprint(newPaquete.intentos) + "," + newPaquete.fechaEntrega + "," + newPaquete.tipoCamion + "," + newPaquete.estado + "\n"); err != nil {
				log.Println(err)
			}
		}
		if idCamion == "Retails 2" {
			newPaquete.tipoCamion = idCamion
			listRetail2 = append(listRetail2, newPaquete)
			f, err := os.OpenFile("C2_RETAIL.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Println(err)
			}
			defer f.Close()
			if _, err := f.WriteString(newPaquete.id + "," + newPaquete.tipo + "," + fmt.Sprint(newPaquete.valor) + "," + newPaquete.origen + "," + newPaquete.destino + "," + fmt.Sprint(newPaquete.intentos) + "," + newPaquete.fechaEntrega + "," + newPaquete.tipoCamion + "," + newPaquete.estado + "\n"); err != nil {
				log.Println(err)
			}
		}
		if idCamion == "normal" {
			newPaquete.tipoCamion = idCamion
			listNormal = append(listNormal, newPaquete)
			newPaquete.tipoCamion = idCamion
			listRetail2 = append(listRetail2, newPaquete)
			f, err := os.OpenFile("C3_NORMAL.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Println(err)
			}
			defer f.Close()
			if _, err := f.WriteString(newPaquete.id + "," + newPaquete.tipo + "," + fmt.Sprint(newPaquete.valor) + "," + newPaquete.origen + "," + newPaquete.destino + "," + fmt.Sprint(newPaquete.intentos) + "," + newPaquete.fechaEntrega + "," + newPaquete.tipoCamion + "," + newPaquete.estado + "\n"); err != nil {
				log.Println(err)
			}
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

func clienteRecibe(maxIntentos int) (int, bool) {
	//Se Envia
	intentos := 0
	time.Sleep(time.Duration(tiempoEntrega) * time.Millisecond)
	if random80() {
		return intentos, true

	}
	//Primer intento
	intentos++
	time.Sleep(time.Duration(tiempoEntrega) * time.Millisecond)
	var rand80 = random80()
	if rand80 || maxIntentos == intentos {
		if maxIntentos == intentos {
			if rand80 {
				return intentos, true
			}
			return intentos, false
		}
		return intentos, true
	}
	//Segundo intento
	intentos++
	time.Sleep(time.Duration(tiempoEntrega) * time.Millisecond)
	var rand80_2 = random80()
	if rand80_2 || maxIntentos == intentos {
		if maxIntentos == intentos {
			if rand80_2 {
				return intentos, true
			}
			return intentos, false
		}
		return intentos, true
	}
	//tercer intento
	intentos++
	time.Sleep(time.Duration(tiempoEntrega) * time.Millisecond)
	var rand80_3 = random80()
	if rand80_3 || maxIntentos == intentos {
		if maxIntentos == intentos {
			if rand80_3 {
				return intentos, true
			}
			return intentos, false
		}
		return intentos, true
	}
	return -1, false
}

func sendEstado(p1 paquete) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, _ := cb.ActualizarEstado(ctx, &pb.EstadoPaquete{Id: p1.id, Estado: p1.estado})
	log.Printf("M: %v", r.GetMessage())
}

func entregaPedido(p1 paquete) {
	var maxInt = 3
	//Si es pyme calculamos el maxInt
	if p1.origen == "pyme" {
		maxInt = int(math.Floor(float64(p1.valor) / 10))
		if maxInt == 0 {
			maxInt++
		}
		if maxInt > 2 {
			maxInt = 2
		}
	}
	p1.estado = "En camino"
	sendEstado(p1)

	var intentos, entrego = clienteRecibe(maxInt)
	t := time.Now()
	p1.fechaEntrega = t.Format("2006-01-02 15:04:05")
	p1.intentos = int32(intentos)
	if entrego {
		p1.estado = "Recibido"
		log.Printf("Paquete entegrado al cliente  id: %v por camion: %v", p1.id, p1.tipoCamion)
	} else {
		p1.estado = "No Recibido"
		log.Printf("Paquete no pudo ser entregado al cliente id: %v por camion: %v ", p1.id, p1.tipoCamion)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, _ := cb.ResultadoEntrega(ctx, &pb.PaqueteRecibido{Id: p1.id, Intentos: p1.intentos, Estado: p1.estado, Tipo: p1.tipo})
	log.Printf("M: %v", r.GetMessage())
}

func entregarPedido(p1 paquete, p2 paquete) {

	if p2.id != "null" && p2.id != "" {
		//Entregamos el paquete con mas valor primero
		if p1.valor >= p2.valor {
			entregaPedido(p1)
			entregaPedido(p2)
		} else {
			entregaPedido(p2)
			entregaPedido(p1)
		}
	} else {
		entregaPedido(p1)
	}
}

func logicSegundo(c pb.LogisticaClienteClient, tipoCamion string, idCamion string) paquete {
	var pq paquete

	ticker := time.NewTicker(100 * time.Millisecond)
	var rr bool
	pq.id = "null"
	var tiempow = 0
	for {
		select {
		case <-ticker.C:
			rr, pq = getPaquete(c, tipoCamion, idCamion)
			tiempow = tiempow + 100
			if rr {
				log.Printf("Segundo paquete arribado en %v millisegundos \n", tiempow)
				return pq
			}
			if tiempow > tiempo {
				log.Printf("Segundo paquete no llego en el tiempo esperado: %v millisegundos\n", tiempo)
				return pq
			}
		}
	}
}

func dispatcher(c pb.LogisticaClienteClient, tipoCamion string, idCamion string, wg *sync.WaitGroup) bool {

	//Recibio un paquete
	var result, primerPaquete = getPaquete(c, tipoCamion, idCamion)
	if result == true {
		//Intentamos obtener el segundo paquete
		var result2, segundoPaquete = getPaquete(c, tipoCamion, idCamion)
		if !result2 {
			var tercerPaquete = logicSegundo(c, tipoCamion, idCamion)
			entregarPedido(primerPaquete, tercerPaquete)
		} else {
			//Procesamos el pedido
			entregarPedido(primerPaquete, segundoPaquete)
		}
	}
	wg.Done()
	return true
}

var cb pb.LogisticaClienteClient

func main() {
	runtime.GOMAXPROCS(3)
	os.Remove("C1_RETAIL.csv")
	os.Remove("C2_RETAIL.csv")
	os.Remove("C3_NORMAL.csv")
	f, err := os.OpenFile("C1_RETAIL.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	if _, err := f.WriteString("id,tipo,valor,origen,destino,intentos,fechaEntrega,tipoCamion,estado\n"); err != nil {
		log.Println(err)
	}
	f1, _ := os.OpenFile("C2_RETAIL.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	if _, err1 := f1.WriteString("id,tipo,valor,origen,destino,intentos,fechaEntrega,tipoCamion,estado\n"); err1 != nil {
		log.Println(err1)
	}
	f2, _ := os.OpenFile("C3_NORMAL.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	if _, err2 := f2.WriteString("id,tipo,valor,origen,destino,intentos,fechaEntrega,tipoCamion,estado\n"); err2 != nil {
		log.Println(err2)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Bienvenido a la simulacion de PrestigioExpress--Camiones \n")
	fmt.Print("Tasa de refresco 500ms \n")
	fmt.Print("Ingrese el tiempo en milisegundos a esperar por el segundo pedido -- 1000ms = 1seg\n")
	text, _ := reader.ReadString('\n')
	tiempo, _ = strconv.Atoi(strings.TrimSuffix(text, "\n"))
	fmt.Print("Ingrese el retardo en milisegundos en entregar un pedido -- 1000ms = 1seg\n")
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
	var w sync.WaitGroup
	var w2 sync.WaitGroup
	var w3 sync.WaitGroup
	go dispatcherR(c, &w)
	go dispatcherN(c, &w2)
	dispatcherR2(c, &w3)
}

func dispatcherR(c pb.LogisticaClienteClient, wg *sync.WaitGroup) {
	r1 := time.NewTicker(500 * time.Millisecond)
	for {
		select {
		case <-r1.C:
			wg.Add(1)
			go dispatcher(c, "retails", "Retails 1", wg)
			wg.Wait()
		}
	}
}

func dispatcherN(c pb.LogisticaClienteClient, wg *sync.WaitGroup) {
	r1 := time.NewTicker(500 * time.Millisecond)
	for {
		select {
		case <-r1.C:
			wg.Add(1)
			go dispatcher(c, "normal", "normal", wg)
			wg.Wait()
		}
	}

}
func dispatcherR2(c pb.LogisticaClienteClient, wg *sync.WaitGroup) {
	r1 := time.NewTicker(500 * time.Millisecond)
	for {
		select {
		case <-r1.C:
			wg.Add(1)
			go dispatcher(c, "retails", "Retails 2", wg)
			wg.Wait()
		}
	}

}
