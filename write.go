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
	"strings"
)

func WriteRouter(path, appname string) {
	path = path + "/routers/commentsRouter_controllers.go"

	var f *os.File
	var err error
	if CheckFileIsExist(path) {
		f, err = os.OpenFile(path, os.O_TRUNC|os.O_RDWR, 0666)
	} else {
		f, err = os.Create(path)
	}
	Check(err)
	write(f, appname)
}

func write(f *os.File, appname string) {
	var content = `package routers

import (
	"net/http"
	"strings"%s

	"github.com/gorilla/sessions"
	"github.com/dotqin/blfgo"
	"%s/controllers"%s
)

func (r *Router) Route(_w http.ResponseWriter, _req *http.Request, _s *sessions.Session) (re string, tp int, data map[string]interface{}) {

	var url = strings.SplitN(_req.RequestURI, "?", 2)[0]

	switch {%s
	default:
		return re, 404, nil
	}
	return re, tp, data
}`

	var useStrconv, useReflect bool
	var importModels string

	var caseContent string
	for _, v := range Routers {
		var ctx string
		m := "controllers." + v.ControllerName
		if v.Receiver != "" {
			m = "c." + v.ControllerName
			ctx = fmt.Sprintf(`
		c := &controllers.%s{blfgo.Controller{"%s", "%s", %d, %t, %s, _w, _req, _s, nil}}
		if !blfgo.Intercept(&c.Controller) {
			return re, 9, nil
		}
		c.Prepare()`, v.Receiver, v.Url, strings.ToUpper(v.Method), v.Flag, v.Csrf, v.JoinCustomArgs())
		}

		var argSatatements, args string

		ln := len(v.ControllerArgTypes)
		if ln > 0 {
			for i, n := range v.ControllerArgNames {
				switch v.ControllerArgTypes[i] {
				case "string":
					argSatatements += fmt.Sprintf(`
		%s := strings.Join(_req.Form["%s"], "")`, n, n)
				case "[]string":
					argSatatements += fmt.Sprintf(`
		%s := _req.Form["%s"]`, n, n)
				case "int":
					useStrconv = true
					argSatatements += fmt.Sprintf(`
		%s, _ := strconv.Atoi(strings.Join(_req.Form["%s"], ""))`, n, n)
				case "interface{}":
					argSatatements += fmt.Sprintf(`
		%s := _req.Form["%s"]`, n, n)
				default:
					df := fmt.Sprintf(`
		var %s %s`, n, v.ControllerArgTypes[i])
					if importModels == "" && strings.Index(v.ControllerArgTypes[i], "models.") > -1 {
						useReflect = true
						importModels = fmt.Sprintf(`
	"%s/models"`, appname)
						df = fmt.Sprintf(`
		var %s %s
		if reflect.ValueOf(%s).Kind() == reflect.Struct {
			%s = %s{}
			blfgo.ParseForm(_req.Form, &%s)
		}`, n, v.ControllerArgTypes[i], n, n, v.ControllerArgTypes[i], n)
					}
					argSatatements += df
				}
				args += n
				if i != ln-1 {
					args += ", "
				}
			}
			ctx = fmt.Sprintf(`
		_req.ParseForm()%s%s`, ctx, argSatatements)
		}

		caseContent += fmt.Sprintf(`
    case _req.Method == "%s" && url == "%s":%s
        re, tp = %s(%s)
		data = c.Data`, strings.ToUpper(v.Method), v.Url, ctx, m, args)
	}

	if PackPre != "" {
		appname = PackPre + "/" + appname
		if importModels != "" {
			importModels = fmt.Sprintf(`
	"%s/models"`, appname)
		}
	}

	var usePacks string
	if useStrconv {
		usePacks += `
	"strconv"`
	}
	if useReflect {
		usePacks += `
	"reflect"`
	}

	_, err := io.WriteString(f, fmt.Sprintf(content, usePacks, appname, importModels, caseContent))
	f.Close()
	Check(err)
}
