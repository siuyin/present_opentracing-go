package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/opentracing/opentracing-go"

	//010 OMIT
	jaeger "github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
	//011 OMIT
)

//012 OMIT
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
	//013 OMIT

	//014 OMIT
	// Example logger and metrics factory. Use github.com/uber/jaeger-client-go/log
	// and github.com/uber/jaeger-lib/metrics respectively to bind to real logging and metrics
	// frameworks.
	jLogger := jaegerlog.StdLogger
	jMetricsFactory := metrics.NullFactory

	// Initialize tracer with a logger and a metrics factory
	closer, err := cfg.InitGlobalTracer(
		"myWebApp", // HL
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
	)
	if err != nil {
		log.Printf("Could not initialize jaeger tracer: %s", err.Error())
		return
	}
	defer closer.Close()
	//015 OMIT

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `<a href="/req">Click to trigger a request</a>`)
	})

	//030 OMIT
	http.HandleFunc("/req", func(w http.ResponseWriter, r *http.Request) {
		sp := opentracing.StartSpan("GET /req")
		sp = sp.SetTag("Operation", "doing work by getting /req")
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
