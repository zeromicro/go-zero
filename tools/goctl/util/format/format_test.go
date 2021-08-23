package format

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplit(t *testing.T) {
	list, err := split("A")
	assert.Nil(t, err)
	assert.Equal(t, []string{"A"}, list)

	list, err = split("goZero")
	assert.Nil(t, err)
	assert.Equal(t, []string{"go", "Zero"}, list)

	list, err = split("Gozero")
	assert.Nil(t, err)
	assert.Equal(t, []string{"Gozero"}, list)

	list, err = split("go_zero")
	assert.Nil(t, err)
	assert.Equal(t, []string{"go", "zero"}, list)

	list, err = split("talGo_zero")
	assert.Nil(t, err)
	assert.Equal(t, []string{"tal", "Go", "zero"}, list)

	list, err = split("GOZERO")
	assert.Nil(t, err)
	assert.Equal(t, []string{"G", "O", "Z", "E", "R", "O"}, list)

	list, err = split("gozero")
	assert.Nil(t, err)
	assert.Equal(t, []string{"gozero"}, list)

	list, err = split("")
	assert.Nil(t, err)
	assert.Equal(t, 0, len(list))

	list, err = split("a_b_CD_EF")
	assert.Nil(t, err)
	assert.Equal(t, []string{"a", "b", "C", "D", "E", "F"}, list)

	list, err = split("_")
	assert.Nil(t, err)
	assert.Equal(t, 0, len(list))

	list, err = split("__")
	assert.Nil(t, err)
	assert.Equal(t, 0, len(list))

	list, err = split("_A")
	assert.Nil(t, err)
	assert.Equal(t, []string{"A"}, list)

	list, err = split("_A_")
	assert.Nil(t, err)
	assert.Equal(t, []string{"A"}, list)

	list, err = split("A_")
	assert.Nil(t, err)
	assert.Equal(t, []string{"A"}, list)

	list, err = split("welcome_to_go_zero")
	assert.Nil(t, err)
	assert.Equal(t, []string{"welcome", "to", "go", "zero"}, list)
}

func TestFileNamingFormat(t *testing.T) {
	testFileNamingFormat(t, "gozero", "welcome_to_go_zero", "welcometogozero")
	testFileNamingFormat(t, "_go#zero_", "welcome_to_go_zero", "_welcome#to#go#zero_")
	testFileNamingFormat(t, "Go#zero", "welcome_to_go_zero", "Welcome#to#go#zero")
	testFileNamingFormat(t, "Go#Zero", "welcome_to_go_zero", "Welcome#To#Go#Zero")
	testFileNamingFormat(t, "Go_Zero", "welcome_to_go_zero", "Welcome_To_Go_Zero")
	testFileNamingFormat(t, "go_Zero", "welcome_to_go_zero", "welcome_To_Go_Zero")
	testFileNamingFormat(t, "goZero", "welcome_to_go_zero", "welcomeToGoZero")
	testFileNamingFormat(t, "GoZero", "welcome_to_go_zero", "WelcomeToGoZero")
	testFileNamingFormat(t, "GOZero", "welcome_to_go_zero", "WELCOMEToGoZero")
	testFileNamingFormat(t, "GoZERO", "welcome_to_go_zero", "WelcomeTOGOZERO")
	testFileNamingFormat(t, "GOZERO", "welcome_to_go_zero", "WELCOMETOGOZERO")
	testFileNamingFormat(t, "GO*ZERO", "welcome_to_go_zero", "WELCOME*TO*GO*ZERO")
	testFileNamingFormat(t, "[GO#ZERO]", "welcome_to_go_zero", "[WELCOME#TO#GO#ZERO]")
	testFileNamingFormat(t, "{go###zero}", "welcome_to_go_zero", "{welcome###to###go###zero}")
	testFileNamingFormat(t, "{go###zerogo_zero}", "welcome_to_go_zero", "{welcome###to###go###zerogo_zero}")
	testFileNamingFormat(t, "GogoZerozero", "welcome_to_go_zero", "WelcomegoTogoGogoZerozero")
	testFileNamingFormat(t, "前缀GoZero后缀", "welcome_to_go_zero", "前缀WelcomeToGoZero后缀")
	testFileNamingFormat(t, "GoZero", "welcometogozero", "Welcometogozero")
	testFileNamingFormat(t, "GoZero", "WelcomeToGoZero", "WelcomeToGoZero")
	testFileNamingFormat(t, "gozero", "WelcomeToGoZero", "welcometogozero")
	testFileNamingFormat(t, "go_zero", "WelcomeToGoZero", "welcome_to_go_zero")
	testFileNamingFormat(t, "Go_Zero", "WelcomeToGoZero", "Welcome_To_Go_Zero")
	testFileNamingFormat(t, "Go_Zero", "", "")
	testFileNamingFormatErr(t, "go", "")
	testFileNamingFormatErr(t, "gOZero", "")
	testFileNamingFormatErr(t, "zero", "")
	testFileNamingFormatErr(t, "goZEro", "welcome_to_go_zero")
	testFileNamingFormatErr(t, "goZERo", "welcome_to_go_zero")
	testFileNamingFormatErr(t, "zerogo", "welcome_to_go_zero")
}

func testFileNamingFormat(t *testing.T, format, in, expected string) {
	format, err := FileNamingFormat(format, in)
	assert.Nil(t, err)
	assert.Equal(t, expected, format)
}

func testFileNamingFormatErr(t *testing.T, format, in string) {
	_, err := FileNamingFormat(format, in)
	assert.Error(t, err)
}
