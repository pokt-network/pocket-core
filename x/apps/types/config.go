package types

var ApplicationCacheSize int64 = 5

func InitConfig(applicationCacheSize int64) {
	ApplicationCacheSize = applicationCacheSize
}
