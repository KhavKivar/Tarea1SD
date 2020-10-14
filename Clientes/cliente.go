package main

import (
	"context"
	"log"
	"time"

	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	pb "google.golang.org/TAREA1SD/Logistica/paquete"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	for true {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Bienvenido a la simulacion de PrestigioExpress \n")
		fmt.Print("Ingrese 1 o 2 para realizar las siguientes tareas\n")
		fmt.Print("1) Realizar un pedido \n")
		fmt.Print("2) Ver estado de un pedido \n")
		fmt.Print("3) Salir \n")
		text, _ := reader.ReadString('\n')
		if text == "1\n" {
			fmt.Print("Ingrese 1 para simular una tienda de Pyme o 2 para retails\n")
			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')
			if text == "1\n" {
				fmt.Print("Se enviaran los pedidos que estan en el archivo pymes.csv\n")
				fmt.Print("Ingrese el tiempo de envio entre los pedidos\n")
				reader := bufio.NewReader(os.Stdin)
				text, _ := reader.ReadString('\n')
				dut, _ := strconv.Atoi(strings.TrimSuffix(text, "\n"))
				csvfile, err := os.Open("pymes.csv")
				if err != nil {
					log.Fatalln("Couldn't open the csv file", err)
				}

				r := csv.NewReader(csvfile)

				conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
				if err != nil {
					log.Fatalf("did not connect: %v", err)
				}
				defer conn.Close()

				for {

					record, err := r.Read()
					if err == io.EOF {
						break
					}
					if err != nil {
						log.Fatal(err)
					}
					if record[0] != "id" {
						c := pb.NewLogisticaClienteClient(conn)
						ctx, cancel := context.WithTimeout(context.Background(), time.Second)
						defer cancel()
						valor, _ := strconv.Atoi(record[2])
						prior, _ := strconv.Atoi(record[5])
						fmt.Printf("Pedido enviado id:%s producto:%s valor:%s tienda:%s destino:%s prioritario:%s \n", record[0], record[1], record[2], record[3], record[4], record[5])
						re, err := c.EnviarPedido(ctx, &pb.Orden{Id: record[0], Producto: record[1], Valor: int32(valor), Tienda: record[3], Destino: record[4], Prioritario: int32(prior)})
						if err != nil {
							log.Fatalf("No se puedo enviar el mensaje: %v \n", err)
						}
						log.Printf("%s", re.GetMessage())
						duration := time.Duration(dut) * time.Second
						time.Sleep(duration)
					}
				}
			}
			if text == "2\n" {
				fmt.Print("Se enviaran los pedidos que estan en el archivo retail.csv\n")
				fmt.Print("Ingrese el tiempo de envio entre los pedidos\n")
				reader := bufio.NewReader(os.Stdin)
				text, _ := reader.ReadString('\n')
				dut, _ := strconv.Atoi(strings.TrimSuffix(text, "\n"))
				csvfile, err := os.Open("retail.csv")
				if err != nil {
					log.Fatalln("Couldn't open the csv file", err)
				}

				r := csv.NewReader(csvfile)

				conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
				if err != nil {
					log.Fatalf("did not connect: %v", err)
				}
				defer conn.Close()

				for {

					record, err := r.Read()
					if err == io.EOF {
						break
					}
					if err != nil {
						log.Fatal(err)
					}
					if record[0] != "id" {
						c := pb.NewLogisticaClienteClient(conn)
						ctx, cancel := context.WithTimeout(context.Background(), time.Second)
						defer cancel()
						valor, _ := strconv.Atoi(record[2])
						fmt.Printf("Pedido enviado id:%s producto:%s valor:%s tienda:%s destino:%s \n", record[0], record[1], record[2], record[3], record[4])
						re, err := c.EnviarPedido(ctx, &pb.Orden{Id: record[0], Producto: record[1], Valor: int32(valor), Tienda: record[3], Destino: record[4]})
						if err != nil {
							log.Fatalf("No se puedo enviar el mensaje: %v \n", err)
						}
						log.Printf("%s", re.GetMessage())
						duration := time.Duration(dut) * time.Second
						time.Sleep(duration)
					}
				}
			}
		}
		if text == "3\n" {
			fmt.Print("Adios\n")
			os.Exit(0)
		}	
	
	}

}
