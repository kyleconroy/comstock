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
		body, err := ioutil.ReadAll(r.Body)
		if err == nil {
			fmt.Println(body)
		}
		fmt.Fprintf(w, "ok")
	})

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
