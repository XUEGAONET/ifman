// Copyright 2021 The ifman authors https://github.com/XUEGAONET/ifman
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	// global version
	version string = "v2.0"
)

// global config
const (
	ctlPort = 11073
)

func loadedModules() string {
	// global modules
	modules := []string{
		"core",
		"test",
		"key",
		"reload",
	}

	var res string
	for i, _ := range modules {
		res += modules[i]

		if i+1 == len(modules) {
			break
		} else {
			res += ", "
		}
	}

	return res
}

func printVersion() {
	var banner string
	banner += fmt.Sprintf("XUEGAONET  https://github.com/XUEGAONET\n")
	banner += fmt.Sprintf("* Component: ifman, Interface Manager\n")
	banner += fmt.Sprintf("* Version: %s\n", version)
	banner += fmt.Sprintf("* Loaded modules: %s\n", loadedModules())

	fmt.Println(banner)
	fmt.Println("Usage: ")
	flag.PrintDefaults()

	os.Exit(0)
}

func main() {
	// config variable
	var (
		configFile string
		module     string
	)

	flag.Usage = printVersion
	flag.StringVar(&configFile, "config", "config.yaml", "yaml config path")
	flag.StringVar(&module, "module", "", "which module you want to use")
	flag.Parse()

	switch module {
	case "test":
		_, err := parseLocalConfig(configFile)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("OK")
		}
	case "core":
		// init dynamic core config
		err := initCoreConfig(configFile)
		if err != nil {
			panic(err)
		}

		// get newest core config and init logger
		lc := getCoreConfig()
		err = initLogger(&lc.Logger)
		if err != nil {
			panic(err)
		}

		// listen socket
		err = ListenCtl(ctlPort)
		if err != nil {
			panic(err)
		}

		// start core service
		err = startCoreService()
		if err != nil {
			panic(err)
		}
	case "key":
		generateWireGuardKeyChain()
	case "reload":
		err := SendReload(ctlPort)
		if err != nil {
			panic(err)
		}
		fmt.Println("send reload succeed, please care about the log ifman output")
	default:
		fmt.Printf("Please specify the module you want to use.\n")
	}
}
