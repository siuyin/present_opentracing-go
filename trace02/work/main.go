package main

import (
	"log"
	"net/http"
	"time"

	"github.com/opentracing/opentracing-go"

	"sourcegraph.com/sourcegraph/appdash"
	ot "sourcegraph.com/sourcegraph/appdash/opentracing"
)

//010 OMIT
func main() {
	opentracing.InitGlobalTracer(
		ot.NewTracer(
			appdash.NewRemoteCollector("localhost:7701")))

	http.HandleFunc("/work", func(w http.ResponseWriter, r *http.Request) {
		spCtx, _ := opentracing.GlobalTracer().Extract(opentracing.TextMap,
			opentracing.HTTPHeadersCarrier(r.Header))
		sp := opentracing.StartSpan("GET /work", opentracing.ChildOf(spCtx))
		defer sp.Finish()

		dbQuery(spCtx)
	})
	//020 OMIT

	log.Fatal(http.ListenAndServe(":8081", nil))
}

//030 OMIT
func dbQuery(sc opentracing.SpanContext) {
	sp := opentracing.StartSpan("DB query", opentracing.ChildOf(sc))
	defer sp.Finish()

	time.Sleep(13 * time.Millisecond) // do some work
}

//040 OMIT
