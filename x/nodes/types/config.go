package types

var ValidatorCacheSize int64 = 5

func InitConfig(validatorCacheSize int64) {
	ValidatorCacheSize = validatorCacheSize
}
