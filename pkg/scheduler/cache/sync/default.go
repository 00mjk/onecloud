package sync

import (
	u "yunion.io/x/pkg/utils"

	o "yunion.io/x/onecloud/cmd/scheduler/options"
	"yunion.io/x/onecloud/pkg/scheduler/cache"
	networks_db "yunion.io/x/onecloud/pkg/scheduler/cache/sync/networks/db"
)

const (
	CacheKind = "SyncCache"

	//GlanceSyncCache       = "Glance"
	NetworkSyncCache      = "Network"
	NetworksDataSyncCache = "NetworkData"
)

func getUpdate(d []interface{}) ([]string, error) {
	return nil, nil
}

func defaultSyncItems() []cache.CachedItem {
	return []cache.CachedItem{
		//newGlanceCache(),
		//newNetworkCache(),
		//newNetworksDataCache(),
	}
}

func noneUpdate(id []string) ([]interface{}, error) {
	return []interface{}{}, nil
}

/*
func newGlanceCache() cache.CachedItem {
	item := new(syncItem)

	item.CachedItem = cache.NewCacheItem(
		GlanceSyncCache,
		viper.GetDuration("cache.glance_cache.ttl"),
		viper.GetDuration("cache.glance_cache.period"),
		imageUUIDKey,
		noneUpdate,
		loadImages,
		getUpdate,
	)
	return item
}*/

func newNetworkCache() cache.CachedItem {
	item := new(syncItem)

	// data from db
	item.CachedItem = cache.NewCacheItem(
		NetworkSyncCache,
		u.ToDuration(o.GetOptions().NetworkCacheTTL),
		u.ToDuration(o.GetOptions().NetworkCachePeriod),
		networks_db.BuilderCacheKey,
		networks_db.UpdateNetworkDescBuilder,
		networks_db.LoadNetworkDescBuilder,
		getUpdate,
	)

	return item
}

func newNetworksDataCache() cache.CachedItem {
	item := new(syncItem)

	item.CachedItem = cache.NewCacheItem(
		NetworksDataSyncCache,
		u.ToDuration(o.GetOptions().NetworkCacheTTL),
		u.ToDuration(o.GetOptions().NetworkCachePeriod),
		BuilderNetworkCacheKey,
		updateNetworksBuilder,
		loadNetworksBuilder,
		getUpdate,
	)
	return item
}
