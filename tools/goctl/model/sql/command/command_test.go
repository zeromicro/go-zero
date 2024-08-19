package command

import (
	_ "embed"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/model/sql/gen"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

var (
	//go:embed testdata/user.sql
	sql string
	cfg = &config.Config{
		NamingFormat: "gozero",
	}
)

func TestFromDDl(t *testing.T) {
	err := gen.Clean()
	assert.Nil(t, err)

	err = fromDDL(ddlArg{
		src:      "./user.sql",
		dir:      pathx.MustTempDir(),
		cfg:      cfg,
		cache:    true,
		database: "go-zero",
		strict:   false,
	})
	assert.Equal(t, errNotMatched, err)

	// case dir is not exists
	unknownDir := filepath.Join(pathx.MustTempDir(), "test", "user.sql")
	err = fromDDL(ddlArg{
		src:      unknownDir,
		dir:      pathx.MustTempDir(),
		cfg:      cfg,
		cache:    true,
		database: "go_zero",
	})
	assert.True(t, func() bool {
		switch err.(type) {
		case *os.PathError:
			return true
		default:
			return false
		}
	}())

	// case empty src
	err = fromDDL(ddlArg{
		dir:      pathx.MustTempDir(),
		cfg:      cfg,
		cache:    true,
		database: "go_zero",
	})
	if err != nil {
		assert.Equal(t, "expected path or path globbing patterns, but nothing found", err.Error())
	}

	tempDir := filepath.Join(pathx.MustTempDir(), "test")
	err = pathx.MkdirIfNotExist(tempDir)
	if err != nil {
		return
	}

	user1Sql := filepath.Join(tempDir, "user1.sql")
	user2Sql := filepath.Join(tempDir, "user2.sql")

	err = os.WriteFile(user1Sql, []byte(sql), os.ModePerm)
	if err != nil {
		return
	}

	err = os.WriteFile(user2Sql, []byte(sql), os.ModePerm)
	if err != nil {
		return
	}

	_, err = os.Stat(user1Sql)
	assert.Nil(t, err)

	_, err = os.Stat(user2Sql)
	assert.Nil(t, err)

	filename := filepath.Join(tempDir, "usermodel.go")
	fromDDL := func(db string) {
		err = fromDDL(ddlArg{
			src:      filepath.Join(tempDir, "user*.sql"),
			dir:      tempDir,
			cfg:      cfg,
			cache:    true,
			database: db,
		})
		assert.Nil(t, err)

		_, err = os.Stat(filename)
		assert.Nil(t, err)
	}

	fromDDL("go_zero")
	_ = os.Remove(filename)
	fromDDL("go-zero")
	_ = os.Remove(filename)
	fromDDL("1gozero")
}

func Test_parseTableList(t *testing.T) {
	testData := []string{"foo", "b*", "bar", "back_up", "foo,bar,b*"}
	patterns := parseTableList(testData)
	actual := patterns.list()
	expected := []string{"foo", "b*", "bar", "back_up"}
	sort.Slice(actual, func(i, j int) bool {
		return actual[i] > actual[j]
	})
	sort.Slice(expected, func(i, j int) bool {
		return expected[i] > expected[j]
	})
	assert.Equal(t, strings.Join(expected, ","), strings.Join(actual, ","))

	matchTestData := map[string]bool{
		"foo":     true,
		"bar":     true,
		"back_up": true,
		"bit":     true,
		"ab":      false,
		"b":       true,
	}
	for v, expected := range matchTestData {
		actual := patterns.Match(v)
		assert.Equal(t, expected, actual)
	}
}
