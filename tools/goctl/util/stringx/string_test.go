package stringx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString_IsEmptyOrSpace(t *testing.T) {
	ret := From("  ").IsEmptyOrSpace()
	assert.Equal(t, true, ret)
	ret2 := From("ll??").IsEmptyOrSpace()
	assert.Equal(t, false, ret2)
	ret3 := From(`
   		`).IsEmptyOrSpace()
	assert.Equal(t, true, ret3)
}

func TestString_Snake2Camel(t *testing.T) {
	ret := From("____this_is_snake").Snake2Camel()
	assert.Equal(t, "ThisIsSnake", ret)

	ret2 := From("测试_test_Data").Snake2Camel()
	assert.Equal(t, "测试TestData", ret2)

	ret3 := From("___").Snake2Camel()
	assert.Equal(t, "", ret3)

	ret4 := From("testData_").Snake2Camel()
	assert.Equal(t, "TestData", ret4)

	ret5 := From("testDataTestData").Snake2Camel()
	assert.Equal(t, "TestDataTestData", ret5)
}

func TestString_Camel2Snake(t *testing.T) {
	ret := From("ThisIsCCCamel").Camel2Snake()
	assert.Equal(t, "this_is_c_c_camel", ret)

	ret2 := From("测试Test_Data_test_data").Camel2Snake()
	assert.Equal(t, "测试_test__data_test_data", ret2)
}
