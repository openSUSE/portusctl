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

type Users []User

type ApplicationToken struct {
	ID          int64  `json:"id"`
	Application string `json:"application"`
}

func (e ApplicationToken) String() string {
	return tabifyStruct(e)
}

type ApplicationTokens []ApplicationToken

type PlainToken struct {
	PlainToken string `json:"plain_token"`
}

func (e PlainToken) String() string {
	return tabifyStruct(e)
}

type PlainTokens []PlainToken

type Team struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Hidden bool   `json:"hidden"`
}

func (e Team) String() string {
	return tabifyStruct(e)
}

type Teams []Team

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

type Namespaces []Namespace

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

type NamespaceMin struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (a NamespaceMin) String() string {
	return fmt.Sprintf("%v (ID: %v)", a.Name, a.ID)
}

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

type Repositories []Repository

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

type Tags []Tag

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

type Registries []Registry
