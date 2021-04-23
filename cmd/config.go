package main

type Logger struct {
	Mode   string `yaml:"mode"`
	Level  string `yaml:"level"`
	Single Single `yaml:"single"`
	Rotate Rotate `yaml:"rotate"`
}

type Single struct {
	Path string `yaml:"path"`
	Perm uint16 `yaml:"perm"`
}

type Rotate struct {
	Dir       string `yaml:"dir"`
	Age       uint32 `yaml:"age"`
	CutPeriod uint32 `yaml:"cut_period"`
}

type Interface struct {
	Name        string                 `yaml:"name"`
	Description string                 `yaml:"description"`
	Status      string                 `yaml:"status"`
	Type        string                 `yaml:"type"`
	Config      map[string]interface{} `yaml:"config"`
	Tso         bool                   `yaml:"tso"`
	Lro         bool                   `yaml:"lro"`
	Gso         bool                   `yaml:"gso"`
	Gro         bool                   `yaml:"gro"`
	RpFilter    string                 `yaml:"rp_filter"`
	Address     []string               `yaml:"address"`
}

type Route struct {
	Dst     string `yaml:"dst"`
	Gw      string `yaml:"gw"`
	TableId uint8  `yaml:"table_id"`
}

type Verify struct {
	CheckPeriod uint32 `yaml:"check_period"`
}

type RootConfig struct {
	Logger    Logger      `yaml:"logger"`
	Interface []Interface `yaml:"interface"`
	Route     []Route     `yaml:"route"`
	Verify    Verify      `yaml:"verify"`
}
