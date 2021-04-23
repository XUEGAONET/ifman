package ethtool

type Offload struct {
	GRO bool
	GSO bool
	TSO bool
	LRO bool
}

func GetAttr() *Offload {
	return &Offload{
		GRO: false,
		GSO: false,
		TSO: false,
		LRO: false,
	}
}
