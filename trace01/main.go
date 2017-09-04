package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

//010 OMIT
import (
	"github.com/opentracing/opentracing-go"

	"sourcegraph.com/sourcegraph/appdash"
	ot "sourcegraph.com/sourcegraph/appdash/opentracing"
)

func main() {

	opentracing.InitGlobalTracer(
		ot.NewTracer(
			appdash.NewRemoteCollector("localhost:7701")))

	//020 OMIT

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `<a href="/req">Click to trigger a request</a>`)
	})

	//030 OMIT
	http.HandleFunc("/req", func(w http.ResponseWriter, r *http.Request) {
		sp := opentracing.StartSpan("GET /req")
		defer sp.Finish()

		time.Sleep(10 * time.Millisecond) // do some work
		fmt.Fprintf(w, `<p>request executed
			<p><a href="/">home</a>`)
	})
	//040 OMIT

	log.Fatal(http.ListenAndServe(":8080", nil))

}
