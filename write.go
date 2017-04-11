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

	"%s/controllers"%s
)

func (r *Router) Route(req *http.Request) (re string, tp int) {

	var url = strings.SplitN(req.RequestURI, "?", 2)[0]

	switch {%s
	}
	return re, tp
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
		c := &controllers.%s{}`, v.Receiver)
		}

		var argSatatements, args string

		ln := len(v.ControllerArgTypes)
		if ln > 0 {
			for i, n := range v.ControllerArgNames {
				switch v.ControllerArgTypes[i] {
				case "string":
					argSatatements += fmt.Sprintf(`
		%s := strings.Join(req.Form["%s"], "")`, n, n)
				case "[]string":
					argSatatements += fmt.Sprintf(`
		%s := req.Form["%s"]`, n, n)
				case "int":
					useStrconv = true
					argSatatements += fmt.Sprintf(`
		%s, _ := strconv.Atoi(strings.Join(req.Form["%s"], ""))`, n, n)
				case "interface{}":
					argSatatements += fmt.Sprintf(`
		%s := req.Form["%s"]`, n, n)
				default:
					df := fmt.Sprintf(`
		var %s %s`, n, v.ControllerArgTypes[i])
					if importModels == "" && strings.Index(v.ControllerArgTypes[i], "models.") > -1 {
						useReflect = true
						importModels = fmt.Sprintf(`
	"github.com/dotqin/blfgo"
	"%s/models"`, appname)
						df = fmt.Sprintf(`
		var %s %s
		if reflect.ValueOf(%s).Kind() == reflect.Struct {
			%s = %s{}
			blfgo.ParseForm(req.Form, &%s)
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
		req.ParseForm()%s%s`, argSatatements, ctx)
		}

		caseContent += fmt.Sprintf(`
    case req.Method == "%s" && url == "%s":%s
        re, tp = %s(%s)`, strings.ToUpper(v.Method), v.Url, ctx, m, args)
	}

	if PackPre != "" {
		appname = PackPre + "/" + appname
		if importModels != "" {
			importModels = fmt.Sprintf(`
	"github.com/dotqin/blfgo"
	"%s/models"`, appname)
		}
	}

	var usePacks string
	if useStrconv {
		usePacks += `
	strconv`
	}
	if useReflect {
		usePacks += `
	reflect`
	}

	_, err := io.WriteString(f, fmt.Sprintf(content, usePacks, appname, importModels, caseContent))
	f.Close()
	Check(err)
}
