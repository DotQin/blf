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

package log

import (
	"fmt"
	"os"

	"github.com/op/go-logging"
)

var (
	log    = logging.MustGetLogger("blf")
	format = logging.MustStringFormatter(
		"BLF %{time:2006-01-02 15:04:05.000} %{color}[%{level:.4s}]%{color:reset} %{message}",
	)
)

func Debug(args ...interface{}) { // lv5
	log.Debug(args...)
}

func Debugf(message string, args ...interface{}) {
	log.Debugf(message, args...)
}

func Info(args ...interface{}) { // lv4
	log.Info(args...)
}

func Infof(message string, args ...interface{}) {
	log.Infof(message, args...)
}

func Notice(args ...interface{}) { // lv3
	log.Notice(args...)
}

func Noticef(message string, args ...interface{}) {
	log.Noticef(message, args...)
}

func Warning(args ...interface{}) { // lv2
	log.Warning(args...)
}

func Warningf(message string, args ...interface{}) {
	log.Warningf(message, args...)
}

func Error(args ...interface{}) { // lv1
	log.Error(args...)
}

func Errorf(message string, args ...interface{}) {
	log.Errorf(message, args...)
}

func GetLevel() string {
	return fmt.Sprintf("%v", logging.GetLevel("blf"))
}

func init() {
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	formatter := logging.NewBackendFormatter(backend, format)
	log.SetBackend(logging.AddModuleLevel(formatter))
	//logging.SetBackend(formatter) //全局
}
