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
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/dotqin/blf/log"
)

type Router struct {
	Url            string
	Method         string
	Flag           int
	Csrf           bool
	Receiver       string
	ControllerName string
	ResponseBody   bool
	CustomArgs     map[string]interface{}
}

var Routers map[string]Router = make(map[string]Router)

func Scan(path string, level int) {

	dirs, err := ioutil.ReadDir(path)
	if err != nil {
		log.Debug("Scan dir error")
		return
	}
	for _, v := range dirs {
		if level == 0 {
			if v.Name() == "controllers" && v.IsDir() {
				log.Debug("Found Controllers Dir")
				Scan(path+"/controllers", level+1)
			}
		} else {
			if v.IsDir() {
				Scan(path+"/"+v.Name(), level+1)
			} else {
				read(path + "/" + v.Name())
			}
		}
	}
}

func read(path string) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Debugf("Parse file error : %s", path)
		return
	} else {
		content := string(bytes)
		sep := "// @router"
		handle(content, sep)
	}
}

func handle(content, sep string) {
	if strings.Index(content, sep) > -1 {
		array := strings.SplitN(content, sep, 2)
		if len(array) == 2 {
			parse(array[1])
			handle(array[1], sep)
		}
	}
}

func parse(code string) {
	code = strings.TrimSpace(code)
	code = TrimLeft(code, "(")
	array := strings.SplitN(code, ")", 2)
	if len(array) == 2 {
		r := Router{}
		r.CustomArgs = make(map[string]interface{})
		fn := strings.TrimSpace(array[1])
		if strings.Index(fn, "func") == 0 {
			foo := strings.SplitN(fn, "\n", 2)
			if len(foo) == 2 {
				bar := TrimLeft(foo[0], "func")
				bar = TrimRight(bar, "{")
				bar = strings.TrimSpace(bar)
				if len(bar) > 2 {
					var fnName string
					getFuName := func(str string) string {
						array := strings.SplitN(strings.TrimSpace(str), "(", 2)
						if len(array) == 2 {
							return strings.TrimSpace(array[0])
						}
						return ""
					}
					if strings.Index(bar, "(") == 0 {
						bar = TrimLeft(bar, "(")
						barSlice := strings.SplitN(bar, ")", 2)
						if len(barSlice) == 2 {
							recSlice := strings.SplitN(strings.TrimSpace(barSlice[0]), " ", 2)
							if len(recSlice) == 2 {
								r.Receiver = TrimLeft(strings.TrimSpace(recSlice[1]), "*")
							}
							fnName = getFuName(barSlice[1])
						} else {
							return
						}
					} else {
						fnName = getFuName(bar)
					}
					if fnName != "" {
						r.ControllerName = fnName
					} else {
						return
					}
				} else {
					return
				}
			}
		} else {
			return
		}
		args := strings.Split(array[0], ",")
		if len(args) > 0 {
			for _, v := range args {
				v = strings.TrimSpace(v)
				arg := strings.SplitN(v, "=", 2)
				if len(arg) == 2 {
					key := strings.TrimSpace(arg[0])
					val := strings.TrimSpace(arg[1])
					switch key {
					case "url":
						r.Url = val
					case "method":
						r.Method = val
					case "flag":
						flag, err := strconv.Atoi(val)
						if err != nil {
							panic("Flag must be a int value")
						} else {
							r.Flag = flag
						}
					case "csrf":
						switch val {
						case "true":
							r.Csrf = true
						case "false":
							r.Csrf = false
						default:
							panic("Csrf must be true or false")
						}
					case "response_body":
						switch val {
						case "true":
							r.ResponseBody = true
						case "false":
							r.ResponseBody = false
						default:
							panic("ResponseBody must be true or false")
						}
					default:
						r.CustomArgs[key] = val
					}
				}
			}
		}

		if _, ok := Routers[r.Url]; !ok {
			Routers[r.Url] = r
		} else {
			panic("Url Repeat : " + r.Url)
		}
	}
}
