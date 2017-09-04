package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"sourcegraph.com/sourcegraph/appdash"

	//010 OMIT

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

		// call a remote service
		wReq, _ := http.NewRequest("GET", "http://localhost:8081/work", nil)
		sp.Tracer().Inject(sp.Context(),
			opentracing.TextMap,
			opentracing.HTTPHeadersCarrier(wReq.Header))
		resp, err := http.DefaultClient.Do(wReq)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		fmt.Fprintf(w, `<p>request executed
			<p><a href="/">home</a>`)
	})
	//040 OMIT

	log.Fatal(http.ListenAndServe(":8080", nil))

}
