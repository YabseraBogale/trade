package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

	})
	err := http.ListenAndServe("localhost:3000", nil)
	if err != nil {
		log.Fatalln(err)
	}
}
