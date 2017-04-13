// Copyright 2017 blf Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/dotqin/blf/log"
	"github.com/fatih/color"
	"github.com/op/go-logging"
)

var AppName, AppPath, PackPre, GOPATH, GOROOT string

func main() {

	AppPath = GetCurrentDirectory()
	GOPATH = os.Getenv("GOPATH")
	GOPATH = strings.Replace(GOPATH, "\\", "/", -1)
	GOROOT = os.Getenv("GOROOT")
	GOROOT = strings.Replace(GOROOT, "\\", "/", -1)

	var help = `
Usage:
    blf command [arguments]

The commands are:
    run    run the app and start a Web server for development
    new    create app

Use "blf help [command]" for more information about a command.
`

	var logo = `
   ____    _       _____
  | __ )  | |     |  ___|
  |  _ \  | |     | |_
  | |_) | | |___  |  _|
  |____/  |_____| |_|
	`

	if len(os.Args) > 1 {
		for _, v := range os.Args[1:] {
			switch v {
			case "help":
				if len(os.Args) > 2 {
					for _, v := range os.Args[2:] {
						switch v {
						case "run":
							fmt.Println("usage: blf run [-nolog]")
						case "new":
							fmt.Println("usage: blf new [appname]")
						default:
							fmt.Printf(`Unknown help topic "%s".  Run 'blf help'.`, v)
						}
					}
				} else {
					fmt.Println(help)
				}
				return
			case "run":
				if len(os.Args) > 2 {
					for _, v := range os.Args[2:] {
						switch v {
						case "-nolog":
							logging.SetLevel(0, "blf")
						}
					}
				}
			case "new":
				if len(os.Args) > 2 {
					AppName = strings.TrimSpace(os.Args[2])
					NewApp(AppPath, AppName)
				} else {
					log.Warning("must enter a name")
				}
				return
			}
		}
	} else {
		fmt.Println(help)
		return
	}

	conf := &Config{}
	conf.InitConfig(AppPath + "/conf/app.conf")
	AppName = conf.Read("default", "appname")

	PackPre = conf.Read("default", "packpre")
	if PackPre == "" {
		PackPre = TrimLeft(TrimRight(TrimLeft(AppPath, GOPATH+"/src"), "/"+AppName), "/")
	}

	log.Debug("AppName\t:", AppName)
	log.Debug("AppPath\t:", AppPath)
	log.Debug("PackPre\t:", PackPre)
	log.Debug("GOPATH\t:", GOPATH)
	log.Debug("GOROOT\t:", GOROOT)

	Scan(AppPath, 0)

	WriteRouter(AppPath, AppName)

	RunCmd(AppPath)

	color.Green(logo)
	log.Notice("Running ...")

	wait()
}

func wait() {
	sig := make(chan os.Signal, 1)

	signal.Notify(sig, os.Interrupt)
	signal.Notify(sig, syscall.SIGTERM)

	for {
		<-sig
		fmt.Println("")
		log.Warning("Stopping ...")
		Kill()
		os.Exit(1)
	}
}
