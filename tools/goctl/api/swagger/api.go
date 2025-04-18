package swagger

import "github.com/zeromicro/go-zero/tools/goctl/api/spec"

func fillAllStructs(api *spec.ApiSpec) {
	var (
		tps         []spec.Type
		structTypes = make(map[string]spec.DefineStruct)
		groups      []spec.Group
	)
	for _, tp := range api.Types {
		structTypes[tp.Name()] = tp.(spec.DefineStruct)
	}

	for _, tp := range api.Types {
		filledTP := fillStruct("", tp, structTypes)
		tps = append(tps, filledTP)
		structTypes[filledTP.Name()] = filledTP.(spec.DefineStruct)
	}

	for _, group := range api.Service.Groups {
		var routes []spec.Route
		for _, route := range group.Routes {
			route.RequestType = fillStruct("", route.RequestType, structTypes)
			route.ResponseType = fillStruct("", route.ResponseType, structTypes)
			routes = append(routes, route)
		}
		group.Routes = routes
		groups = append(groups, group)
	}
	api.Service.Groups = groups
	api.Types = tps
}

func fillStruct(parent string, tp spec.Type, allTypes map[string]spec.DefineStruct) spec.Type {
	switch val := tp.(type) {
	case spec.DefineStruct:
		var members []spec.Member
		for _, member := range val.Members {
			switch memberType := member.Type.(type) {
			case spec.PointerType:
				member.Type = spec.PointerType{
					RawName: memberType.RawName,
					Type:    fillStruct(val.Name(), memberType.Type, allTypes),
				}
			case spec.ArrayType:
				member.Type = spec.ArrayType{
					RawName: memberType.RawName,
					Value:   fillStruct(val.Name(), memberType.Value, allTypes),
				}
			case spec.MapType:
				member.Type = spec.MapType{
					RawName: memberType.RawName,
					Key:     memberType.Key,
					Value:   fillStruct(val.Name(), memberType.Value, allTypes),
				}
			case spec.DefineStruct:
				if parent != memberType.Name() { // avoid recursive struct
					if st, ok := allTypes[memberType.Name()]; ok {
						member.Type = fillStruct("", st, allTypes)
					}
				}
			case spec.NestedStruct:
				member.Type = fillStruct("", member.Type, allTypes)
			}
			members = append(members, member)
		}
		if len(members) == 0 {
			st, ok := allTypes[val.RawName]
			if ok {
				members = st.Members
			}
		}
		val.Members = members
		return val
	case spec.NestedStruct:
		var members []spec.Member
		for _, member := range val.Members {
			switch memberType := member.Type.(type) {
			case spec.PointerType:
				member.Type = spec.PointerType{
					RawName: memberType.RawName,
					Type:    fillStruct(val.Name(), memberType.Type, allTypes),
				}
			case spec.ArrayType:
				member.Type = spec.ArrayType{
					RawName: memberType.RawName,
					Value:   fillStruct(val.Name(), memberType.Value, allTypes),
				}
			case spec.MapType:
				member.Type = spec.MapType{
					RawName: memberType.RawName,
					Key:     memberType.Key,
					Value:   fillStruct(val.Name(), memberType.Value, allTypes),
				}
			case spec.DefineStruct:
				if parent != memberType.Name() { // avoid recursive struct
					if st, ok := allTypes[memberType.Name()]; ok {
						member.Type = fillStruct("", st, allTypes)
					}
				}
			case spec.NestedStruct:
				if parent != memberType.Name() {
					if st, ok := allTypes[memberType.Name()]; ok {
						member.Type = fillStruct("", st, allTypes)
					}
				}
			}
			members = append(members, member)
		}
		if len(members) == 0 {
			st, ok := allTypes[val.RawName]
			if ok {
				members = st.Members
			}
		}
		val.Members = members
		return val
	case spec.PointerType:
		return spec.PointerType{
			RawName: val.RawName,
			Type:    fillStruct(parent, val.Type, allTypes),
		}
	case spec.ArrayType:
		return spec.ArrayType{
			RawName: val.RawName,
			Value:   fillStruct(parent, val.Value, allTypes),
		}
	case spec.MapType:
		return spec.MapType{
			RawName: val.RawName,
			Key:     val.Key,
			Value:   fillStruct(parent, val.Value, allTypes),
		}
	default:
		return tp
	}
}
