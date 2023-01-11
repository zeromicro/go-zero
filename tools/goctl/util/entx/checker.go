package entx

import (
	"strings"
)

// IsTimeProperty returns true when the string contains time suffix
func IsTimeProperty(prop string) bool {
	if strings.HasSuffix(prop, "_at") || strings.HasSuffix(prop, "_time") {
		return true
	}
	return false
}

// IsUpperProperty returns true when the string
// contains Ent upper string such as uuid, api and id
func IsUpperProperty(prop string) bool {
	if strings.Contains(prop, "uuid") || strings.Contains(prop, "api") ||
		strings.Contains(prop, "id") {
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
