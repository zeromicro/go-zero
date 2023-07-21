package proto2api

import "strings"

// SkipRpcName returns true when the rpc should be skipped
func SkipRpcName(rpcName, modelName string) bool {
	if rpcName == ("create"+modelName) ||
		rpcName == ("update"+modelName) ||
		rpcName == ("delete"+modelName) ||
		rpcName == ("get"+modelName+"List") ||
		rpcName == ("get"+modelName+"ById") {
		return true
	}
	return false
}

func SkipBaseMessage(messageName string) bool {
	switch messageName {
	case "Empty", "IDReq", "IDsReq", "UUIDsReq", "UUIDReq", "BaseResp", "PageInfoReq",
		"BaseMsg", "BaseIDResp", "BaseUUIDResp":
		return true
	}

	return false
}

func FindTypeContentIndex(data string) int {
	beginIndex := strings.Index(data, "type ")
	for i := beginIndex; i < len(data); i++ {
		if data[i] == '(' {
			return i
		}
	}
	return -1
}

func FindServiceContentIndex(data string) int {
	beginIndex := strings.Index(data, "service ")
	for i := beginIndex; i < len(data); i++ {
		if data[i] == '{' {
			return i
		}
	}
	return -1
}
