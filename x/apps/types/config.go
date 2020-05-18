package types

var ApplicationCacheSize int64

func InitConfig(applicationCacheSize int64) {
	ApplicationCacheSize = applicationCacheSize
}
