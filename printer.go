// Copyright (C) 2017-2019 Miquel Sabaté Solà <msabate@suse.com>
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

func capitalize(str string) string {
	r, n := utf8.DecodeRuneInString(str)
	return string(unicode.ToUpper(r)) + str[n:]
}

func printHealth(data Health) error {
	success := true
	keys, values := "", ""

	for k, v := range data {
		if !globalConfig.quiet {
			keys += strings.ToUpper(k) + "\t"
			values += capitalize(v.Message) + "\t"
		}

		if !v.Success {
			success = false
			if globalConfig.quiet {
				break
			}
		}
	}

	if !globalConfig.quiet {
		keys = strings.TrimRight(keys, "\t")
		values = strings.TrimRight(values, "\t")

		w := tabwriter.NewWriter(os.Stdout, 0, 4, 4, ' ', tabwriter.TabIndent)
		fmt.Fprintln(w, keys)
		fmt.Fprintln(w, values)
		w.Flush()
	}

	if success {
		return nil
	}
	return errors.New("some services are not healthy")
}

func messagesError(messages map[string][]string) error {
	str := []string{}

	for k, v := range messages {
		str = append(str, "- "+capitalize(k)+": "+strings.Join(v, ". "))
	}
	return errors.New(strings.Join(str, "\n"))
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
		_ = json.Unmarshal([]byte(bodyStr), &data)
		printList(data)
	case kindApplicationToken:
		data := ApplicationTokens{}
		_ = json.Unmarshal([]byte(bodyStr), &data)
		printList(data)
	case kindPlainToken:
		data := PlainTokens{}
		_ = json.Unmarshal([]byte(bodyStr), &data)
		printList(data)
	case kindTeam:
		data := Teams{}
		_ = json.Unmarshal([]byte(bodyStr), &data)
		printList(data)
	case kindNamespace:
		data := Namespaces{}
		_ = json.Unmarshal([]byte(bodyStr), &data)
		printList(data)
	case kindTag:
		data := Tags{}
		_ = json.Unmarshal([]byte(bodyStr), &data)
		printList(data)
	case kindRegistry:
		data := Registries{}
		_ = json.Unmarshal([]byte(bodyStr), &data)
		printList(data)
	case kindRepository:
		data := Repositories{}
		_ = json.Unmarshal([]byte(bodyStr), &data)
		printList(data)
	case kindValidate:
		data := Validate{}
		_ = json.Unmarshal([]byte(bodyStr), &data)
		return printValidate(data)
	case kindHealth:
		data := Health{}
		_ = json.Unmarshal([]byte(bodyStr), &data)
		return printHealth(data)
	case kindBootstrap:
		data := PlainTokens{}
		_ = json.Unmarshal([]byte(bodyStr), &data)
		return printBootstrap(data[0].PlainToken)
	default:
	}
	return nil
}

func printBootstrap(token string) error {
	fmt.Printf("You can now use the user '%v' with the following token: %v\n",
		retainedValues["username"], token)
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
