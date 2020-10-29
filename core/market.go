package core

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

// Market market info
type Market struct {
	// asset id
	AssetID string `sql:"size:36;PRIMARY_KEY" json:"asset_id"`
	// symbol
	Symbol string `sql:"size:20;unique_index:symbol_idx" json:"symbol"`
	// 已借出的资产
	TotalBorrows        decimal.Decimal `sql:"type:decimal(20,8)" json:"total_borrows"`
	TotalBorrowInterest decimal.Decimal `sql:"type:decimal(20,8)" json:"total_borrow_interest"`
	TotalSupplyInterest decimal.Decimal `sql:"type:decimal(20,8)" json:"total_supply_interest"`
	// ctoken asset id
	CTokenAssetID string `sql:"size:36" json:"ctoken_asset_id"`
	// ctoken symbol
	CTokenSymbol string `sql:"size:20" json:"ctoken_symbol"`
	// CToken 数量
	CTokens decimal.Decimal `sql:"type:decimal(20,8)" json:"ctokens"`
	// 初始兑换率
	InitExchangeRate decimal.Decimal `sql:"type:decimal(20,8);default:1" json:"init_exchange_rate"`
	// 平台储备金率 (0, 1), 默认为 0.10
	ReserveFactor decimal.Decimal `sql:"type:decimal(20,8)" json:"reserve_factor"`
	// 清算激励因子 (0, 1)
	LiquidationIncentive decimal.Decimal `sql:"type:decimal(20,8)" json:"liquidation_incentive"`
	//抵押因子 = 可借贷价值 / 抵押资产价值，目前compound设置为0.75
	CollateralFactor decimal.Decimal `sql:"type:decimal(20,8)" json:"collateral_factor"`
	//触发清算因子 [0.05, 0.9]
	CloseFactor decimal.Decimal `sql:"type:decimal(20,8)" json:"close_factor"`
	//基础利率 per year, 0.025
	BaseRate decimal.Decimal `sql:"type:decimal(20,8)" json:"base_rate"`
	// The multiplier of utilization rate that gives the slope of the interest rate. per year
	Multiplier decimal.Decimal `sql:"type:decimal(20,8)" json:"multiplier"`
	// The multiplierPerBlock after hitting a specified utilization point. per year
	JumpMultiplier decimal.Decimal `sql:"type:decimal(20,8)" json:"jump_multiplier"`
	// Kink
	Kink      decimal.Decimal `sql:"type:decimal(20,8)" json:"kink"`
	CreatedAt time.Time       `sql:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time       `sql:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// IMarketStore asset store interface
type IMarketStore interface {
	Save(ctx context.Context, market *Market) error
	Find(ctx context.Context, assetID, symbol string) (*Market, error)
	All(ctx context.Context) ([]*Market, error)
}

// IMarketService market interface
type IMarketService interface {
	SaveUtilizationRate(ctx context.Context, symbol string, rate decimal.Decimal, block int64) error
	GetUtilizationRate(ctx context.Context, symbol string, block int64) (decimal.Decimal, error)
	SaveBorrowRatePerBlock(ctx context.Context, symbol string, rate decimal.Decimal, block int64) error
	GetBorrowRatePerBlock(ctx context.Context, symbol string, block int64) (decimal.Decimal, error)
	GetBorrowRate(ctx context.Context, symbol string, block int64) (decimal.Decimal, error)
	SaveSupplyRatePerBlock(ctx context.Context, symbol string, rate decimal.Decimal, block int64) error
	GetSupplyRatePerBlock(ctx context.Context, symbol string, block int64) (decimal.Decimal, error)
	GetSupplyRate(ctx context.Context, symbol string, block int64) (decimal.Decimal, error)

	CurUtilizationRate(ctx context.Context, market *Market) (decimal.Decimal, error)
	CurExchangeRate(ctx context.Context, market *Market) (decimal.Decimal, error)
	CurBorrowRate(ctx context.Context, market *Market) (decimal.Decimal, error)
	CurBorrowRatePerBlock(ctx context.Context, market *Market) (decimal.Decimal, error)
	CurSupplyRate(ctx context.Context, market *Market) (decimal.Decimal, error)
	CurSupplyRatePerBlock(ctx context.Context, market *Market) (decimal.Decimal, error)
	CurTotalCash(ctx context.Context, market *Market) (decimal.Decimal, error)
	CurTotalBorrow(ctx context.Context, market *Market) (decimal.Decimal, error)
	CurTotalReserves(ctx context.Context, market *Market) (decimal.Decimal, error)
	CurTotalBorrowInterest(ctx context.Context, market *Market) (decimal.Decimal, error)
	CurTotalSupplyInterest(ctx context.Context, market *Market) (decimal.Decimal, error)

	Mint(ctx context.Context, market *Market) error
}
