package snapshot

import (
	"compound/core"
	"compound/worker"
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/bluele/gcache"
	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pkg/logger"
	"github.com/fox-one/pkg/property"
	"github.com/fox-one/pkg/store/db"
	"github.com/robfig/cron/v3"
)

// Worker snapshot worker
type Worker struct {
	worker.BaseJob
	config        *core.Config
	dapp          *mixin.Client
	property      property.Store
	db            *db.DB
	marketStore   core.IMarketStore
	supplyStore   core.ISupplyStore
	borrowStore   core.IBorrowStore
	walletService core.IWalletService
	blockService  core.IBlockService
	priceService  core.IPriceOracleService
	marketService core.IMarketService
	supplyService core.ISupplyService
	borrowService core.IBorrowService
	snapshotCache gcache.Cache
}

const (
	checkPointKey = "compound_snapshot_checkpoint"
	limit         = 500
)

// New new snapshot worker
func New(
	config *core.Config,
	dapp *mixin.Client,
	property property.Store,
	db *db.DB,
	marketStore core.IMarketStore,
	supplyStore core.ISupplyStore,
	borrowStore core.IBorrowStore,
	walletService core.IWalletService,
	priceSrv core.IPriceOracleService,
	blockService core.IBlockService,
	marketSrv core.IMarketService,
	supplyService core.ISupplyService,
	borrowService core.IBorrowService,
) *Worker {
	job := Worker{
		config:        config,
		dapp:          dapp,
		property:      property,
		db:            db,
		marketStore:   marketStore,
		supplyStore:   supplyStore,
		borrowStore:   borrowStore,
		walletService: walletService,
		blockService:  blockService,
		priceService:  priceSrv,
		marketService: marketSrv,
		supplyService: supplyService,
		borrowService: borrowService,
		snapshotCache: gcache.New(limit).LRU().Build(),
	}

	l, _ := time.LoadLocation(job.config.App.Location)
	job.Cron = cron.New(cron.WithLocation(l))
	spec := "@every 1s"
	job.Cron.AddFunc(spec, job.Run)
	job.OnWork = func() error {
		return job.onWork(context.Background())
	}

	return &job
}

func (w *Worker) onWork(ctx context.Context) error {
	log := logger.FromContext(ctx)
	checkPoint, err := w.property.Get(ctx, checkPointKey)
	if err != nil {
		log.WithError(err).Errorf("read property error: %s", checkPointKey)
		return err
	}

	snapshots, next, err := w.walletService.PullSnapshots(ctx, checkPoint.String(), limit)
	if err != nil {
		log.WithError(err).Error("pull snapshots error")
		return err
	}

	if len(snapshots) == 0 {
		return errors.New("no more snapshots")
	}

	for _, snapshot := range snapshots {
		if snapshot.UserID == "" {
			continue
		}

		// if w.snapshotCache.Has(snapshot.ID) {
		// 	continue
		// }

		if err := w.handleSnapshot(ctx, snapshot); err != nil {
			return err
		}

		// w.snapshotCache.Set(snapshot.ID, nil)
	}

	if checkPoint.String() != next {
		if err := w.property.Save(ctx, checkPointKey, next); err != nil {
			log.WithError(err).Errorf("update property error: %s", checkPointKey)
			return err
		}
	}

	return nil
}

func (w *Worker) handleSnapshot(ctx context.Context, snapshot *core.Snapshot) error {
	if snapshot.UserID == w.config.BlockWallet.ClientID {
		return w.handleBlockEvent(ctx, snapshot)
	} else if snapshot.UserID == w.config.Mixin.ClientID {
		// main wallet
		var action core.Action
		e := json.Unmarshal([]byte(snapshot.Memo), &action)
		if e != nil {
			return nil
		}
		service := action[core.ActionKeyService]
		if service == core.ActionServiceSupply {
			return w.handleSupplyEvent(ctx, snapshot)
		} else if service == core.ActionServiceRedeem {
			return w.handleSupplyRedeemEvent(ctx, snapshot)
		} else if service == core.ActionServiceBorrow {
			return w.handleBorrowEvent(ctx, snapshot)
		} else if service == core.ActionServiceRepay {
			return w.handleBorrowRepayEvent(ctx, snapshot)
		} else if service == core.ActionServiceMint {
			return w.handleMintEvent(ctx, snapshot)
		} else {
			return w.handleRefundEvent(ctx, snapshot)
		}

	}
	return nil
}