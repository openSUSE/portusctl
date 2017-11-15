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
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

const (
	kindApplicationToken = iota + 1
	kindNamespace
	kindPlainToken
	kindRepository
	kindRegistry
	kindTag
	kindTeam
	kindUser
)

// Resource represents a resource as defined by this application.
type Resource struct {
	bundler       string
	kind          int
	optional      []string
	prefix        string
	required      []string
	returnedKind  int
	synonims      []string
	superresource int
	subresources  []string
}

// String returns a pretty version of the resource.
func (r *Resource) String() string {
	var singular string

	if len(r.synonims) == 2 {
		singular = r.synonims[0]
	} else {
		singular = r.synonims[1]
	}

	return strings.Replace(singular, "_", " ", -1)
}

// FullName returns the name to be used for URL paths representing this
// resource.
func (r *Resource) FullName() string {
	return r.synonims[len(r.synonims)-1]
}

// Path returns the full URL path of the given resource with the given extra
// arguments. The resource prefix will also be considered when building this
// path.
func (r *Resource) Path(args []string) string {
	path := filepath.Join(pathPrefix, r.prefix, r.FullName())
	for _, v := range args {
		path = filepath.Join(path, v)
	}
	return path
}

// ReturnedKind returns the type being returned on create actions for the
// current resource. Use this instead of directly accessing the `returnedKind`
// attribute of Resource.
func (r *Resource) ReturnedKind() int {
	if r.returnedKind < kindApplicationToken {
		r.returnedKind = r.kind
	}
	return r.returnedKind
}

// TODO: maybe we should turn that into a map or a proper slice
func findResourceByID(id int) *Resource {
	for _, res := range availableResources {
		if res.kind == id {
			return &res
		}
	}
	return nil
}

func findResource(resource string) *Resource {
	if resource == "" {
		return nil
	}

	for _, res := range availableResources {
		if res.synonims[0][0] == resource[0] {
			for _, part := range res.synonims {
				if part == resource {
					return &res
				} else if len(part) > len(resource) {
					break
				}
			}
		} else if res.synonims[0][0] > resource[0] {
			break
		}
	}
	return nil
}

var availableResources = []Resource{
	{
		// TODO: restrict so it's not used on GET
		required:      []string{"id", "application"},
		kind:          kindApplicationToken,
		returnedKind:  kindPlainToken,
		superresource: kindUser,
		synonims:      []string{"at", "application_token", "application_tokens"},
	},
	{
		kind:         kindNamespace,
		optional:     []string{"description"},
		required:     []string{"name", "team"},
		subresources: []string{"repositories"},
		synonims:     []string{"n", "namespace", "namespaces"},
	},
	{
		kind:         kindRepository,
		subresources: []string{"tags"},
		synonims:     []string{"r", "repository", "repositories"},
	},
	{
		kind:     kindRegistry,
		synonims: []string{"re", "registry", "registries"},
	},
	{
		kind:     kindTag,
		synonims: []string{"tag", "tags"},
	},
	{
		kind:         kindTeam,
		optional:     []string{"description"},
		required:     []string{"name"},
		subresources: []string{"members", "namespaces"},
		synonims:     []string{"t", "team", "teams"},
	},
	{
		bundler:      "user",
		kind:         kindUser,
		optional:     []string{"display_name"},
		required:     []string{"username", "email", "password"},
		subresources: []string{"application_tokens"},
		synonims:     []string{"u", "user", "users"},
	},
}

func create(resource *Resource, args []string) error {
	b, err := generateBody(resource, args)
	if err != nil {
		return err
	}

	res, err := request("POST", resource.Path(nil), b)
	if err != nil {
		return err
	}

	// TODO: quiet flag ?
	fmt.Printf("Created '%v' successfully!\n\n", resource.String())

	body, _ := ioutil.ReadAll(res.Body)
	switch globalConfig.format {
	case jsonFmt:
		fmt.Print(string(body))
	default:
		return prettyPrint(resource.ReturnedKind(), body, true)
	}
	return nil
}

func delete(resource *Resource, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("Expecting exactly 1 argument, %v given", len(args))
	}

	super := findResourceByID(resource.superresource)
	if super != nil {
		resource.prefix = super.FullName()
	}
	res, err := request("DELETE", resource.Path(args), nil)
	if err != nil {
		return err
	}
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Printf("%v\n", string(body))

	// TODO: quiet flag ?
	fmt.Printf("Deleted '%v' successfully!\n", resource.String())
	return nil
}

func get(resource *Resource, args []string) error {
	rsrc, single, err := parseArguments(resource, args)
	if err != nil {
		return err
	}

	res, err := request("GET", resource.Path(args), nil)
	if err != nil {
		return err
	}

	body, _ := ioutil.ReadAll(res.Body)
	switch globalConfig.format {
	case jsonFmt:
		fmt.Print(string(body))
	default:
		return prettyPrint(rsrc.kind, body, single)
	}

	return nil
}

func update(resource *Resource, args []string) error {
	return nil
}
