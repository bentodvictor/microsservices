package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

// Cria estrutura de cupom
type Coupon struct {
	Code string
}

// criar uma estrutura para armazenar infinitos arrays de cupons
type Coupons struct {
	Coupon []Coupon
}

type Result struct {
	Status string
}

var coupons Coupons

// metodo para verificar se cupom eh valido
func (c Coupons) Check(code string) string {
	// percorre todos os cupoms da lista
	for _, item := range c.Coupon {
		// se encontra cupom, valida cupom
		if code == item.Code {
			return "valid"
		}
	}
	return "invalid"
}

func main() {
	// cria cupom
	coupon := Coupon{
		Code: "abc",
	}

	// adiciona cupom na lista
	coupons.Coupon = append(coupons.Coupon, coupon)

	http.HandleFunc("/", home)
	http.ListenAndServe(":9092", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	coupon := r.PostFormValue("coupon")
	ccNumber := r.PostFormValue("ccNumber")

	// verifica cc (microsservico d)
	resultCC := makeHttpCall("http://localhost:9093", ccNumber)

	// verifica cupom
	valid := coupons.Check(coupon)

	result := Result{Status: valid}

	if result.Status == "valid" {
		if resultCC.Status == "validCard" {
			result.Status = "aproved"
		} else {
			result.Status = "invalidCard"
		}
	} else {
		result.Status = "invalidCoupon"
	}

	jsonResult, err := json.Marshal(result)
	if err != nil {
		log.Fatal("Error converting json")
	}

	fmt.Fprintf(w, string(jsonResult))

}

// invoca microsservico d (verificacao de credcard)
func makeHttpCall(urlMicroservice string, ccNumber string) Result {

	values := url.Values{}
	values.Add("ccNumber", ccNumber)

	res, err := http.PostForm(urlMicroservice, values)
	if err != nil {
		result := Result{Status: "Servidor fora do ar!"}
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