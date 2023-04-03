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

package util

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// ConvertValidateTagToSwagger converts the validator tag to swagger comments.
func ConvertValidateTagToSwagger(tagData string) ([]string, error) {
	if tagData == "" || !strings.Contains(tagData, "validate") {
		return nil, nil
	}

	validateData := ExtractValidateString(tagData)

	return ConvertTagToComment(validateData)
}

// ExtractValidateString extracts the validator's string.
func ExtractValidateString(data string) string {
	beginIndex := strings.Index(data, "validate")
	if beginIndex == -1 {
		return ""
	}
	firstQuotationMark := 0
	for i := beginIndex; i < len(data); i++ {
		if data[i] == '"' && firstQuotationMark == 0 {
			firstQuotationMark = i
		} else if data[i] == '"' && firstQuotationMark != 0 {
			return data[firstQuotationMark+1 : i]
		}
	}
	return ""
}

// ConvertTagToComment converts validator tag to comments.
func ConvertTagToComment(tagString string) ([]string, error) {
	var result []string
	vals := strings.Split(tagString, ",")
	for _, v := range vals {
		if strings.Contains(v, "required") {
			result = append(result, "// required : true\n")
		}

		if strings.Contains(v, "min") || strings.Contains(v, "max") {
			result = append(result, fmt.Sprintf("// %s\n", strings.Replace(v, "=", " length : ", -1)))
		}

		if strings.Contains(v, "len") {
			tagSplit := strings.Split(v, "=")
			_, tagNum := tagSplit[0], tagSplit[1]
			result = append(result, fmt.Sprintf("// max length : %s\n", tagNum))
			result = append(result, fmt.Sprintf("// min length : %s\n", tagNum))
		}

		if strings.Contains(v, "gt") || strings.Contains(v, "gte") ||
			strings.Contains(v, "lt") || strings.Contains(v, "lte") {
			tagSplit := strings.Split(v, "=")
			tag, tagNum := tagSplit[0], tagSplit[1]
			if strings.Contains(tagNum, ".") {
				bitSize := len(tagNum) - strings.Index(tagNum, ".") - 1
				n, err := strconv.ParseFloat(tagNum, bitSize)
				if err != nil {
					return nil, errors.New("failed to convert the number in validate tag")
				}

				switch tag {
				case "gte":
					result = append(result, fmt.Sprintf("// min : %.*f\n", bitSize, n))
				case "gt":
					result = append(result, fmt.Sprintf("// min : %.*f\n", bitSize, n+1/math.Pow(10, float64(bitSize))))
				case "lte":
					result = append(result, fmt.Sprintf("// max : %.*f\n", bitSize, n))
				case "lt":
					result = append(result, fmt.Sprintf("// max : %.*f\n", bitSize, n-1/math.Pow(10, float64(bitSize))))
				}
			} else {
				n, err := strconv.Atoi(tagNum)
				if err != nil {
					return nil, errors.New("failed to convert the number in validate tag")
				}

				switch tag {
				case "gte":
					result = append(result, fmt.Sprintf("// min : %d\n", n))
				case "gt":
					result = append(result, fmt.Sprintf("// min : %d\n", n))
				case "lte":
					result = append(result, fmt.Sprintf("// max : %d\n", n))
				case "lt":
					result = append(result, fmt.Sprintf("// max : %d\n", n))
				}
			}

		}
	}
	return result, nil
}

// HasCustomValidation returns true if the comment has validations.
func HasCustomValidation(data string) bool {
	lowerCase := strings.ToLower(data)
	if strings.Contains(lowerCase, "max") || strings.Contains(lowerCase, "min") ||
		strings.Contains(lowerCase, "required") {
		return true
	}
	return false
}
