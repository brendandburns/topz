package main

import (
	"net/http"

	"log"

	"github.com/brendandburns/topz/pkg/topz"
)

func main() {
	http.HandleFunc("/topz", topz.HandleTopz)

	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
