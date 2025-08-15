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

	// FpLatestBlockVoted tracks the latest block height each FP voted on
	FpLatestBlockVoted = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "finality_gadget_fp_latest_block_voted",
		Help: "Latest block height that each finality provider voted on",
	}, []string{"fp_pubkey"})

	// FpMissedBlocks tracks the total number of blocks missed by each FP
	FpMissedBlocks = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "finality_gadget_fp_missed_blocks_total",
		Help: "Total number of blocks missed by each finality provider",
	}, []string{"fp_pubkey"})

	// FpLatestVotingPower tracks each FP's voting power for the latest processed block
	FpLatestVotingPower = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "finality_gadget_fp_latest_voting_power",
		Help: "Latest voting power of each finality provider",
	}, []string{"fp_pubkey"})

	// LatestFinalizedBlockHeight tracks the height of the latest finalized block
	LatestFinalizedBlockHeight = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "finality_gadget_latest_finalized_block_height",
		Help: "Height of the latest finalized block",
	})
)

// Init initializes the metrics registry
func Init(logger *zap.Logger) {
	// Metrics are automatically registered via promauto
	logger.Info("Prometheus metrics initialized",
		zap.String("blocks_metric", "finality_gadget_finalized_blocks_total"),
		zap.String("fp_voting_metric", "finality_gadget_fp_latest_block_voted"),
		zap.String("fp_missed_blocks_metric", "finality_gadget_fp_missed_blocks_total"),
		zap.String("fp_voting_power_metric", "finality_gadget_fp_latest_voting_power"),
		zap.String("latest_finalized_metric", "finality_gadget_latest_finalized_block_height"))
}
