package main

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"sync"
)

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
	Mode     string `yaml:"mode"`
	Level    string `yaml:"level"`
	SyslogOn bool   `yaml:"syslog_on"`
	Single   Single `yaml:"single"`
	Rotate   Rotate `yaml:"rotate"`
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

func parseLocalConfig(path string) (*Config, error) {
	b, err := readFile(path)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	c := Config{}
	err = yaml.Unmarshal(b, &c)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &c, nil
}

func readFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer f.Close()

	return ioutil.ReadAll(f)
}

type DynamicConfig struct {
	conf *Config
	path string
	lock sync.RWMutex
}

var coreConfig DynamicConfig

func getCoreConfig() *Config {
	coreConfig.lock.RLock()
	defer coreConfig.lock.RUnlock()

	return coreConfig.conf
}

func refreshCoreConfig() error {
	coreConfig.lock.Lock()
	defer coreConfig.lock.Unlock()

	c, err := parseLocalConfig(coreConfig.path)
	if err != nil {
		return errors.WithStack(err)
	}

	coreConfig.conf = c

	logrus.Infof("core config refreshed")

	return nil
}

func initCoreConfig(path string) error {
	coreConfig = DynamicConfig{
		conf: nil,
		path: path,
		lock: sync.RWMutex{},
	}

	c, err := parseLocalConfig(path)
	if err != nil {
		return errors.WithStack(err)
	}

	coreConfig.lock.Lock()
	defer coreConfig.lock.Unlock()
	coreConfig.conf = c

	return nil
}
