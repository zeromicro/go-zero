package pathx

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gookit/color"
	"github.com/zeromicro/go-zero/tools/goctl/internal/version"
)

// NL defines a new line.
const (
	NL              = "\n"
	goctlDir        = ".goctl"
	gitDir          = ".git"
	autoCompleteDir = ".auto_complete"
	cacheDir        = "cache"
)

var goctlHome string

// RegisterGoctlHome register goctl home path.
func RegisterGoctlHome(home string) {
	goctlHome = home
}

// CreateIfNotExist creates a file if it is not exists.
func CreateIfNotExist(file string) (*os.File, error) {
	_, err := os.Stat(file)
	if !os.IsNotExist(err) {
		return nil, fmt.Errorf("%s already exist", file)
	}

	return os.Create(file)
}

// RemoveIfExist deletes the specified file if it is exists.
func RemoveIfExist(filename string) error {
	if !FileExists(filename) {
		return nil
	}

	return os.Remove(filename)
}

// RemoveOrQuit deletes the specified file if read a permit command from stdin.
func RemoveOrQuit(filename string) error {
	if !FileExists(filename) {
		return nil
	}

	fmt.Printf("%s exists, overwrite it?\nEnter to overwrite or Ctrl-C to cancel...",
		color.New(color.BgRed, color.Bold).Render(filename))
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	return os.Remove(filename)
}

// FileExists returns true if the specified file is exists.
func FileExists(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}

// FileNameWithoutExt returns a file name without suffix.
func FileNameWithoutExt(file string) string {
	return strings.TrimSuffix(file, filepath.Ext(file))
}

// GetGoctlHome returns the path value of the goctl, the default path is ~/.goctl, if the path has
// been set by calling the RegisterGoctlHome method, the user-defined path refers to.
func GetGoctlHome() (home string, err error) {
	defer func() {
		if err != nil {
			return
		}
		info, err := os.Stat(home)
		if err == nil && !info.IsDir() {
			os.Rename(home, home+".old")
			MkdirIfNotExist(home)
		}
	}()
	if len(goctlHome) != 0 {
		home = goctlHome
		return
	}
	home, err = GetDefaultGoctlHome()
	return
}

// GetDefaultGoctlHome returns the path value of the goctl home where Join $HOME with .goctl.
func GetDefaultGoctlHome() (string, error) {
	var goctlHomeDir = goctlDir
	home, err := os.UserHomeDir()
	if err == nil {
		goctlHomeDir = filepath.Join(home, goctlDir)
	}

	_ = MkdirIfNotExist(goctlHomeDir)
	return goctlHomeDir, nil
}

// GetGitHome returns the git home of goctl.
func GetGitHome() (string, error) {
	goctlH, err := GetGoctlHome()
	if err != nil {
		return "", err
	}

	return filepath.Join(goctlH, gitDir), nil
}

// GetAutoCompleteHome returns the auto_complete home of goctl.
func GetAutoCompleteHome() (string, error) {
	goctlH, err := GetGoctlHome()
	if err != nil {
		return "", err
	}

	return filepath.Join(goctlH, autoCompleteDir), nil
}

// GetCacheDir returns the cache dit of goctl.
func GetCacheDir() (string, error) {
	goctlH, err := GetGoctlHome()
	if err != nil {
		return "", err
	}

	return filepath.Join(goctlH, cacheDir), nil
}

// GetTemplateDir returns the category path value in GoctlHome where could get it by GetGoctlHome.
func GetTemplateDir(category string) (string, error) {
	home, err := GetGoctlHome()
	if err != nil {
		return "", err
	}
	if home == goctlHome {
		// backward compatible, it will be removed in the feature
		// backward compatible start.
		beforeTemplateDir := filepath.Join(home, version.GetGoctlVersion(), category)
		entries, _ := os.ReadDir(beforeTemplateDir)
		infos := make([]fs.FileInfo, 0, len(entries))
		for _, entry := range entries {
			info, err := entry.Info()
			if err != nil {
				continue
			}
			infos = append(infos, info)
		}
		var hasContent bool
		for _, e := range infos {
			if e.Size() > 0 {
				hasContent = true
			}
		}
		if hasContent {
			return beforeTemplateDir, nil
		}
		// backward compatible end.

		return filepath.Join(home, category), nil
	}

	return filepath.Join(home, version.GetGoctlVersion(), category), nil
}

// InitTemplates creates template files GoctlHome where could get it by GetGoctlHome.
func InitTemplates(category string, templates map[string]string) error {
	dir, err := GetTemplateDir(category)
	if err != nil {
		return err
	}

	if err := MkdirIfNotExist(dir); err != nil {
		return err
	}

	for k, v := range templates {
		if err := createTemplate(filepath.Join(dir, k), v, false); err != nil {
			return err
		}
	}

	return nil
}

// CreateTemplate writes template into file even it is exists.
func CreateTemplate(category, name, content string) error {
	dir, err := GetTemplateDir(category)
	if err != nil {
		return err
	}
	return createTemplate(filepath.Join(dir, name), content, true)
}

// Clean deletes all templates and removes the parent directory.
func Clean(category string) error {
	dir, err := GetTemplateDir(category)
	if err != nil {
		return err
	}
	return os.RemoveAll(dir)
}

// LoadTemplate gets template content by the specified file.
func LoadTemplate(category, file, builtin string) (string, error) {
	dir, err := GetTemplateDir(category)
	if err != nil {
		return "", err
	}

	file = filepath.Join(dir, file)
	if !FileExists(file) {
		return builtin, nil
	}

	content, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// SameFile compares the between path if the same path,
// it maybe the same path in case case-ignore, such as:
// /Users/go_zero and /Users/Go_zero, as far as we know,
// this case maybe appear on macOS and Windows.
func SameFile(path1, path2 string) (bool, error) {
	stat1, err := os.Stat(path1)
	if err != nil {
		return false, err
	}

	stat2, err := os.Stat(path2)
	if err != nil {
		return false, err
	}

	return os.SameFile(stat1, stat2), nil
}

func createTemplate(file, content string, force bool) error {
	if FileExists(file) && !force {
		return nil
	}

	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(content)
	return err
}

// MustTempDir creates a temporary directory.
func MustTempDir() string {
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		log.Fatalln(err)
	}

	return dir
}

func Copy(src, dest string) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()

	dir := filepath.Dir(dest)
	err = MkdirIfNotExist(dir)
	if err != nil {
		return err
	}
	w, err := os.Create(dest)
	if err != nil {
		return err
	}
	w.Chmod(os.ModePerm)
	defer w.Close()
	_, err = io.Copy(w, f)
	return err
}

func Hash(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = f.Close()
	}()
	hash := md5.New()
	_, err = io.Copy(hash, f)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
