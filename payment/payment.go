package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
	"github.com/streadway/amqp"
	"github.com/wesleywillians/go-rabbitmq/queue"
)

type Result struct {
	Status string
}

// estrutura para ordem de servico
type Order struct {
	ID       uuid.UUID
	Coupon   string
	CcNumber string
}

//  Toda vez que uma nova ordem for gerada, adiciona novo ID
func NewOrder() Order {
	return Order{ID: uuid.NewV4()}
}

const (
	invalidCoupon   = "invalidCoupon"
	invalidCard     = "invalidCard"
	ConnectionError = "connection error"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env")
	}
}

func main() {
	// canal de comunicao do go para amqp delivery eh criado
	messageChannel := make(chan amqp.Delivery, 0)

	// canal de comunicacao para bebbitmq
	rabbitMQ := queue.NewRabbitMQ()
	ch := rabbitMQ.Connect()
	defer ch.Close()

	// tadas as mensagens consumidas pelo rabbit sao enviadas para messageChanel
	rabbitMQ.Consume(messageChannel)

	// mensageem eh processada no loop
	for msg := range messageChannel {
		process(msg)
	}

}

// Toda vez que recebemos uma msg do rabbit ira cair no process
func process(msg amqp.Delivery) {
	order := NewOrder()
	json.Unmarshal(msg.Body, &order)

	// verifica cupom (microsservico c)
	resultCoupon := makeHttpCall("http://localhost:9092", order.Coupon, order.CcNumber)

	switch resultCoupon.Status {
	case invalidCoupon:
		log.Println("Order: ", order.ID, ": invalid coupon!")
	case invalidCard:
		log.Println("Order: ", order.ID, ": invalid cc number!")
	case ConnectionError:
		// nao coloca novamente na fila caso microsservico c esteja fora do ar
		msg.Reject(false)
		log.Println("Order: ", order.ID, ": could not process!")
	default:
		log.Println("Order: ", order.ID, ": Processed!")
	}
}

// invoca microsservico c (verificacao de cupom)
func makeHttpCall(urlMicroservice string, coupon string, ccNumber string) Result {

	values := url.Values{}
	values.Add("coupon", coupon)
	values.Add("ccNumber", ccNumber)

	res, err := http.PostForm(urlMicroservice, values)
	if err != nil {
		result := Result{Status: ConnectionError}
		return result
	}

	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Error processing result")
	}

	result := Result{}

	json.Unmarshal(data, &result)

	return result

}
