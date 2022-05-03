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
			Path:  "api/sample_file/remote",
			Cmd:   "goctl api --o sample.api --remote https://github.com/zeromicro/go-zero-template --branch main",
		},
		{
			IsDir: true,
			Path:  "api/dart/legacy/true",
			Cmd:   "goctl api --o sample.api && goctl api dart --api sample.api --dir . --hostname 127.0.0.1 --legacy true",
		},
		{
			IsDir: true,
			Path:  "api/dart/legacy/false",
			Cmd:   "goctl api --o sample.api && goctl api dart --api sample.api --dir . --hostname 127.0.0.1 --legacy true",
		},
		{
			IsDir: true,
			Path:  "api/doc",
			Cmd:   "goctl api --o sample.api && goctl api doc --dir . --o .",
		},
		{
			Path:    "api/format/unformat.api",
			Content: unformatApi,
			Cmd:     "goctl api format --dir . --iu",
		},
		{
			IsDir: true,
			Path:  "api/go/style/default",
			Cmd:   "goctl api --o sample.api && goctl api go --api sample.api --dir .",
		},
		{
			IsDir: true,
			Path:  "api/go/style/goZero",
			Cmd:   "goctl api --o sample.api && goctl api go --api sample.api --dir . --style goZero",
		},
		{
			IsDir: true,
			Path:  "api/java",
			Cmd:   "goctl api --o sample.api && goctl api java --api sample.api --dir .",
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
			Path:  "api/ts",
			Cmd:   "goctl api --o sample.api && goctl api ts --api sample.api --dir . --unwrap --webapi .",
		},
		{
			IsDir: true,
			Path:  "api/validate",
			Cmd:   "goctl api --o sample.api && goctl api validate --api sample.api",
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
			Path:  "model/mongo/cache",
			Cmd:   "goctl model mongo --dir . --type user --style goZero -c",
		},
		{
			IsDir: true,
			Path:  "model/mongo/nocache",
			Cmd:   "goctl model mongo --dir . --type user",
		},
		{
			Content: userSql,
			Path:    "model/mysql/ddl/user.sql",
			Cmd:     "goctl model mysql ddl --database user --dir cache --src user.sql -c",
		},
		{
			Content: userSql,
			Path:    "model/mysql/ddl/user.sql",
			Cmd:     "goctl model mysql ddl --database user --dir nocache --src user.sql",
		},
		{
			IsDir: true,
			Path:  "model/mysql/datasource",
			Cmd:   "goctl model mysql datasource --url $DSN --dir cache --table=* -c",
		},
		{
			IsDir: true,
			Path:  "model/mysql/datasource",
			Cmd:   "goctl model mysql datasource --url $DSN --dir cache --table=*",
		},
		{
			IsDir: true,
			Path:  "model/rpc/new",
			Cmd:   "goctl rpc new greet",
		},
		{
			IsDir: true,
			Path:  "model/rpc/template",
			Cmd:   "goctl rpc template --o greet.proto",
		},
		{
			IsDir: true,
			Path:  "model/rpc/protoc",
			Cmd:   "goctl rpc template --o greet.proto && goctl rpc protoc greet.proto --go_out=. --go-grpc_out=. --zrpc_out=.",
		},
	}
)
