package gen

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/logx"
)

var (
	source = "CREATE TABLE `test` (\n  `id` bigint NOT NULL AUTO_INCREMENT,\n  `nanosecond` bigint NOT NULL DEFAULT '0',\n  `data` varchar(255) DEFAULT '',\n  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,\n  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,\n  PRIMARY KEY (`id`),\n  UNIQUE KEY `nanosecond_unique` (`nanosecond`)\n) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;"
)

func TestCacheModel(t *testing.T) {
	logx.Disable()
	_ = Clean()
	g := NewDefaultGenerator(source, "./testmodel")
	err := g.Start(true)
	assert.Nil(t, err)
}
