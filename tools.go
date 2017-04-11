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
	"os"
	"path/filepath"
	"strings"

	"github.com/dotqin/blf/log"
)

// getCurrentDirectory 获取当前工作路径
func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Error(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

// TrimLeft 从s左侧去掉t
func TrimLeft(s, t string) string {
	if n := strings.Index(s, t); n > -1 {
		return s[n+len(t):]
	}
	return s
}

// TrimRight 从s右侧去掉t
func TrimRight(s, t string) string {
	if n := strings.LastIndex(s, t); n > -1 {
		return s[:n]
	}
	return s
}

// CheckFileIsExist 检查文件或文件夹是否存在
func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

// Check 检查处理错误
func Check(e error) {
	if e != nil {
		panic(e)
	}
}
