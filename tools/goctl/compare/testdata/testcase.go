package testdata

import _ "embed"

var (
	//go:embed unformat.api
	unformatApi string
	//go:embed kotlin.api
	kotlinApi string
	//go:embed user.sql
	userSql string

	list = Files{
		{
			IsDir: true,
			Path:  "version",
			Cmd:   "goctl --version",
		},
		{
			IsDir: true,
			Path:  "api/sample_file/local",
			Cmd:   "goctl api --o sample.api",
		},
		{
			IsDir: true,
			Path:  "api/sample_file/local/assign",
			Cmd:   "goctl api --o=sample.api",
		},
		{
			IsDir: true,
			Path:  "api/sample_file/local/assign/shorthand",
			Cmd:   "goctl api -o=sample.api",
		},
		{
			IsDir: true,
			Path:  "api/sample_file/remote",
			Cmd:   "goctl api --o sample.api --remote https://github.com/zeromicro/go-zero-template --branch main",
		},
		{
			IsDir: true,
			Path:  "api/sample_file/remote/shorthand",
			Cmd:   "goctl api -o sample.api -remote https://github.com/zeromicro/go-zero-template -branch main",
		},
		{
			IsDir: true,
			Path:  "api/sample_file/remote/assign",
			Cmd:   "goctl api --o=sample.api --remote https://github.com/zeromicro/go-zero-template --branch=main",
		},
		{
			IsDir: true,
			Path:  "api/sample_file/remote/assign/shorthand",
			Cmd:   "goctl api -o=sample.api -remote https://github.com/zeromicro/go-zero-template -branch=main",
		},
		{
			IsDir: true,
			Path:  "api/dart/legacy/true",
			Cmd:   "goctl api --o sample.api && goctl api dart --api sample.api --dir . --hostname 127.0.0.1 --legacy true",
		},
		{
			IsDir: true,
			Path:  "api/dart/legacy/true/shorthand",
			Cmd:   "goctl api -o sample.api && goctl api dart -api sample.api -dir . -hostname 127.0.0.1 -legacy true",
		},
		{
			IsDir: true,
			Path:  "api/dart/legacy/true/assign",
			Cmd:   "goctl api --o=sample.api && goctl api dart --api=sample.api --dir=. --hostname=127.0.0.1 --legacy=true",
		},
		{
			IsDir: true,
			Path:  "api/dart/legacy/true/assign/shorthand",
			Cmd:   "goctl api -o=sample.api && goctl api dart -api=sample.api -dir=. -hostname=127.0.0.1 -legacy=true",
		},
		{
			IsDir: true,
			Path:  "api/dart/legacy/false",
			Cmd:   "goctl api --o sample.api && goctl api dart --api sample.api --dir . --hostname 127.0.0.1 --legacy true",
		},
		{
			IsDir: true,
			Path:  "api/dart/legacy/false/shorthand",
			Cmd:   "goctl api -o sample.api && goctl api dart -api sample.api -dir . -hostname 127.0.0.1 -legacy true",
		},
		{
			IsDir: true,
			Path:  "api/dart/legacy/false/assign",
			Cmd:   "goctl api --o=sample.api && goctl api dart --api=sample.api --dir=. --hostname=127.0.0.1 --legacy=true",
		},
		{
			IsDir: true,
			Path:  "api/dart/legacy/false/assign/shorthand",
			Cmd:   "goctl api -o=sample.api && goctl api dart -api=sample.api -dir=. -hostname=127.0.0.1 -legacy=true",
		},
		{
			IsDir: true,
			Path:  "api/doc",
			Cmd:   "goctl api --o sample.api && goctl api doc --dir . --o .",
		},
		{
			IsDir: true,
			Path:  "api/doc/shorthand",
			Cmd:   "goctl api -o sample.api && goctl api doc -dir . -o .",
		},
		{
			IsDir: true,
			Path:  "api/doc/assign",
			Cmd:   "goctl api --o=sample.api && goctl api doc --dir=. --o=.",
		},
		{
			IsDir: true,
			Path:  "api/doc/assign/shorthand",
			Cmd:   "goctl api -o=sample.api && goctl api doc -dir=. -o=.",
		},
		{
			Path:    "api/format/unformat.api",
			Content: unformatApi,
			Cmd:     "goctl api format --dir . --iu",
		},
		{
			Path:    "api/format/shorthand/unformat.api",
			Content: unformatApi,
			Cmd:     "goctl api format -dir . -iu",
		},
		{
			Path:    "api/format/assign/unformat.api",
			Content: unformatApi,
			Cmd:     "goctl api format --dir=. --iu",
		},
		{
			Path:    "api/format/assign/shorthand/unformat.api",
			Content: unformatApi,
			Cmd:     "goctl api format -dir=. -iu",
		},
		{
			IsDir: true,
			Path:  "api/go/style/default",
			Cmd:   "goctl api --o sample.api && goctl api go --api sample.api --dir .",
		},
		{
			IsDir: true,
			Path:  "api/go/style/default/shorthand",
			Cmd:   "goctl api -o sample.api && goctl api go -api sample.api -dir .",
		},
		{
			IsDir: true,
			Path:  "api/go/style/assign/default",
			Cmd:   "goctl api --o=sample.api && goctl api go --api=sample.api --dir=.",
		},
		{
			IsDir: true,
			Path:  "api/go/style/assign/default/shorthand",
			Cmd:   "goctl api -o=sample.api && goctl api go -api=sample.api -dir=.",
		},
		{
			IsDir: true,
			Path:  "api/go/style/goZero",
			Cmd:   "goctl api --o sample.api && goctl api go --api sample.api --dir . --style goZero",
		},
		{
			IsDir: true,
			Path:  "api/go/style/goZero/shorthand",
			Cmd:   "goctl api -o sample.api && goctl api go -api sample.api -dir . -style goZero",
		},
		{
			IsDir: true,
			Path:  "api/go/style/goZero/assign",
			Cmd:   "goctl api --o=sample.api && goctl api go --api=sample.api --dir=. --style=goZero",
		},
		{
			IsDir: true,
			Path:  "api/go/style/goZero/assign/shorthand",
			Cmd:   "goctl api -o=sample.api && goctl api go -api=sample.api -dir=. -style=goZero",
		},
		{
			IsDir: true,
			Path:  "api/java",
			Cmd:   "goctl api --o sample.api && goctl api java --api sample.api --dir .",
		},
		{
			IsDir: true,
			Path:  "api/java/shorthand",
			Cmd:   "goctl api -o sample.api && goctl api java -api sample.api -dir .",
		},
		{
			IsDir: true,
			Path:  "api/java/assign",
			Cmd:   "goctl api --o=sample.api && goctl api java --api=sample.api --dir=.",
		},
		{
			IsDir: true,
			Path:  "api/java/shorthand/assign",
			Cmd:   "goctl api -o=sample.api && goctl api java -api=sample.api -dir=.",
		},
		{
			IsDir: true,
			Path:  "api/new/style/default",
			Cmd:   "goctl api new greet",
		},
		{
			IsDir: true,
			Path:  "api/new/style/goZero",
			Cmd:   "goctl api new greet --style goZero",
		},
		{
			IsDir: true,
			Path:  "api/new/style/goZero/assign",
			Cmd:   "goctl api new greet --style=goZero",
		},
		{
			IsDir: true,
			Path:  "api/new/style/goZero/shorthand",
			Cmd:   "goctl api new greet -style goZero",
		},
		{
			IsDir: true,
			Path:  "api/new/style/goZero/shorthand/assign",
			Cmd:   "goctl api new greet -style=goZero",
		},
		{
			IsDir: true,
			Path:  "api/ts",
			Cmd:   "goctl api --o sample.api && goctl api ts --api sample.api --dir . --unwrap --webapi .",
		},
		{
			IsDir: true,
			Path:  "api/ts/shorthand",
			Cmd:   "goctl api -o sample.api && goctl api ts -api sample.api -dir . -unwrap -webapi .",
		},
		{
			IsDir: true,
			Path:  "api/ts/assign",
			Cmd:   "goctl api --o=sample.api && goctl api ts --api=sample.api --dir=. --unwrap --webapi=.",
		},
		{
			IsDir: true,
			Path:  "api/ts/shorthand/assign",
			Cmd:   "goctl api -o=sample.api && goctl api ts -api=sample.api -dir=. -unwrap -webapi=.",
		},
		{
			IsDir: true,
			Path:  "api/validate",
			Cmd:   "goctl api --o sample.api && goctl api validate --api sample.api",
		},
		{
			IsDir: true,
			Path:  "api/validate/shorthand",
			Cmd:   "goctl api -o sample.api && goctl api validate -api sample.api",
		},
		{
			IsDir: true,
			Path:  "api/validate/assign",
			Cmd:   "goctl api --o=sample.api && goctl api validate --api=sample.api",
		},
		{
			IsDir: true,
			Path:  "api/validate/shorthand/assign",
			Cmd:   "goctl api -o=sample.api && goctl api validate -api=sample.api",
		},
		{
			IsDir: true,
			Path:  "env/show",
			Cmd:   "goctl env > env.txt",
		},
		{
			IsDir: true,
			Path:  "env/check",
			Cmd:   "goctl env check -f -v",
		},
		{
			IsDir: true,
			Path:  "env/install",
			Cmd:   "goctl env install -v",
		},
		{
			IsDir: true,
			Path:  "kube",
			Cmd:   "goctl kube deploy --image alpine --name foo --namespace foo --o foo.yaml --port 8888",
		},
		{
			IsDir: true,
			Path:  "kube/shorthand",
			Cmd:   "goctl kube deploy -image alpine -name foo -namespace foo -o foo.yaml -port 8888",
		},
		{
			IsDir: true,
			Path:  "kube/assign",
			Cmd:   "goctl kube deploy --image=alpine --name=foo --namespace=foo --o=foo.yaml --port=8888",
		},
		{
			IsDir: true,
			Path:  "kube/shorthand/assign",
			Cmd:   "goctl kube deploy -image=alpine -name=foo -namespace=foo -o=foo.yaml -port=8888",
		},
		{
			IsDir: true,
			Path:  "model/mongo/cache",
			Cmd:   "goctl model mongo --dir . --type user --style goZero -c",
		},
		{
			IsDir: true,
			Path:  "model/mongo/cache/shorthand",
			Cmd:   "goctl model mongo -dir . -type user -style goZero -c",
		},
		{
			IsDir: true,
			Path:  "model/mongo/cache/assign",
			Cmd:   "goctl model mongo --dir=. --type=user --style=goZero -c",
		},
		{
			IsDir: true,
			Path:  "model/mongo/cache/shorthand/assign",
			Cmd:   "goctl model mongo -dir=. -type=user -style=goZero -c",
		},
		{
			IsDir: true,
			Path:  "model/mongo/nocache",
			Cmd:   "goctl model mongo --dir . --type user",
		},
		{
			IsDir: true,
			Path:  "model/mongo/nocache/shorthand",
			Cmd:   "goctl model mongo -dir . -type user",
		},
		{
			IsDir: true,
			Path:  "model/mongo/nocache/assign",
			Cmd:   "goctl model mongo --dir=. --type=user",
		},
		{
			IsDir: true,
			Path:  "model/mongo/nocache/shorthand/assign",
			Cmd:   "goctl model mongo -dir=. -type=user",
		},
		{
			Content: userSql,
			Path:    "model/mysql/ddl/user.sql",
			Cmd:     "goctl model mysql ddl --database user --dir cache --src user.sql -c",
		},
		{
			Content: userSql,
			Path:    "model/mysql/ddl/shorthand/user.sql",
			Cmd:     "goctl model mysql ddl -database user -dir cache -src user.sql -c",
		},
		{
			Content: userSql,
			Path:    "model/mysql/ddl/assign/user.sql",
			Cmd:     "goctl model mysql ddl --database=user --dir=cache --src=user.sql -c",
		},
		{
			Content: userSql,
			Path:    "model/mysql/ddl/shorthand/assign/user.sql",
			Cmd:     "goctl model mysql ddl -database=user -dir=cache -src=user.sql -c",
		},
		{
			Content: userSql,
			Path:    "model/mysql/ddl/user.sql",
			Cmd:     "goctl model mysql ddl --database user --dir nocache --src user.sql",
		},
		{
			Content: userSql,
			Path:    "model/mysql/ddl/shorthand/user.sql",
			Cmd:     "goctl model mysql ddl -database user -dir nocache -src user.sql",
		},
		{
			Content: userSql,
			Path:    "model/mysql/ddl/assign/user.sql",
			Cmd:     "goctl model mysql ddl --database=user --dir=nocache --src=user.sql",
		},
		{
			Content: userSql,
			Path:    "model/mysql/ddl/shorthand/assign/user.sql",
			Cmd:     "goctl model mysql ddl -database=user -dir=nocache -src=user.sql",
		},
		{
			IsDir: true,
			Path:  "model/mysql/datasource",
			Cmd:   `goctl model mysql datasource --url $DSN --dir cache --table "*" -c`,
		},
		{
			IsDir: true,
			Path:  "model/mysql/datasource/shorthand",
			Cmd:   `goctl model mysql datasource -url $DSN -dir cache -table "*" -c`,
		},
		{
			IsDir: true,
			Path:  "model/mysql/datasource/shorthand2",
			Cmd:   `goctl model mysql datasource -url $DSN -dir cache -t "*" -c`,
		},
		{
			IsDir: true,
			Path:  "model/mysql/datasource/assign",
			Cmd:   `goctl model mysql datasource --url=$DSN --dir=cache --table="*" -c`,
		},
		{
			IsDir: true,
			Path:  "model/mysql/datasource/shorthand/assign",
			Cmd:   `goctl model mysql datasource -url=$DSN -dir=cache -table="*" -c`,
		},
		{
			IsDir: true,
			Path:  "model/mysql/datasource/shorthand2/assign",
			Cmd:   `goctl model mysql datasource -url=$DSN -dir=cache -t="*" -c`,
		},
		{
			IsDir: true,
			Path:  "model/mysql/datasource",
			Cmd:   `goctl model mysql datasource --url $DSN --dir nocache --table "*" -c`,
		},
		{
			IsDir: true,
			Path:  "model/mysql/datasource/shorthand",
			Cmd:   `goctl model mysql datasource -url $DSN -dir nocache -table "*" -c`,
		},
		{
			IsDir: true,
			Path:  "model/mysql/datasource/shorthand2",
			Cmd:   `goctl model mysql datasource -url $DSN -dir nocache -t "*" -c`,
		},
		{
			IsDir: true,
			Path:  "model/mysql/datasource/assign",
			Cmd:   `goctl model mysql datasource --url=$DSN --dir=nocache --table="*" -c`,
		},
		{
			IsDir: true,
			Path:  "model/mysql/datasource/shorthand/assign",
			Cmd:   `goctl model mysql datasource -url=$DSN -dir=nocache -table="*" -c`,
		},
		{
			IsDir: true,
			Path:  "model/mysql/datasource/shorthand2/assign",
			Cmd:   `goctl model mysql datasource -url=$DSN -dir=nocache -t="*" -c`,
		},
		{
			IsDir: true,
			Path:  "rpc/new",
			Cmd:   "goctl rpc new greet",
		},
		{
			IsDir: true,
			Path:  "rpc/template",
			Cmd:   "goctl rpc template --o greet.proto",
		},
		{
			IsDir: true,
			Path:  "rpc/template/shorthand",
			Cmd:   "goctl rpc template -o greet.proto",
		},
		{
			IsDir: true,
			Path:  "rpc/template/assign",
			Cmd:   "goctl rpc template --o=greet.proto",
		},
		{
			IsDir: true,
			Path:  "rpc/template/shorthand/assign",
			Cmd:   "goctl rpc template -o=greet.proto",
		},
		{
			IsDir: true,
			Path:  "rpc/protoc",
			Cmd:   "goctl rpc template --o greet.proto && goctl rpc protoc greet.proto --go_out . --go-grpc_out . --zrpc_out .",
		},
		{
			IsDir: true,
			Path:  "rpc/protoc/assign",
			Cmd:   "goctl rpc template --o=greet.proto && goctl rpc protoc greet.proto --go_out=. --go-grpc_out=. --zrpc_out=.",
		},
	}
)
