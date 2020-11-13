package block

import (
	"compound/core"
	"compound/internal/compound"
	"context"
	"time"
)

type service struct {
	config *core.Config
}

// New new block service
func New(config *core.Config) core.IBlockService {
	return &service{
		config: config,
	}
}

//CurrentBlock 获取当前块高度
func (s *service) CurrentBlock(ctx context.Context) (int64, error) {
	current, e := compound.CurrentBlock(ctx, compound.SecondsPerBlock, s.config.App.Genesis)
	if e != nil {
		return 0, e
	}
	return current, nil
}

func (s *service) GetBlock(ctx context.Context, t time.Time) (int64, error) {
	block, e := compound.GetBlockByTime(ctx, compound.SecondsPerBlock, s.config.App.Genesis, t)
	if e != nil {
		return 0, e
	}
	return block, nil
}
