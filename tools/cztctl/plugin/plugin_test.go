package plugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPluginAndArgs(t *testing.T) {
	bin, args := getPluginAndArgs("android")
	assert.Equal(t, "android", bin)
	assert.Equal(t, "", args)

	bin, args = getPluginAndArgs("android=")
	assert.Equal(t, "android", bin)
	assert.Equal(t, "", args)

	bin, args = getPluginAndArgs("android=-javaPackage com.tal")
	assert.Equal(t, "android", bin)
	assert.Equal(t, "-javaPackage com.tal", args)

	bin, args = getPluginAndArgs("android=-javaPackage com.tal --lambda")
	assert.Equal(t, "android", bin)
	assert.Equal(t, "-javaPackage com.tal --lambda", args)

	bin, args = getPluginAndArgs(`https://test-xjy-file.obs.cn-east-2.myhuaweicloud.com/202012/8a7ab6e1-e639-49d1-89cf-2ae6127a1e90n=-v 1`)
	assert.Equal(t, "https://test-xjy-file.obs.cn-east-2.myhuaweicloud.com/202012/8a7ab6e1-e639-49d1-89cf-2ae6127a1e90n", bin)
	assert.Equal(t, "-v 1", args)
}
