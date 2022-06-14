package types

var ValidatorCacheSize int64 = 50000

func InitConfig(validatorCacheSize int64) {
	ValidatorCacheSize = validatorCacheSize
}
