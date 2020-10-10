package session

import (
	"context"

	"compound/core"

	"github.com/bluele/gcache"
)

type cacheSession struct {
	core.Session
	tokens gcache.Cache
}

func (s *cacheSession) Login(ctx context.Context, accessToken string) (*core.User, error) {
	if user, err := s.tokens.Get(accessToken); err == nil {
		return user.(*core.User), nil
	}

	user, err := s.Session.Login(ctx, accessToken)
	if err != nil {
		return nil, err
	}

	_ = s.tokens.Set(accessToken, user)
	return user, nil
}
