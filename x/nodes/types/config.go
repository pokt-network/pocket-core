package types

var ValidatorCacheSize int64 = 10000

func InitConfig(validatorCacheSize int64) {
	ValidatorCacheSize = validatorCacheSize
}
