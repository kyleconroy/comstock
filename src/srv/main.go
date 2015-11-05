package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"syslog"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5500"
	}

	http.HandleFunc("/logs", func(w http.ResponseWriter, r *http.Request) {
		count, err := strconv.ParseInt(r.Header.Get("Logplex-Msg-Count"), 10, 32)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		msgs, err := syslog.ParseFrame(body, int(count))
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		fmt.Fprintf(w, "parsed %d messages\n", len(msgs))
	})

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
