package main

import(
	"fmt"
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	"os"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

var ingresos float64
var gastos float64
var balance_total float64

var balance_prod float64

type final struct{
	Id 			string `json: "id"`
	Seguimiento	string `json: "seguimiento"`
	Tipo		string `json: "tipo"`
	Valor		int32 `json: "valor"`
	Intentos	int32 `json: "intentos`
	Estado		string `json: "estado"`
	Balance		float64 
}

var ordenes []final

func jsonToStruct(a []byte) final{
	var est final

	err := json.Unmarshal([]byte(a), &est)
	if err != nil{
		fmt.Println(err)
	}

	return est
}

func calcularBalance(nuevo final){

	balance_prod = 0
	
	if(nuevo.Tipo == "retail"){
		balance_prod= float64(nuevo.Valor) - float64((10 * (nuevo.Intentos - 1)))
		ingresos = float64(nuevo.Valor)
		gastos = float64(10 * (nuevo.Intentos - 1))
	}else if(nuevo.Tipo == "prioritario"){
		if(nuevo.Estado == "Recibido"){
			balance_prod = float64(nuevo.Valor - (10 * (nuevo.Intentos - 1)))
			ingresos = float64(nuevo.Valor)
			gastos = float64(10 * (nuevo.Intentos - 1))
		}else{
			balance_prod = float64(nuevo.Valor - (10 * (nuevo.Intentos - 1)))
			ingresos = (0.3 * float64(nuevo.Valor))
			gastos = float64(10 * (nuevo.Intentos - 1))
		}
	}else{
		if(nuevo.Estado == "Recibido"){
			balance_prod = float64(nuevo.Valor - (10 * (nuevo.Intentos - 1)))
			ingresos = float64(nuevo.Valor)
			gastos = float64(10 * (nuevo.Intentos - 1))
			}else{
				balance_prod = float64(nuevo.Valor - (10 * (nuevo.Intentos - 1)))
				gastos = float64(10 * (nuevo.Intentos - 1))
			}
	}
}

func actualizarCSV(nuevo final){

	f, err := os.OpenFile("ordenes.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString("ID: " + nuevo.Id + "," + "Seguimiento: " + nuevo.Seguimiento + "," + "Valor: " + fmt.Sprint(nuevo.Valor) + "," + "Tipo: " + nuevo.Tipo + "," + "Intentos: " + fmt.Sprint(nuevo.Intentos) + "," + "Estado: " + nuevo.Estado + "," + "Balance: " + fmt.Sprintf("%f",nuevo.Balance) + "\n"); err != nil {
		log.Println(err)
	}
}

func main(){

	balance_total = 0
	ingresos = 0
	gastos = 0

	log.Printf("Bienvenido al sistema de finanzas de PrestigioExpress")
	
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

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
			log.Printf("-------------------------------------------------------------------------------------")
			log.Printf("El balance actual es:")
			log.Printf("      Ingresos: %f dignipesos		Gastos: %f dignipesos", ingresos, gastos)
			log.Printf("			Balance Total: %f dignipesos", balance_total)
			log.Printf("-------------------------------------------------------------------------------------")
			log.Printf(" [*] Esperando actualizaciones de logistica. Presiona CTRL + C para salir.")
		}
	}()

	log.Printf(" [*] Esperando actualizaciones de logistica. Presiona CTRL + C para salir.")
	<-forever
}






