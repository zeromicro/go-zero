package main

import "github.com/zeromicro/go-zero/tools/goctl/compare/cmd"

// EXPRIMENTAL: compare goctl generated code results between old and new.
// TODO: BEFORE RUNNING: export DSN=$datasource
// TODO: AFTER RUNNING: diff --recursive old_fs new_fs

func main() {
	cmd.Execute()
}
