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

import "fmt"

// User represents a Portus user.
type User struct {
	ID          int64  `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Admin       bool   `json:"admin"`
	NamespaceID int64  `json:"namespace_id"`
	DisplayName string `json:"display_name"`
}

func (u User) String() string {
	return tabifyStruct(u)
}

// Users is a list of users.
type Users []User

// ApplicationToken is the representation of an application token.
type ApplicationToken struct {
	ID          int64  `json:"id"`
	Application string `json:"application"`
}

func (e ApplicationToken) String() string {
	return tabifyStruct(e)
}

// ApplicationTokens is a list of application tokens.
type ApplicationTokens []ApplicationToken

// PlainToken is the response given by Portus after creating a valid application
// token.
type PlainToken struct {
	PlainToken string `json:"plain_token"`
}

func (e PlainToken) String() string {
	return tabifyStruct(e)
}

// PlainTokens is a list of plain tokens.
type PlainTokens []PlainToken

// Team is the representation of a team.
type Team struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Hidden bool   `json:"hidden"`
}

func (e Team) String() string {
	return tabifyStruct(e)
}

// Teams is a list of teams.
type Teams []Team

// Namespace is the representation of a namespace.
type Namespace struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	TeamID      int64  `json:"team_id"`
	Visibility  string `json:"visibility"`
	Global      bool   `json:"global"`
}

func (e Namespace) String() string {
	return tabifyStruct(e)
}

// Namespaces is a list of namespaces.
type Namespaces []Namespace

// Author is the information given for tags in regards to the user who pushed
// said tag.
type Author struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (a Author) String() string {
	if a.Name == "" {
		return "<unknown>"
	}
	return fmt.Sprintf("%v (ID: %v)", a.Name, a.ID)
}

// NamespaceMin represents the information as given when fetching repositories
// about the namespace it belongs said repository.
type NamespaceMin struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (a NamespaceMin) String() string {
	return fmt.Sprintf("%v (ID: %v)", a.Name, a.ID)
}

// Repository is the representation of a repository.
// TODO: groupped tags and such
type Repository struct {
	ID        int64        `json:"id"`
	Name      string       `json:"name"`
	FullName  string       `json:"full_name"`
	Namespace NamespaceMin `json:"namespace"`
	Stars     int64        `json:"stars"`
	TagsCount int64        `json:"tags_count"`
}

func (e Repository) String() string {
	return tabifyStruct(e)
}

// Repositories is a list of repositories.
type Repositories []Repository

// Tag is the representation of a tag.
// TODO: vulnerabilities
type Tag struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Author  Author `json:"author"`
	Digest  string `json:"digest"`
	ImageID string `json:"image_id"`
}

func (e Tag) String() string {
	return tabifyStruct(e)
}

// Tags is a list of tags.
type Tags []Tag

// Registry is the representation of a registry.
type Registry struct {
	ID               int64  `json:"id"`
	Name             string `json:"name"`
	Hostname         string `json:"hostname"`
	ExternalHostname string `json:"external_hostname"`
	Secure           bool   `json:"use_ssl"`
}

func (e Registry) String() string {
	return tabifyStruct(e)
}

// Registries is a list of registries.
type Registries []Registry

// Validate is the response from /validate endpoints.
type Validate struct {
	Messages map[string][]string `json:"messages"`
	Valid    bool                `json:"valid"`
}

// HealthStatus is a component of the /health response.
type HealthStatus struct {
	Message string `json:"msg"`
	Success bool   `json:"success"`
}

// Health is the response as given by the server on a /health request.
type Health map[string]HealthStatus

type gitInfo struct {
	Branch string `json:"branch"`
	Commit string `json:"commit"`
	Tag    string `json:"tag"`
}

// Version represents the response as given by Portus on a /version request
// (except for the `PortusctlVersion` field, which is expected to be filled by
// portusctl).
type Version struct {
	APIVersions      []string `json:"api-versions"`
	Git              gitInfo  `json:"git"`
	PortusctlVersion string   `json:"portusctl-version"`
	Version          string   `json:"version"`
}
