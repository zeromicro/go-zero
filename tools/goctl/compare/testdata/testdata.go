package testdata

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/logrusorgru/aurora"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

type (
	File struct {
		IsDir        bool
		Path         string
		AbsolutePath string
		Content      string
		Cmd          string
	}

	Files []File
)

func (f File) execute(goctl string) error {
	printDir := f.Path
	dir := f.AbsolutePath
	if !f.IsDir {
		printDir = filepath.Dir(printDir)
		dir = filepath.Dir(dir)
	}
	printCommand := strings.ReplaceAll(fmt.Sprintf("cd %s && %s", printDir, f.Cmd), "goctl", filepath.Base(goctl))
	command := strings.ReplaceAll(fmt.Sprintf("cd %s && %s", dir, f.Cmd), "goctl", goctl)
	fmt.Println(aurora.BrightGreen(printCommand))
	cmd := exec.Command("sh", "-c", command)
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (fs Files) execute(goctl string) error {
	for _, f := range fs {
		err := f.execute(goctl)
		if err != nil {
			return err
		}
	}
	return nil
}

func mustGetTestData(baseDir string) (Files, Files) {
	if len(baseDir) == 0 {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatalln(err)
		}
		baseDir = dir
	}
	baseDir, err := filepath.Abs(baseDir)
	if err != nil {
		return nil, nil
	}
	createFile := func(baseDir string, data File) (File, error) {
		fp := filepath.Join(baseDir, data.Path)
		dir := filepath.Dir(fp)
		if data.IsDir {
			dir = fp
		}
		if err := pathx.MkdirIfNotExist(dir); err != nil {
			return data, err
		}
		data.AbsolutePath = fp
		if data.IsDir {
			return data, nil
		}

		return data, ioutil.WriteFile(fp, []byte(data.Content), os.ModePerm)
	}
	oldDir := filepath.Join(baseDir, "old_fs")
	newDir := filepath.Join(baseDir, "new_fs")
	os.RemoveAll(oldDir)
	os.RemoveAll(newDir)
	var oldFiles, newFiles []File
	for _, data := range list {
		od, err := createFile(oldDir, data)
		if err != nil {
			log.Fatalln(err)
		}
		oldFiles = append(oldFiles, od)
		nd, err := createFile(newDir, data)
		if err != nil {
			log.Fatalln(err)
		}
		newFiles = append(newFiles, nd)
	}
	return oldFiles, newFiles
}

func MustRun(baseDir string) {
	oldFiles, newFiles := mustGetTestData(baseDir)
	goctlOld, err := exec.LookPath("goctl.old")
	must(err)
	goctlNew, err := exec.LookPath("goctl")
	must(err)
	fmt.Println(aurora.BrightBlue("========================goctl.old======================="))
	must(oldFiles.execute(goctlOld))
	fmt.Println()
	fmt.Println(aurora.BrightBlue("========================goctl.new======================="))
	must(newFiles.execute(goctlNew))
}

func must(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
