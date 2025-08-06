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

	// BlockVoters tracks which FPs voted for each block
	BlockVoters = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "finality_gadget_block_voters",
		Help: "List of finality providers who voted for each block",
	}, []string{"block_height", "fp_pubkeys"})

	// FpVotingPowerPerBlock tracks each FP's voting power at each block
	FpVotingPowerPerBlock = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "finality_gadget_fp_voting_power_per_block",
		Help: "Voting power of each finality provider at each block height",
	}, []string{"block_height", "fp_pubkey"})

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
		zap.String("block_voters_metric", "finality_gadget_block_voters"),
		zap.String("fp_voting_power_metric", "finality_gadget_fp_voting_power_per_block"),
		zap.String("latest_finalized_metric", "finality_gadget_latest_finalized_block_height"))
}
