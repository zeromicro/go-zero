package main

import (
	"github.com/tal-tech/go-zero/core/lang"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/gen"
)

// go run .
func main() {
	// generating with cache
	withCacheGenerator := gen.NewDefaultGenerator("./test.sql", "./withcachemodel")
	lang.Must(withCacheGenerator.Start(true))

	// generating without cache
	withoutGenerator := gen.NewDefaultGenerator("./test.sql", "./withoutcachemodel")
	lang.Must(withoutGenerator.Start(false))
}
