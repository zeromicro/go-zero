package sqlx

type LogConf struct {
	// Disable Stmt Log, default is `false`.
	DisableStmtLog bool `json:",default=false"`
	// Disable sql log, default is `false`.
	DisableSqlLog bool `json:",default=false"`
}

func SetUp(c LogConf) {
	if c.DisableSqlLog {
		DisableLog()
	} else {
		EnableLog()
	}

	if c.DisableStmtLog {
		DisableStmtLog()
	} else {
		EnableStmtLog()
	}
}
