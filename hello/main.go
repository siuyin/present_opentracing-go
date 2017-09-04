package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `<a href="/req">Click to trigger a request</a>`)
	})

	http.HandleFunc("/req", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Millisecond) // do some work
		fmt.Fprintf(w, `<p>request executed
			<p><a href="/">home</a>`)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))

}
