package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/streadway/amqp"
	"github.com/joho/godotenv"
	"github.com/wesleywillians/go-rabbitmq/queue"
)

// estrutura para armazenar o resultado das requisições
type Result struct {
	Status string
}

// estrutura para ordems
type Order struct {
	Coupon string
	CcNumber string
}

// funcao init eh invocada assim que programa eh executado
// armazena variaves .env
func init() {
	err := godotenv.Load();
	if err != nil {
		log.Fatal("Error loading .env");
	}
}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/process", process)
	http.ListenAndServe(":9090", nil)
}

// carrega templates e request e response (w = request e r = response)
func home(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("templates/home.html"))
	// home nao tem resultados
	t.Execute(w, Result{})	
}

// função para processamento do checkout
func process(w http.ResponseWriter, r *http.Request) {

	// obtém coupon e ccNumber para criar a ordem
	coupon := r.PostFormValue("coupon")
	ccNumber := r.PostFormValue("cc-number")

	order := Order {
		Coupon: coupon,
		CcNumber: ccNumber,
	}

	// converter para json
	jsonOrder, err := json.Marshal(order)
	if err != nil {
		log.Fatal("Error parsing to json")
	}

	// cria nova fila e se conecta
	rabbitMQ := queue.NewRabbitMQ()
	ch := rabbitMQ.Connect()

	// cria exchange para ordens
	if err := ch.ExchangeDeclare(
		"orders_ex", 	// name
		"direct",     	// type
		true,        	// durable
		false,         	// auto-deleted
		false,         	// internal
		false,         	// no-wait
		nil,           	// arguments
	); err != nil {
		log.Fatal(err, "Failed to declare an exchange")
	}

	// cria dlx, quando algum microsservico cair
	if err := ch.ExchangeDeclare(
		"dlx", 	// name
		"direct",     	// type
		true,        	// durable
		false,         	// auto-deleted
		false,         	// internal
		false,         	// no-wait
		nil,           	// arguments
	); err != nil {
		log.Fatal(err, "Failed to declare an exchange")
	}

	// cria queue de ordens, indicando dlx para enviar a ordem caso microsservico caia no processo 
	if _, err := ch.QueueDeclare(
		"orders", // name
		true,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		amqp.Table{"x-dead-letter-exchange": "dlx"},
	); err != nil {
		log.Fatal(err, "Failed to declare an QueueDeclare")
	}
	
	// cria queue de orders_dlq, e tenta reenviar para  orders_ex (sua dead letter) a cada 3s  
	if _, err := ch.QueueDeclare(
		"orders_dlq", // name
		true,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		amqp.Table{"x-dead-letter-exchange": "orders_ex", "x-message-ttl": 3000},
	); err != nil {
		log.Fatal(err, "Failed to declare an QueueDeclare")
	}

	// faz bind orders_ex - orders
	if err := ch.QueueBind(
		"orders",        // queue name
		"",             // routing key
		"orders_ex", // exchange
		false,
		nil,
	); err != nil {
		log.Fatal(err, "Failed to declare an QueueBind")
	}

	// faz bind dlx - orders_dlq
	if err := ch.QueueBind(
		"orders_dlq",        // queue name
		"",             // routing key
		"dlx", // exchange
		false,
		nil,
	); err != nil {
		log.Fatal(err, "Failed to declare an QueueBind")
	}
	
	defer ch.Close()	// Fecha conexão após tudo ser executado


	//  envia a mensagem: mensagem, tipo, exchange e routing key
	err = rabbitMQ.Notify(string(jsonOrder), "application/json", "orders_ex", "")
	if err != nil {
		log.Fatal("Error sending message to the queue")
	}

	// invoca template
	t := template.Must(template.ParseFiles("templates/process.html"))
	t.Execute(w, "")
}