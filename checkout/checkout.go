package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
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