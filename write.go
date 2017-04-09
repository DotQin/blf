package main

import (
	"fmt"
	"io"
	"os"
)

func WriteRouter(path, appname, packpre string) {
	path = path + "/routers/commentsRouter_controllers.go"

	var f *os.File
	var err error
	if CheckFileIsExist(path) {
		f, err = os.OpenFile(path, os.O_TRUNC|os.O_RDWR, 0666)
	} else {
		f, err = os.Create(path)
	}
	Check(err)
	write(f, appname, packpre)
}

func write(f *os.File, appname, packpre string) {
	var content = `package routers

import (
	"%s/%s/controllers"
)

func (r *Router) Route(url string) string {

	var content string

	switch url {%s
	}
	return content
}` // TODO switch url and method

	var caseContent string
	for _, v := range Routers {
		var c string
		m := "controllers." + v.ControllerName
		if v.Receiver != "" {
			m = "c." + v.ControllerName
			c = fmt.Sprintf(`
		c := &controllers.%s{}`, v.Receiver)
		}
		caseContent += fmt.Sprintf(`
    case "%s":%s
        content = %s()
`,
			v.Url, c, m)
	}

	_, err := io.WriteString(f, fmt.Sprintf(content, packpre, appname, caseContent))
	Check(err)
}
