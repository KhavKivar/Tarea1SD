package main

import (
	"context"
	"log"
	"time"

	pb "google.golang.org/TAREA1SD/Logistica/paquete"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50052"
)

func dispatcher(c pb.LogisticaClienteClient, tipoCamion string, idCamion string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SolicitudPaquetes(ctx, &pb.TipoCamion{Tipo: tipoCamion})
	if err != nil {
		log.Fatalf("Error al obtener el paquete: %v", err)
	}
	if r.GetId() != "null" {
		log.Printf("Paquete recibido: id:%v, Origen: %v, Destino:%v, Camion encargado: %v", r.GetId(), r.GetOrigen(), r.GetDestino(), idCamion)
	}

}
func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	c := pb.NewLogisticaClienteClient(conn)
	r1 := time.NewTicker(500 * time.Millisecond)
	r2 := time.NewTicker(500 * time.Millisecond)
	n0 := time.NewTicker(500 * time.Millisecond)
	for {
		select {
		case <-r1.C:
			go dispatcher(c, "retails", "Retails 1")
		case <-r2.C:
			go dispatcher(c, "retails", "Retails 2")
		case <-n0.C:
			go dispatcher(c, "normal", "normal")
		}
	}

}
