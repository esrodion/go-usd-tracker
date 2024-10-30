package metrics

import (
	"context"
	"go-usdtrub/pkg/logger"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type HandlerFunc = func(w http.ResponseWriter, r *http.Request)
type wrapperKey string

const WrapperKey wrapperKey = "wrapper"

var ZeroHandler *handlerWrapper = nil // pass to context value to disable metrics collection

func HandlerHTTP() http.Handler {
	return promhttp.Handler()
}

func WrapHandlerFunc(name, route string, handler HandlerFunc) HandlerFunc {
	wrapper := newHandlerWrapper(name, route)

	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		wrapper.count.Inc()
		handler(w, r.WithContext(context.WithValue(r.Context(), WrapperKey, wrapper)))
		duration := time.Since(startTime).Seconds()
		wrapper.duration.Observe(duration)
	}
}

func GRPCMethodHead(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, WrapperKey, grpcWrapper(name))
}

func GRPCMethodTail(ctx context.Context, name string, startTime time.Time) {
	wrapper, ok := ctx.Value(WrapperKey).(*handlerWrapper)
	if !ok {
		logger.Logger().Sugar().Named("GRPC metrics").Error("context does not contain metrics wrapper")
		return
	}
	duration := time.Since(startTime).Seconds()
	wrapper.dbAccess.Observe(duration)
}

func DBAccessDuration(ctx context.Context, startTime time.Time) {
	wrapper, ok := ctx.Value(WrapperKey).(*handlerWrapper)
	if !ok {
		logger.Logger().Sugar().Named("DB metrics").Error("context does not contain metrics wrapper")
		return
	}
	if wrapper == nil {
		return
	}
	duration := time.Since(startTime).Seconds()
	wrapper.dbAccess.Observe(duration)
}

//// Service

type handlerWrapper struct {
	count    prometheus.Counter
	duration prometheus.Histogram
	dbAccess prometheus.Histogram
}

func newHandlerWrapper(name, route string) *handlerWrapper {
	wrapper := &handlerWrapper{
		count: prometheus.NewCounter(prometheus.CounterOpts{
			Name: name + "_requests_total",
			Help: "Total number of requests to " + route,
		}),

		duration: prometheus.NewHistogram(prometheus.HistogramOpts{
			Name:    name + "_duration_seconds",
			Help:    "Request duration in seconds to " + route,
			Buckets: prometheus.LinearBuckets(0.1, 0.1, 10),
		}),

		dbAccess: prometheus.NewHistogram(prometheus.HistogramOpts{
			Name:    name + "_DB_duration_seconds",
			Help:    "DB usage duration in seconds of " + route,
			Buckets: prometheus.LinearBuckets(0.1, 0.1, 10),
		}),
	}

	prometheus.MustRegister(wrapper.count)
	prometheus.MustRegister(wrapper.duration)
	prometheus.MustRegister(wrapper.dbAccess)

	return wrapper
}

//// gRPC wrappers

var grpcWrappersM sync.RWMutex
var grpcWrappers map[string]*handlerWrapper

func init() {
	grpcWrappers = make(map[string]*handlerWrapper)
}

func grpcWrapper(key string) *handlerWrapper {
	grpcWrappersM.RLock()
	hw, ok := grpcWrappers[key]
	grpcWrappersM.RUnlock()

	if !ok {
		grpcWrappersM.Lock()
		defer grpcWrappersM.Unlock()

		hw, ok = grpcWrappers[key]
		if ok {
			return hw
		}

		hw = newHandlerWrapper(key, "gRPC:"+key)
		grpcWrappers[key] = hw
	}

	return hw
}
