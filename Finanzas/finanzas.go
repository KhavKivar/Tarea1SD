package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

// Funcion que printea errores de RabbitMQ
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

var ingresos float64
var gastos float64
var balance_total float64

var balance_prod float64

var todos_los_datos [3]float64

// Struct en las que los mensajes en JSON se mapean
type final struct {
	Id          string `json: "id"`
	Seguimiento string `json: "seguimiento"`
	Tipo        string `json: "tipo"`
	Valor       int32  `json: "valor"`
	Intentos    int32  `json: "intentos`
	Estado      string `json: "estado"`
	Balance     float64
}

var ordenes []final

// Funcion que transforma el mensaje en JSON que llega de Rabbit en struct
func jsonToStruct(a []byte) final {
	var est final

	err := json.Unmarshal([]byte(a), &est)
	if err != nil {
		fmt.Println(err)
	}

	return est
}

// Funcion que calcula el balance de cada producto que se actualiza a Recibido o No Recibido dependiendo de su tipo
func calcularBalance(nuevo final) {

	balance_prod = 0
	ingresos = 0
	gastos = 0

	if nuevo.Tipo == "retail" {
		balance_prod = float64(nuevo.Valor) - float64(10*nuevo.Intentos)
		ingresos = float64(nuevo.Valor)
		gastos = float64(10 * nuevo.Intentos)
	} else if nuevo.Tipo == "prioritario" {
		if nuevo.Estado == "Recibido" {
			balance_prod = float64(nuevo.Valor) - float64(10*nuevo.Intentos)
			ingresos = float64(nuevo.Valor)
			gastos = float64(10 * nuevo.Intentos)
		} else {
			balance_prod = float64(nuevo.Valor) - float64(10*nuevo.Intentos)
			ingresos = (0.3 * float64(nuevo.Valor))
			gastos = float64(10 * nuevo.Intentos)
		}
	} else {
		if nuevo.Estado == "Recibido" {
			balance_prod = float64(nuevo.Valor) - float64(10*nuevo.Intentos)
			ingresos = float64(nuevo.Valor)
			gastos = float64(10 * nuevo.Intentos)
		} else {
			balance_prod = float64(nuevo.Valor) - float64(10*nuevo.Intentos)
			gastos = float64(10 * nuevo.Intentos)
		}
	}
}

// Funcion que crea el csv para llevar registro de ordenes completadas
func actualizarCSV(nuevo final) {

	f, err := os.OpenFile("ordenes.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString("ID: " + nuevo.Id + "," + "Seguimiento: " + nuevo.Seguimiento + "," + "Valor: " + fmt.Sprint(nuevo.Valor) + "," + "Tipo: " + nuevo.Tipo + "," + "Intentos: " + fmt.Sprint(nuevo.Intentos) + "," + "Estado: " + nuevo.Estado + "," + "Balance: " + fmt.Sprintf("%f", nuevo.Balance) + "\n"); err != nil {
		log.Println(err)
	}
}

// La funcion main tiene toda la conexion de receiver de RabbitMQ, sumado a una funcion de go para actualizar el balance e imprimir datos en pantalla
func main() {

	balance_total = 0

	log.Printf("Bienvenido al sistema de finanzas de PrestigioExpress")

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Fallo en conectar con RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Fallo al abrir el canal")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Fallo al crear la queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Fallo en sacar un mensaje de la queue")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Actualizacion de logistica recibida. Recalculando balance.")
			actualizacion := jsonToStruct(d.Body)
			calcularBalance(actualizacion)
			actualizacion.Balance = balance_prod
			actualizarCSV(actualizacion)
			ordenes = append(ordenes, actualizacion)
			balance_total = balance_total + ingresos - gastos
			todos_los_datos[1] = 0
			todos_los_datos[2] = 0

			todos_los_datos[0] = balance_total
			todos_los_datos[1] = todos_los_datos[1] + ingresos
			todos_los_datos[2] = todos_los_datos[2] + gastos

			log.Printf("-------------------------------------------------------------------------------------")
			log.Printf("El estado de la entrega es %s, es de tipo %s, su valor es de %d y los reintentos fueron %d", actualizacion.Estado, actualizacion.Tipo, actualizacion.Valor, actualizacion.Intentos)
			log.Printf("El balance actual es:")
			log.Printf("   Ingreso del paquete: %f dignipesos || Gastos asociados a la entrega: %f dignipesos", ingresos, gastos)
			log.Printf("			Balance Total: %f dignipesos", balance_total)
			log.Printf("-------------------------------------------------------------------------------------")
			log.Printf(" [*] Esperando actualizaciones de logistica. Presiona CTRL + C para salir.")
		}
	}()

	log.Printf(" [*] Esperando actualizaciones de logistica. Presiona CTRL + C para salir.")
	<-forever
}
