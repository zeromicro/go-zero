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

package initlogic

import (
	_ "embed"
	"errors"

	"github.com/spf13/cobra"
)

var (
	// VarStringTarget describes the target.
	VarStringTarget string
	// VarStringModelName describes the model name
	VarStringModelName string
	// VarStringOutputPath describes the output directory
	VarStringOutputPath string
)

func Gen(_ *cobra.Command, _ []string) error {
	err := Validate()
	if err != nil {
		return err
	}

	ctx := &CoreGenContext{
		Target:    VarStringTarget,
		ModelName: VarStringModelName,
		Output:    VarStringOutputPath,
	}

	return DoGen(ctx)
}

func DoGen(g *CoreGenContext) error {
	if g.Target == "core" {
		return GenCore(g)
	} else if g.Target == "other" {
		return OtherGen(g)
	}
	return errors.New("invalid target, try \"core\" or \"other\"")
}

func Validate() error {
	if VarStringTarget == "" {
		return errors.New("the target cannot be empty, use --target to set it")
	} else if VarStringModelName == "" {
		return errors.New("the model name cannot be empty, use --model_name to set it")
	}
	return nil
}
