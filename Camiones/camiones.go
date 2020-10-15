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

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	c := pb.NewLogisticaClienteClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.SolicitudPaquetes(ctx, &pb.TipoCamion{Tipo: "retails"})
	if err != nil {
		log.Fatalf("Error al obtener el paquete: %v", err)
	}
	log.Printf("Paquete recibido: id:%v, tipo: %v, valor :%v ,intentos: %v,estado:%v, seguimiento:%s", r.GetId(), r.GetTipo(), r.GetValor(), r.GetIntentos(), r.GetEstado(), r.GetSeguimiento())

}
