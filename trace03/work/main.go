package main

import (
	"encoding/json"
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

		carrier := map[string]string{}
		opentracing.GlobalTracer().Inject(spCtx,
			opentracing.TextMap,
			opentracing.TextMapCarrier(carrier))

		b, _ := json.Marshal(carrier)
		networkDBQuery(b)
	})
	//020 OMIT

	log.Fatal(http.ListenAndServe(":8081", nil))
}

//030 OMIT
// we simulate a message queue call - say over apache kafka
func networkDBQuery(b []byte) {
	carrier := map[string]string{}
	json.Unmarshal(b, &carrier)

	sc, _ := opentracing.GlobalTracer().Extract(opentracing.TextMap,
		opentracing.TextMapCarrier(carrier))
	sp := opentracing.StartSpan("Net DB query", opentracing.FollowsFrom(sc))
	defer sp.Finish()

	time.Sleep(13 * time.Millisecond) // do some work
}

//040 OMIT
