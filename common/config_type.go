package common

type Config struct {
	Logger    Logger        `yaml:"logger"`
	Interface []interface{} `yaml:"interface"`
	Addr      []Addr        `yaml:"addr"`
	RpFilter  []RpFilter    `yaml:"rp_filter"`
	Learning  []Learning    `yaml:"learning"`
	Common    Common        `yaml:"common"`
}

type Common struct {
	CheckPeriodSec uint16 `yaml:"check_period_sec"`
}

type Logger struct {
	Mode   string `yaml:"mode"`
	Level  string `yaml:"level"`
	Single Single `yaml:"single"`
	Rotate Rotate `yaml:"rotate"`
}

type Rotate struct {
	Dir       string `yaml:"dir"`
	MaxAgeSec uint32 `yaml:"max_age_sec"`
	PeriodSec uint32 `yaml:"period_sec"`
}

type Single struct {
	Path   string `yaml:"path"`
	Permit int    `yaml:"permit"`
}
