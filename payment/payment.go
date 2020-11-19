package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type Result struct {
	Status string
}

func main() {
	http.HandleFunc("/", home)
	http.ListenAndServe(":9091", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	// obtem cupom e numero da conta
	coupon := r.PostFormValue("coupon")
	ccNumber := r.PostFormValue("ccNumber")

	// verifica cupom (microsservico c)
	resultCoupon := makeHttpCall("http://localhost:9092", coupon, ccNumber)

	// cria estrutura por padrao como declined
	result := Result{Status: "declined"}

	if resultCoupon.Status == "invalidCupom" {
		result.Status = "invalid coupon"
	} else {
		if resultCoupon.Status == "invalidCc" {
			result.Status = "invalid cc number"
		} else {
			result.Status = "aproved"
		}
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		log.Fatal("Error processing json")
	}

	fmt.Fprintf(w, string(jsonData))
}

// invoca microsservico c (verificacao de cupom)
func makeHttpCall(urlMicroservice string, coupon string, ccNumber string) Result {

	values := url.Values{}
	values.Add("coupon", coupon)
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
