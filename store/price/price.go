package price

import (
	"compound/core"
	"context"
	"time"

	"github.com/fox-one/pkg/store/db"
	"github.com/jinzhu/gorm"
)

type priceStore struct {
	db *db.DB
}

// New new price store
func New(db *db.DB) core.IPriceStore {
	return &priceStore{
		db: db,
	}
}

func init() {
	db.RegisterMigrate(func(db *db.DB) error {
		tx := db.Update().Model(core.Price{})

		if err := tx.AutoMigrate(core.Price{}).Error; err != nil {
			return err
		}

		return nil
	})
}

func (s *priceStore) Create(ctx context.Context, price *core.Price) error {
	return s.db.Update().Where("asset_id=? and block_number=?", price.AssetID, price.BlockNumber).FirstOrCreate(price).Error
}

func (s *priceStore) FindByAssetBlock(ctx context.Context, assetID string, blockNumber int64) (*core.Price, bool, error) {
	var price core.Price
	if e := s.db.View().Where("asset_id=? and block_number=?", assetID, blockNumber).Find(&price).Error; e != nil {
		return nil, gorm.IsRecordNotFoundError(e), e
	}
	return &price, false, nil
}

func (s *priceStore) Update(ctx context.Context, price *core.Price, version int64) error {
	if version > price.Version {
		price.Version = version
		return s.db.Update().Model(price).Updates(price).Error
	}

	return nil
}

func (s *priceStore) DeleteByTime(ctx context.Context, t time.Time) error {
	return s.db.Update().Where("created_at < ?", t).Delete(core.Price{}).Error
}
