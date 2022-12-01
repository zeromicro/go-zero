import "base.api"

@server(
	group: init
)

service {{.name}} {
	// Initialize database | 初始化数据库
	@handler initDatabase
	get /init/database returns (BaseMsgResp)
}
