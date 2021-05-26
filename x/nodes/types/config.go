package types

var ValidatorCacheSize int64 = 100000

func InitConfig(validatorCacheSize int64) {
	ValidatorCacheSize = validatorCacheSize
}
