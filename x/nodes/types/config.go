package types

var ValidatorCacheSize int64

func InitConfig(validatorCacheSize int64) {
	ValidatorCacheSize = validatorCacheSize
}
