package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/Shopify/sarama"
	"github.com/kyleconroy/comstock/syslog"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5500"
	}

	topic := os.Getenv("COMSTOCK_KAFKA_TOPIC")
	producer, err := newProducer()
	if err != nil {
		log.Fatalf("Error creating Kafka producer: %s", err)
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

		for _, ll := range msgs {
			msg := &sarama.ProducerMessage{
				Topic: topic,
				Key:   sarama.ByteEncoder(fmt.Sprintf("%s|%s", ll.Hostname, ll.Application)),
				Value: sarama.ByteEncoder(ll.Body),
			}
			producer.Input() <- msg
		}

		fmt.Fprintf(w, "parsed %d messages\n", len(msgs))
	})

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
