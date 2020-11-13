package util

import (
	"fmt"
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
)

func DecomposeType(t string) (result []string, err error) {
	add := func(tp string) error {
		ret, err := DecomposeType(tp)
		if err != nil {
			return err
		}

		result = append(result, ret...)
		return nil
	}
	if strings.HasPrefix(t, "map") {
		t = strings.ReplaceAll(t, "map", "")
		if t[0] == '[' {
			pos := strings.Index(t, "]")
			if pos > 1 {
				if err = add(t[1:pos]); err != nil {
					return
				}
				if len(t) > pos+1 {
					err = add(t[pos+1:])
					return
				}
			}
		}
	} else if strings.HasPrefix(t, "[]") {
		if len(t) > 2 {
			err = add(t[2:])
			return
		}
	} else if strings.HasPrefix(t, "*") {
		err = add(t[1:])
		return
	} else {
		result = append(result, t)
		return
	}

	err = fmt.Errorf("bad type %q", t)
	return
}

func GetAllTypes(api *spec.ApiSpec, route spec.Route) []spec.Type {
	var rts []spec.Type
	types := api.Types
	getTypeRecursive(route.RequestType, types, &rts)
	getTypeRecursive(route.ResponseType, types, &rts)
	return rts
}

func GetLocalTypes(api *spec.ApiSpec, route spec.Route) []spec.Type {
	sharedTypes := GetSharedTypes(api)
	isSharedType := func(ty spec.Type) bool {
		for _, item := range sharedTypes {
			if item.Name == ty.Name {
				return true
			}
		}
		return false
	}

	var rts = GetAllTypes(api, route)

	var result []spec.Type
	for _, item := range rts {
		if !isSharedType(item) {
			result = append(result, item)
		}
	}
	return result
}

func getTypeRecursive(ty spec.Type, allTypes []spec.Type, result *[]spec.Type) {
	isCustomType := func(name string) (*spec.Type, bool) {
		for _, item := range allTypes {
			if item.Name == name {
				return &item, true
			}
		}
		return nil, false
	}
	if len(ty.Name) > 0 {
		*result = append(*result, ty)
	}
	for _, member := range ty.Members {
		decomposedItems, _ := DecomposeType(member.Type)
		if len(decomposedItems) == 0 {
			continue
		}
		var customTypes []spec.Type
		for _, item := range decomposedItems {
			c, e := isCustomType(item)
			if e {
				customTypes = append(customTypes, *c)
			}
		}
		for _, ty := range customTypes {
			hasAppend := false
			for _, item := range *result {
				if ty.Name == item.Name {
					hasAppend = true
					break
				}

			}
			if !hasAppend {
				getTypeRecursive(ty, allTypes, result)
			}
		}
	}
}

func GetSharedTypes(api *spec.ApiSpec) []spec.Type {
	types := api.Types
	var result []spec.Type
	var container []spec.Type
	hasInclude := func(all []spec.Type, ty spec.Type) bool {
		for _, item := range all {
			if item.Name == ty.Name {
				return true
			}
		}
		return false
	}
	for _, route := range api.Service.Routes() {
		var rts []spec.Type
		getTypeRecursive(route.RequestType, types, &rts)
		getTypeRecursive(route.ResponseType, types, &rts)
		for _, item := range rts {
			if len(item.Name) == 0 {
				continue
			}
			if hasInclude(container, item) {
				hasAppend := false
				for _, r := range result {
					if item.Name == r.Name {
						hasAppend = true
						break
					}

				}
				if !hasAppend {
					result = append(result, item)
				}
			} else {
				container = append(container, item)
			}
		}
	}
	return result
}
