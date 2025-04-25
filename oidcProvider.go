package core

import (
	"context"
	"github.com/Nerzal/gocloak/v13"
	"github.com/ReneKroon/ttlcache"
	"time"
)

var users = newCache()

func newCache() *ttlcache.Cache {
	cache := ttlcache.NewCache()
	cache.SetTTL(time.Nanosecond * 250)
	return cache
}

type oidcProvider struct {
	*gocloak.GoCloak
	fetcher dataFetcher
}

func (p *oidcProvider) GetUser(token string) (*UserInfo, error) {
	if user, ok := users.Get(token); ok {
		logger.Debug("hit cache")
		return user.(*UserInfo), nil
	}
	user, err := p.fetcher.GetUserInfo(context.Background(), token, libConfig.KeyCloak.RealmName)
	logger.Debug("fetched user data", "userIsFound", user != nil, "error", err)
	if user != nil {
		userInfo := UserInfo{user}
		users.Set(token, &userInfo)
		return &userInfo, nil
	}
	return nil, err
}

func getOidcProvider(f ...dataFetcher) *oidcProvider {
	oidcService := gocloak.NewClient(libConfig.KeyCloak.Url, func(cloak *gocloak.GoCloak) {
		// additional hooks
	})

	if len(f) == 0 {
		f = []dataFetcher{oidcService}
	}
	oidc := &oidcProvider{oidcService, f[0]}
	return oidc
}
