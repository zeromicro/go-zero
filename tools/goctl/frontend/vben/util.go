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

package vben

import "strings"

func ConvertGoTypeToTsType(goType string) string {
	switch goType {
	case "int", "uint", "int8", "uint8", "int16", "uint16", "int32", "uint32", "uint64", "int64", "float", "float32", "float64":
		goType = "number"
	case "[]int", "[]uint", "[]int32", "[]int64", "[]uint32", "[]uint64", "[]float", "[]float32", "[]float64":
		goType = "number[]"
	case "string":
		goType = "string"
	case "[]string":
		goType = "string[]"
	case "bool":
		goType = "boolean"

	}
	return goType
}

func FindBeginEndOfLocaleField(data, target string) (int, int) {

	begin := strings.Index(data, target)

	if begin == -1 {
		return -1, -1
	}

	var end int

	for i := begin; i < len(data); i++ {
		if data[i] == '}' {
			end = i + 2
			break
		}
	}

	return begin, end
}
