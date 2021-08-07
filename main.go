/**
 * Copyright (c) 2021 Harris <ic0xgkk@gmail.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of
 * this software and associated documentation files (the "Software"), to deal in
 * the Software without restriction, including without limitation the rights to
 * use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
 * the Software, and to permit persons to whom the Software is furnished to do so,
 * subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
 * FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 * COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
 * IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 * CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

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
	sockFile = "/var/run/ifman.socket"
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
		err = ListenSocket(sockFile)
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
		err := SendReload(sockFile)
		if err != nil {
			panic(err)
		}
		fmt.Println("send reload succeed, please care about the log ifman output")
	default:
		fmt.Printf("Please specify the module you want to use.\n")
	}
}
