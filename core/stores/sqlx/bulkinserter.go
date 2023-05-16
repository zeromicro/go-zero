package sqlx

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/executors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stringx"
)

const (
	flushInterval = time.Second
	maxBulkRows   = 1000
	valuesKeyword = "values"
)

var emptyBulkStmt bulkStmt

type (
	// ResultHandler defines the method of result handlers.
	ResultHandler func(sql.Result, error)

	// A BulkInserter is used to batch insert records.
	// Postgresql is not supported yet, because of the sql is formated with symbol `$`.
	// Oracle is not supported yet, because of the sql is formated with symbol `:`.
	BulkInserter struct {
		executor *executors.PeriodicalExecutor
		inserter *dbInserter
		stmt     bulkStmt
	}

	bulkStmt struct {
		prefix      string
		valueFormat string
		suffix      string
	}
)

// NewBulkInserter returns a BulkInserter.
func NewBulkInserter(sqlConn SqlConn, stmt string) (*BulkInserter, error) {
	bkStmt, err := parseInsertStmt(stmt)
	if err != nil {
		return nil, err
	}

	inserter := &dbInserter{
		sqlConn: sqlConn,
		stmt:    bkStmt,
	}

	return &BulkInserter{
		executor: executors.NewPeriodicalExecutor(flushInterval, inserter),
		inserter: inserter,
		stmt:     bkStmt,
	}, nil
}

// Flush flushes all the pending records.
func (bi *BulkInserter) Flush() {
	bi.executor.Flush()
}

// Insert inserts given args.
func (bi *BulkInserter) Insert(args ...any) error {
	value, err := format(bi.stmt.valueFormat, args...)
	if err != nil {
		return err
	}

	bi.executor.Add(value)

	return nil
}

// SetResultHandler sets the given handler.
func (bi *BulkInserter) SetResultHandler(handler ResultHandler) {
	bi.executor.Sync(func() {
		bi.inserter.resultHandler = handler
	})
}

// UpdateOrDelete runs update or delete queries, which flushes pending records first.
func (bi *BulkInserter) UpdateOrDelete(fn func()) {
	bi.executor.Flush()
	fn()
}

// UpdateStmt updates the insert statement.
func (bi *BulkInserter) UpdateStmt(stmt string) error {
	bkStmt, err := parseInsertStmt(stmt)
	if err != nil {
		return err
	}

	bi.executor.Flush()
	bi.executor.Sync(func() {
		bi.inserter.stmt = bkStmt
	})

	return nil
}

type dbInserter struct {
	sqlConn       SqlConn
	stmt          bulkStmt
	values        []string
	resultHandler ResultHandler
}

func (in *dbInserter) AddTask(task any) bool {
	in.values = append(in.values, task.(string))
	return len(in.values) >= maxBulkRows
}

func (in *dbInserter) Execute(bulk any) {
	values := bulk.([]string)
	if len(values) == 0 {
		return
	}

	stmtWithoutValues := in.stmt.prefix
	valuesStr := strings.Join(values, ", ")
	stmt := strings.Join([]string{stmtWithoutValues, valuesStr}, " ")
	if len(in.stmt.suffix) > 0 {
		stmt = strings.Join([]string{stmt, in.stmt.suffix}, " ")
	}
	result, err := in.sqlConn.Exec(stmt)
	if in.resultHandler != nil {
		in.resultHandler(result, err)
	} else if err != nil {
		logx.Errorf("sql: %s, error: %s", stmt, err)
	}
}

func (in *dbInserter) RemoveAll() any {
	values := in.values
	in.values = nil
	return values
}

func parseInsertStmt(stmt string) (bulkStmt, error) {
	lower := strings.ToLower(stmt)
	pos := strings.Index(lower, valuesKeyword)
	if pos <= 0 {
		return emptyBulkStmt, fmt.Errorf("bad sql: %q", stmt)
	}

	var columns int
	right := strings.LastIndexByte(lower[:pos], ')')
	if right > 0 {
		left := strings.LastIndexByte(lower[:right], '(')
		if left > 0 {
			values := lower[left+1 : right]
			values = stringx.Filter(values, func(r rune) bool {
				return r == ' ' || r == '\t' || r == '\r' || r == '\n'
			})
			fields := strings.FieldsFunc(values, func(r rune) bool {
				return r == ','
			})
			columns = len(fields)
		}
	}

	var variables int
	var valueFormat string
	var suffix string
	left := strings.IndexByte(lower[pos:], '(')
	if left > 0 {
		right = strings.IndexByte(lower[pos+left:], ')')
		if right > 0 {
			values := lower[pos+left : pos+left+right]
			for _, x := range values {
				if x == '?' {
					variables++
				}
			}
			valueFormat = stmt[pos+left : pos+left+right+1]
			suffix = strings.TrimSpace(stmt[pos+left+right+1:])
		}
	}

	if variables == 0 {
		return emptyBulkStmt, fmt.Errorf("no variables: %q", stmt)
	}
	if columns > 0 && columns != variables {
		return emptyBulkStmt, fmt.Errorf("columns and variables mismatch: %q", stmt)
	}

	return bulkStmt{
		prefix:      stmt[:pos+len(valuesKeyword)],
		valueFormat: valueFormat,
		suffix:      suffix,
	}, nil
}
