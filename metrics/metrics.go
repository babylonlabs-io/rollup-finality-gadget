package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
)

var (
	// FinalizedBlocksTotal tracks the total number of finalized blocks processed
	FinalizedBlocksTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "finality_gadget_finalized_blocks_total",
		Help: "The total number of finalized blocks processed by the finality gadget",
	})
)

// Init initializes the metrics registry
func Init(logger *zap.Logger) {
	// Metrics are automatically registered via promauto
	logger.Info("Prometheus metrics initialized", zap.String("metric", "finality_gadget_finalized_blocks_total"))
}
