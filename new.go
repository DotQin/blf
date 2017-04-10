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
	"io"
	"os"

	"github.com/dotqin/blf/log"
)

func NewApp(path, appname string) {
	apppath := path + "/" + appname
	if CheckFileIsExist(apppath) {
		log.Errorf("%s already exists", appname)
	} else {
		err := os.Mkdir(apppath, 0777)
		if err != nil {
			log.Errorf("%s", err)
		} else {
			os.Mkdir(apppath+"/conf", 0777)
			os.Mkdir(apppath+"/controllers", 0777)
			os.Mkdir(apppath+"/routers", 0777)
			os.Mkdir(apppath+"/static", 0777)
			os.Mkdir(apppath+"/views", 0777)
			createAppConf(path, appname)
			createTestController(apppath)
			createRouter(apppath)
			createMain(apppath, appname)
			log.Debug("App Create Success!")
			log.Debugf("Run : cd %s && blf run", appname)
		}
	}
}

func createAppConf(path, appname string) {
	var content = `[default]

# 应用名称
appname = %s

# 监听端口
httpport = 8080

# 应用的模式，默认是 dev，为开发模式
runmode = dev

# 开启session
sessionon = true

# 包前缀
# packpre = %s
`
	PackPre = TrimLeft(path, GOPATH+"/src/")
	createFile(path+"/"+appname+"/conf/app.conf", fmt.Sprintf(content, appname, PackPre))
}

func createTestController(apppath string) {
	var content = `package controllers

import "github.com/dotqin/blfgo"

type TestController struct {
	blfgo.Controller
}

// @router(url=/, method=get, flag=100, csrf=false)
func (c *TestController) Home() (content string) {
	content = "Welcome to Blfgo !"
	return content
}

// @router(url=/test1, method=get, flag=101, csrf=false, test=1, response_body=true)
func (c *TestController) TestA() (content string) {
	content = "from test11"
	return content
}

// @router(url=/test2, method=post, flag=102, csrf=true)
func (c *TestController) TestB() (content string) {
	content = "from test22"
	return content
}
`
	createFile(apppath+"/controllers/test_controller.go", content)
}

func createRouter(apppath string) {
	var content = `package routers

import "github.com/dotqin/blfgo"

type Router struct {
}

func init() {
	blfgo.Router = &Router{}
}
`
	createFile(apppath+"/routers/router.go", content)
}

func createMain(apppath, appname string) {
	var content = `package main

import (
	"log"
	"net/http"

	"github.com/dotqin/blfgo"
	_ "%s/%s/routers"
)

func main() {
	err := http.ListenAndServe(":8080", &blfgo.BlfHandler{})
	if err != nil {
		log.Fatal("Blfgo :", err)
	}
}`
	createFile(apppath+"/main.go", fmt.Sprintf(content, PackPre, appname))
}

func createFile(path, content string) {
	var f *os.File
	var err error
	f, err = os.Create(path)
	Check(err)
	_, err = io.WriteString(f, content)
	Check(err)
}
