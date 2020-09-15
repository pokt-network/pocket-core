package types

//type Upgrade struct {
//	Height  int64  `json:"Height"`
//	Version string `json:"Version"`
//}

func NewUpgrade(height int64, version string) Upgrade {
	return Upgrade{
		Height:  height,
		Version: version,
	}
}

func (u Upgrade) UpgradeHeight() int64 {
	return u.Height
}

func (u Upgrade) UpgradeVersion() string {
	return u.Version
}
