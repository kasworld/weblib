// Copyright 2015,2016,2017,2018,2019 SeukWon Kang (kasworld@gmail.com)
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package weblib

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func WebFormBegin(title string, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w,
		`<html>
	<head>
	<title>%v</title>
	</head>
	<body>`, title,
	)
}

func WebFormEnd(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w,
		`
		</body>
		</html>`)
}

func SetFresh(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return fmt.Errorf("form err : %v", err)
	}
	refreshStr := r.Form.Get("refresh")
	if refreshStr == "" {
		return nil
	}
	if _, err := strconv.ParseInt(refreshStr, 0, 64); err != nil {
		return fmt.Errorf("parse err : %v", err)
	}
	w.Header().Add("Refresh", refreshStr)
	return nil
}

func GetIntByName(
	name string, errorval int, w http.ResponseWriter, r *http.Request) int {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "parse error", 405)
		return errorval
	}

	valuestring := r.Form.Get(name)
	valuestring = strings.TrimSpace(valuestring)
	if valuestring == "" {
		return errorval
	}
	value, err := strconv.ParseInt(valuestring, 0, 64)
	if err != nil {
		http.Error(w, "parse error", 405)
		return errorval
	}

	return int(value)
}

func GetFloat64ByName(
	name string, errorval float64, w http.ResponseWriter, r *http.Request) float64 {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "parse error", 405)
		return errorval
	}

	valuestring := r.Form.Get(name)
	valuestring = strings.TrimSpace(valuestring)
	if valuestring == "" {
		return errorval
	}
	value, err := strconv.ParseFloat(valuestring, 64)
	if err != nil {
		http.Error(w, "parse error", 405)
		return errorval
	}

	return float64(value)
}

func GetStringByName(
	name string, errorval string, w http.ResponseWriter, r *http.Request) string {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "parse error", 405)
		return errorval
	}

	valuestring := r.Form.Get(name)
	valuestring = strings.TrimSpace(valuestring)
	if valuestring == "" {
		return errorval
	}
	return valuestring
}

func GetPage(w http.ResponseWriter, r *http.Request) int {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "parse error", 405)
		return 0
	}
	pagestr := r.Form.Get("page")
	if pagestr == "" {
		return 0
	}
	page, err := strconv.ParseInt(pagestr, 0, 64)
	if err != nil {
		http.Error(w, "parse error", 405)
		return 0
	}
	return int(page)
}

func WebLog(logFilename string, size int64, w http.ResponseWriter, r *http.Request) error {
	SetFresh(w, r)
	filter := GetStringByName("filter", "", w, r)

	file, err := os.OpenFile(logFilename, os.O_RDONLY, 0660)
	if err != nil {
		return fmt.Errorf("Can't Open Log File (%s), %v", logFilename, err)
	}

	fi, fierr := file.Stat()
	if fierr != nil {
		return fmt.Errorf("Can't Open Log File (%s) info, %v", logFilename, fierr)
	}

	if fi.Size() < size {
		size = fi.Size()
	}

	_, seekerr := file.Seek(-size, os.SEEK_END)
	if seekerr != nil {
		return fmt.Errorf("Seek Failed (%s), %v", logFilename, seekerr)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := scanner.Text()
		if filter == "" || strings.Contains(s, filter) {
			fmt.Fprintln(w, s)
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("Read Failed (%s), %v", logFilename, err)
	}
	return nil
}
