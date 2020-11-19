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
	fmt.Print(len(ccNumber))

	result := Result{Status: "invalidCc"}

	// verifica cc number
	if (len(ccNumber) < 16 || len(ccNumber) > 16) {
		result.Status = "invalidCc"
	} else {
		result.Status = "validCc"
	}

	jsonResult, err := json.Marshal(result)
	if err != nil {
		log.Fatal("Error converting json")
	}

	fmt.Fprintf(w, string(jsonResult))

}