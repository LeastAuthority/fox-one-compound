package snapshot

import (
	"compound/core"
	"compound/pkg/mtg"
	"context"
	"errors"

	"github.com/fox-one/pkg/logger"
	foxuuid "github.com/fox-one/pkg/uuid"
	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
)

// handle unpledge event
func (w *Payee) handleUnpledgeEvent(ctx context.Context, output *core.Output, userID, followID string, body []byte) error {
	log := logger.FromContext(ctx).WithField("worker", "unpledge")

	var ctokenAsset uuid.UUID
	var unpledgedAmount decimal.Decimal

	if _, err := mtg.Scan(body, &ctokenAsset, &unpledgedAmount); err != nil {
		return w.handleRefundEvent(ctx, output, userID, followID, core.ActionTypeUnpledge, core.ErrInvalidArgument)
	}

	log.Infof("ctokenAssetID:%s, amount:%s", ctokenAsset.String(), unpledgedAmount)
	unpledgedAmount = unpledgedAmount.Truncate(8)
	ctokenAssetID := ctokenAsset.String()
	market, isRecordNotFound, e := w.marketStore.FindByCToken(ctx, ctokenAssetID)
	if isRecordNotFound {
		log.Warningln("market not found")
		return w.handleRefundEvent(ctx, output, userID, followID, core.ActionTypeUnpledge, core.ErrMarketNotFound)
	}
	if e != nil {
		log.WithError(e).Errorln("find market error")
		return e
	}

	if w.marketService.IsMarketClosed(ctx, market) {
		return w.handleRefundEvent(ctx, output, userID, followID, core.ActionTypeUnpledge, core.ErrMarketClosed)
	}

	supply, isRecordNotFound, e := w.supplyStore.Find(ctx, userID, market.CTokenAssetID)
	if isRecordNotFound {
		log.Warningln("supply not found")
		return w.handleRefundEvent(ctx, output, userID, followID, core.ActionTypeUnpledge, core.ErrSupplyNotFound)
	}
	if e != nil {
		log.WithError(e).Errorln("find supply error")
		return e
	}

	if unpledgedAmount.GreaterThan(supply.Collaterals) {
		log.Errorln(errors.New("insufficient collaterals"))
		return w.handleRefundEvent(ctx, output, userID, followID, core.ActionTypeUnpledge, core.ErrInsufficientCollaterals)
	}

	//accrue interest
	if e = w.marketService.AccrueInterest(ctx, market, output.CreatedAt); e != nil {
		log.Errorln(e)
		return e
	}

	if e = w.marketStore.Update(ctx, market, output.ID); e != nil {
		log.WithError(e).Errorln("update market error")
		return e
	}

	// market transaction
	marketTransaction := core.BuildMarketUpdateTransaction(ctx, market, foxuuid.Modify(output.TraceID, "update_market"))
	if e = w.transactionStore.Create(ctx, marketTransaction); e != nil {
		log.WithError(e).Errorln("create transaction error")
		return e
	}

	// check liqudity
	liquidity, e := w.accountService.CalculateAccountLiquidity(ctx, userID)
	if e != nil {
		log.Errorln(e)
		return e
	}

	price := market.Price
	exchangeRate, e := w.marketService.CurExchangeRate(ctx, market)
	if e != nil {
		log.Errorln(e)
		return e
	}
	unpledgedTokenLiquidity := unpledgedAmount.Mul(exchangeRate).Mul(market.CollateralFactor).Mul(price)
	if unpledgedTokenLiquidity.GreaterThan(liquidity) {
		log.Errorln(errors.New("insufficient liquidity"))
		return w.handleRefundEvent(ctx, output, userID, followID, core.ActionTypeUnpledge, core.ErrInsufficientLiquidity)
	}

	if output.ID > supply.Version {
		supply.Collaterals = supply.Collaterals.Sub(unpledgedAmount).Truncate(16)
		if e = w.supplyStore.Update(ctx, supply, output.ID); e != nil {
			log.Errorln(e)
			return e
		}
	}

	// transaction
	extra := core.NewTransactionExtra()
	extra.Put(core.TransactionKeyCTokenAssetID, ctokenAssetID)
	extra.Put(core.TransactionKeyAmount, unpledgedAmount)
	extra.Put(core.TransactionKeySupply, core.ExtraSupply{
		UserID:        supply.UserID,
		CTokenAssetID: supply.CTokenAssetID,
		Collaterals:   supply.Collaterals,
	})
	transaction := core.BuildTransactionFromOutput(ctx, userID, followID, core.ActionTypeUnpledge, output, &extra)
	if e = w.transactionStore.Create(ctx, transaction); e != nil {
		log.WithError(e).Errorln("create transaction error")
		return e
	}

	// add transfer
	transferAction := core.TransferAction{
		Source:   core.ActionTypeUnpledgeTransfer,
		FollowID: followID,
	}
	return w.transferOut(ctx, userID, followID, output.TraceID, market.CTokenAssetID, unpledgedAmount, &transferAction)
}
