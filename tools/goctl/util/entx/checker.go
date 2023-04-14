// Copyright 2023 The Ryan SU Authors. All Rights Reserved.
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

package entx

import (
	"strings"
)

// IsTimeProperty returns true when the string contains time suffix
func IsTimeProperty(prop string) bool {
	if prop == "time.Time" {
		return true
	}
	return false
}

// IsUpperProperty returns true when the string
// contains Ent upper string such as uuid, api and id
func IsUpperProperty(prop string) bool {
	prop = strings.ToLower(prop)

	if strings.Contains(prop, "uuid") || strings.Contains(prop, "api") ||
		strings.Contains(prop, "id") || strings.Contains(prop, "url") {
		return true
	}

	return false
}

// IsBaseProperty returns true when prop name is
// id, created_at, updated_at, deleted_at
func IsBaseProperty(prop string) bool {
	if prop == "id" || prop == "created_at" || prop == "updated_at" || prop == "deleted_at" {
		return true
	}
	return false
}

// IsGoTypeNotPrototype returns true when property type is
// prototype but not go type
func IsGoTypeNotPrototype(prop string) bool {
	if prop == "int" || prop == "uint" || prop == "[16]byte" {
		return true
	}
	return false
}

// IsUUIDType returns true when prop is Ent's UUID type
func IsUUIDType(prop string) bool {
	if prop == "[16]byte" {
		return true
	}
	return false
}

// IsOnlyEntType returns true when the type is only in ent schema. e.g. uint8
func IsOnlyEntType(t string) bool {
	switch t {
	case "int8", "uint8", "int16", "uint16":
		return true
	default:
		return false
	}
}
