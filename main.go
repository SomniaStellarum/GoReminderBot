package main

import (
	"log"
	"net/http"
)

func main() {
	s := Server{}
	s.routes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
