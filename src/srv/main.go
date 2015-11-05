package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5500"
	}

	http.HandleFunc("/logs", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Header.Get("Logplex-Msg-Count"))
		fmt.Println(r.Header.Get("Logplex-Frame-Id"))
		body, err := ioutil.ReadAll(r.Body)
		if err == nil {
			fmt.Println(string(body))
		}
		fmt.Fprintf(w, "ok")
	})

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
