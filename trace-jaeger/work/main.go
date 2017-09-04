package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/opentracing/opentracing-go"
	otl "github.com/opentracing/opentracing-go/log"

	jaeger "github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
)

//010 OMIT
func main() {
	// Sample configuration for testing. Use constant sampling to sample every trace
	// and enable LogSpan to log every span via configured Logger.
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
		},
	}

	// Example logger and metrics factory. Use github.com/uber/jaeger-client-go/log
	// and github.com/uber/jaeger-lib/metrics respectively to bind to real logging and metrics
	// frameworks.
	jLogger := jaegerlog.StdLogger
	jMetricsFactory := metrics.NullFactory

	// Initialize tracer with a logger and a metrics factory
	//011 OMIT
	closer, err := cfg.InitGlobalTracer(
		"myWebSvc",
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
	)
	//012 OMIT
	if err != nil {
		log.Printf("Could not initialize jaeger tracer: %s", err.Error())
		return
	}
	defer closer.Close()

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
// we simulate a message queue call - say over apache kafka, NATS, zeroMQ etc.
func networkDBQuery(b []byte) {
	carrier := map[string]string{}
	json.Unmarshal(b, &carrier)

	sc, _ := opentracing.GlobalTracer().Extract(opentracing.TextMap,
		opentracing.TextMapCarrier(carrier))
	sp := opentracing.StartSpan("Net DB query", opentracing.FollowsFrom(sc))
	defer sp.Finish()

	sp.LogFields(otl.String("start", "about to start work"))
	time.Sleep(13 * time.Millisecond) // do some work
	sp.LogFields(otl.Int("finishedWorkExpectedToTakeMilliseconds", 13))
}

//040 OMIT
