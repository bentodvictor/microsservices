package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/hashicorp/go-retryablehttp"
)

// estrutura para armazenar o resultado das requisições
type Result struct {
	Status string
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

	// invoca microsserviço dois (2)
	result := makeHttpCall("http://localhost:9091", r.FormValue("coupon"), r.FormValue("cc-number"))

	t := template.Must(template.ParseFiles("templates/home.html"))
	t.Execute(w, result)
}

// sempre retorna uma estrutura Result
func makeHttpCall(urlMicroservice string, coupon string, ccNumber string) Result {

	// recebe dados do html
	values := url.Values{}
	values.Add("coupon", coupon)
	values.Add("ccNumber", ccNumber)

	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 5

	// faz requisição e trata possivel erro
	res, err := retryClient.PostForm(urlMicroservice, values)
	if err != nil {
		result := Result{Status: "Server is down!"}
		return result
	}
	// termina conexão ao final da execução - somente no final
	defer res.Body.Close()

	// pega o corpo do resultado da chamada anterior
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Error processing result")
	}

	result := Result{}

	// pega o resultado (body) e converte para struct
	json.Unmarshal(data, &result)

	return result

}
