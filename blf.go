package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/dotqin/blf/log"
	"github.com/op/go-logging"
)

var AppName, AppPath, PackPre, GOPATH string

func main() {

	var help = `
Usage:
    blf command [arguments]

The commands are:
    run    run the app and start a Web server for development
    new    create app

Use "blf help [command]" for more information about a command.
`
	AppPath = getCurrentDirectory()
	GOPATH = os.Getenv("GOPATH")

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
						default:
							AppPath += "/" + v
							AppPath = TrimRight(AppPath, "/")
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
	log.Notice("Running ...")

	conf := &Config{}
	conf.InitConfig(AppPath + "/conf/app.conf")
	AppName = conf.Read("default", "appname")

	log.Debug("AppPath\t:", AppPath)
	log.Debug("GOPATH\t:", GOPATH)

	PackPre = conf.Read("default", "packpre")
	if PackPre == "" {
		PackPre = TrimRight(TrimLeft(AppPath, GOPATH+"/src/"), "/"+AppName)
	}
	log.Debug("PackPre\t:", PackPre)

	Scan(AppPath, 0)

	WriteRouter(AppPath, AppName, PackPre)

	RunCmd(AppPath)

	wait()
}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Debug(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

func TrimLeft(s, t string) string {
	if n := strings.Index(s, t); n > -1 {
		return s[n+len(t):]
	}
	return s
}

func TrimRight(s, t string) string {
	if n := strings.LastIndex(s, t); n > -1 {
		return s[:n]
	}
	return s
}

func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func Check(e error) {
	if e != nil {
		panic(e)
	}
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
