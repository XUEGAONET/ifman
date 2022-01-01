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
	"github.com/XUEGAONET/ifman/common"
	"log"
)

func main() {
	// config variable
	var (
		configFile string
	)

	flag.StringVar(&configFile, "config", "config.yaml", "yaml config path")
	flag.Parse()

	// init dynamic core config
	err := initGlobalConfig(configFile)
	if err != nil {
		log.Fatalln(err)
	}

	// get newest global config and init logger
	conf := getGlobalConfig()
	err = initLogger(&conf.Logger)
	if err != nil {
		log.Fatalln(err)
	}

	// listen socket
	_, err = NewGrpcServer(uint16(common.GrpcPort))
	if err != nil {
		log.Fatalln(err)
	}

	// start core service
	err = startService()
	if err != nil {
		log.Fatalln(err)
	}

}
