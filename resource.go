// Copyright (C) 2017-2018 Miquel Sabaté Solà <msabate@suse.com>
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
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
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
	kindValidate
	kindHealth
	kindVersion
)

// Resource represents a resource as defined by this application.
// NOTE: maybe in the future we want to make `subresources` more fine-grained so
// it can also be specified the permissions on the subresource level. For now
// it's fine because then the path that will be built is wrong and so the server
// responds with a http.StatusMethodNotAllowed.
type Resource struct {
	action           int
	bundler          string
	kind             int
	optional         []string
	prefix           string
	required         []string
	returnedKind     int
	synonims         []string
	superresource    int
	subresources     []string
	supportedActions []int
	validate         []string
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

func (r *Resource) PluralName() string {
	str := capitalize(r.FullName())
	return strings.Replace(str, "_", " ", -1)
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
	path := filepath.Join(v1Prefix, r.prefix, r.FullName())
	for _, v := range args {
		path = filepath.Join(path, v)
	}
	return path
}

// ReturnedKind returns the type being returned on create actions for the
// current resource. Use this instead of directly accessing the `returnedKind`
// attribute of Resource.
func (r *Resource) ReturnedKind() int {
	if r.action == validateAction {
		r.returnedKind = kindValidate
	} else if r.returnedKind < kindApplicationToken {
		r.returnedKind = r.kind
	}
	return r.returnedKind
}

// SetAction sets the given action to the resource if possible: resources have a
// constraint on which actions are supported. If the action is not supported,
// then an error is returned. Used this function instead of directly accessing
// the `action` attribute of Resource.
func (r *Resource) SetAction(action int) error {
	for _, v := range r.supportedActions {
		if v == action {
			r.action = action
			return nil
		}
	}
	return fmt.Errorf("Action not supported for resource '%v'", r.String())
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
		required:         []string{"id", "application"},
		kind:             kindApplicationToken,
		returnedKind:     kindPlainToken,
		superresource:    kindUser,
		supportedActions: []int{postAction, deleteAction},
		synonims:         []string{"at", "application_token", "application_tokens"},
	},
	{
		kind:             kindNamespace,
		optional:         []string{"description"},
		required:         []string{"name", "team"},
		subresources:     []string{"repositories"},
		supportedActions: []int{getAction, postAction, validateAction},
		synonims:         []string{"n", "namespace", "namespaces"},
		validate:         []string{"name"},
	},
	{
		kind:             kindRepository,
		subresources:     []string{"tags"},
		supportedActions: []int{getAction},
		synonims:         []string{"r", "repository", "repositories"},
	},
	{
		kind:             kindRegistry,
		supportedActions: []int{getAction, validateAction},
		synonims:         []string{"re", "registry", "registries"},
		validate:         []string{"name", "hostname", "external_hostname", "use_ssl", "only"},
	},
	{
		kind:             kindTag,
		supportedActions: []int{getAction},
		synonims:         []string{"tag", "tags"},
	},
	{
		kind:             kindTeam,
		optional:         []string{"description"},
		required:         []string{"name"},
		subresources:     []string{"members", "namespaces"},
		supportedActions: []int{getAction, postAction},
		synonims:         []string{"t", "team", "teams"},
	},
	{
		bundler:          "user",
		kind:             kindUser,
		optional:         []string{"display_name"},
		required:         []string{"username", "email", "password"},
		subresources:     []string{"application_tokens"},
		supportedActions: []int{getAction, postAction, putAction, deleteAction},
		synonims:         []string{"u", "user", "users"},
	},
}

// createUpdate is a method to share code between the `create` and the `update`
// methods.
func createUpdate(resource *Resource, method string, args, prefix []string) error {
	b, err := generateBody(resource, args)
	if err != nil {
		return err
	}

	res, err := request(method, resource.Path(prefix), "", b)
	if err != nil {
		return err
	}

	if !globalConfig.quiet {
		if method == "POST" {
			fmt.Printf("Created '%v' successfully!\n\n", resource.String())
		} else {
			fmt.Printf("Updated '%v' successfully!\n\n", resource.String())
		}
	}
	return printAndQuit(res, resource.ReturnedKind(), true)
}

func create(resource *Resource, args []string) error {
	return createUpdate(resource, "POST", args, nil)
}

func delete(resource *Resource, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("Expecting exactly 1 argument, %v given", len(args))
	}

	super := findResourceByID(resource.superresource)
	if super != nil {
		resource.prefix = super.FullName()
	}
	res, err := request("DELETE", resource.Path(args), "", nil)
	if err != nil {
		return err
	}
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Printf("%v\n", string(body))

	if !globalConfig.quiet {
		fmt.Printf("Deleted '%v' successfully!\n", resource.String())
	}
	return nil
}

func get(resource *Resource, args []string) error {
	rsrc, single, err := parseArguments(resource, args)
	if err != nil {
		return err
	}

	res, err := request("GET", resource.Path(args), "", nil)
	if err != nil {
		return err
	}
	return printAndQuit(res, rsrc.kind, single)
}

func update(resource *Resource, args []string) error {
	if len(args) < 2 {
		return errors.New("not enough parameters")
	}

	prefix := []string{args[0]}
	return createUpdate(resource, "PUT", args[1:], prefix)
}

func validate(resource *Resource, args []string) error {
	if len(args) < 1 {
		return errors.New("not enough parameters")
	}

	extracted, err := extractArguments(resource, args, true)
	if err != nil {
		return err
	}

	values := url.Values{}
	for k, v := range extracted {
		values.Set(k, v)
	}
	query := values.Encode()

	res, err := request("GET", resource.Path([]string{"validate"}), query, nil)
	if err != nil {
		return err
	}
	return printAndQuit(res, resource.ReturnedKind(), false)
}
