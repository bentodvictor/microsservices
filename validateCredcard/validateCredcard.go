package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Result struct {
	Status string
}

func main() {
	http.HandleFunc("/", home)
	http.ListenAndServe(":9093", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	ccNumber := r.PostFormValue("ccNumber")
	result := Result{Status: "invalidCard"}

	// verifica cc number
	if len(ccNumber) < 16 || len(ccNumber) > 16 {
		result.Status = "invalidCard"
	} else {
		result.Status = "validCard"
	}

	jsonResult, err := json.Marshal(result)
	if err != nil {
		log.Fatal("Error converting json")
	}

	fmt.Fprintf(w, string(jsonResult))

}
