// Copyright (C) 2017 Miquel Sabaté Solà <msabate@suse.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strings"
	"text/tabwriter"
	"unicode"
	"unicode/utf8"
)

func header(elements interface{}, join bool) string {
	t := reflect.TypeOf(elements).Elem()

	header := ""
	for i := 0; i < t.NumField(); i++ {
		header += t.Field(i).Name + "\t"
	}
	return strings.TrimRight(header, "\t")
}

func tabifyStruct(element interface{}) string {
	str := ""

	v := reflect.ValueOf(element)
	for i := 0; i < v.NumField(); i++ {
		str += fmt.Sprintf("%v\t", v.Field(i))
	}
	return str
}

func printList(elements interface{}) {
	w := tabwriter.NewWriter(os.Stdout, 0, 4, 4, ' ', tabwriter.TabIndent)
	fmt.Fprintln(w, header(elements, true))

	v := reflect.ValueOf(elements)
	for i := 0; i < v.Len(); i++ {
		fmt.Fprintln(w, v.Index(i))
	}
	w.Flush()
}

func printValidate(data Validate) error {
	if data.Valid {
		if !globalConfig.quiet {
			fmt.Println("Valid")
		}
		return nil
	}
	return messagesError(data.Messages)
}

func messagesError(messages map[string][]string) error {
	str := ""
	for k, list := range messages {
		str += k + ":\n"
		for _, e := range list {
			r, n := utf8.DecodeRuneInString(e)
			str += "  - " + string(unicode.ToUpper(r)) + e[n:] + "\n"
		}
	}
	return errors.New(strings.TrimRight(str, "\n"))
}

func prettyPrint(kind int, body []byte, single bool) error {
	var bodyStr string

	if len(body) == 0 || string(body) == "null" {
		return nil
	}

	if single {
		bodyStr = "[" + string(body) + "]"
	} else {
		bodyStr = string(body)
	}

	// TODO: generilize
	switch kind {
	case kindUser:
		data := Users{}
		err := json.Unmarshal([]byte(bodyStr), &data)
		if err != nil {
			return err
		}
		printList(data)
	case kindApplicationToken:
		data := ApplicationTokens{}
		err := json.Unmarshal([]byte(bodyStr), &data)
		if err != nil {
			return err
		}
		printList(data)
	case kindPlainToken:
		data := PlainTokens{}
		err := json.Unmarshal([]byte(bodyStr), &data)
		if err != nil {
			return err
		}
		printList(data)
	case kindTeam:
		data := Teams{}
		err := json.Unmarshal([]byte(bodyStr), &data)
		if err != nil {
			return err
		}
		printList(data)
	case kindNamespace:
		data := Namespaces{}
		err := json.Unmarshal([]byte(bodyStr), &data)
		if err != nil {
			return err
		}
		printList(data)
	case kindTag:
		data := Tags{}
		err := json.Unmarshal([]byte(bodyStr), &data)
		if err != nil {
			return err
		}
		printList(data)
	case kindRegistry:
		data := Registries{}
		err := json.Unmarshal([]byte(bodyStr), &data)
		if err != nil {
			return err
		}
		printList(data)
	case kindRepository:
		data := Repositories{}
		err := json.Unmarshal([]byte(bodyStr), &data)
		if err != nil {
			return err
		}
		printList(data)
	case kindValidate:
		data := Validate{}
		err := json.Unmarshal([]byte(bodyStr), &data)
		if err != nil {
			return err
		}
		return printValidate(data)
	default:
	}
	return nil
}

func printAndQuit(response *http.Response, kind int, single bool) error {
	body, _ := ioutil.ReadAll(response.Body)

	switch globalConfig.format {
	case jsonFmt:
		fmt.Print(string(body))
	default:
		return prettyPrint(kind, body, single)
	}
	return nil
}
