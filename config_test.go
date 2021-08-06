package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestParse(t *testing.T) {
	conf, err := parseLocalConfig("config.yaml")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", conf)

	inf := conf.Interface[0]
	i := inf.(map[string]interface{})
	typ := i["type"]
	if typ == "ipip" {
		b, err := yaml.Marshal(inf)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\n", b)

		ipip := IPTun{
			BaseLink: BaseLink{},
			Ttl:      0,
			Tos:      0,
			LocalIP:  "",
			RemoteIP: "",
		}
		err = yaml.Unmarshal(b, &ipip)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%+v\n", ipip)
	}
}
