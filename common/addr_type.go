package common

type Addr struct {
	// 接口名称
	Name string `yaml:"name"`
	// 接口地址，需要带掩码
	Address string `yaml:"address"`
	// 点到点模式（Point-to-Point）
	PtpMode bool `yaml:"ptp_mode"`
	// 点到点模式下的对端前缀
	PeerPrefix string `yaml:"peer_prefix"`
}
